package django

import (
	"djcontrol/config"
	"djcontrol/funcs"
	"djcontrol/term"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func Deploy() {
	project, err := config.InputDjangoProject()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conr := term.GetConnection(project.Server.Ip, "root")
	conu := term.GetConnection(project.Server.Ip, project.Server.MainUser)

	fmt.Printf("You are going to deploy %q to server %q %s.\nCopy media files? y/n> ", project.Name, project.Server.Code, project.Server.Ip)

	var doCopyMedia string
	fmt.Scanln(&doCopyMedia)

	fmt.Print("Install virtual environment? y/n> ")
	var doInstallVenv string
	fmt.Scanln(&doInstallVenv)

	prepareLocalFiles(project)
	makeExecutables(project)
	copyProjectFiles(conu, project)

	// create and install venv after project files have been copied
	// need requirements/production.txt file
	if doInstallVenv == "y" {
		createVirtualEnv(conu, project)
		installVirtualEnv(conu, project)
	}

	if project.IsStaticDirSeparate {
		copyStaticFiles(conu, project)
	}

	if doCopyMedia == "y" {
		copyMediaFiles(conu, project)
	}

	applyUserRights(conr, project)
	processNginx(conr, project)
	processSupervisor(conr, project)
	touchReload(conu, project)
}

func prepareLocalFiles(project *config.Project) {
	// set 755 for all directories and 644 for all project's local files

	cmd := fmt.Sprintf("find %s -type d -exec chmod 755 {} +", project.ProjectDirLocal)
	funcs.RunCommandSh(cmd)

	cmd = fmt.Sprintf("find %s -type f -exec chmod 644 {} +", project.ProjectDirLocal)
	funcs.RunCommandSh(cmd)
}

func makeExecutables(project *config.Project) {
	for _, file := range project.Executables {
		file = project.ProjectDirLocal + file
		file = path.Clean(file)
		cmd := fmt.Sprintf("chmod +x %s", file)
		funcs.RunCommandSh(cmd)
	}
}

func createVirtualEnv(con *ssh.Client, project *config.Project) {
	fmt.Println("Creating virtual environment...")
	s := fmt.Sprintf("virtualenv %s --python=%s", project.VenvDirServer, project.PythonPathServer)

	// create virtualenv if directory doesn't exist
	cmd := fmt.Sprintf("[ -d \"%s\" ] || %s", project.VenvDirServer, s)
	term.RunLongCommand(con, cmd)
}

func installVirtualEnv(con *ssh.Client, project *config.Project) {
	fmt.Println("Installing virtual environment...")

	reqPath := fmt.Sprintf("%s%sproduction.txt", project.ProjectDirServer, project.RequirementsDir)
	reqPath = path.Clean(reqPath)
	cmd := fmt.Sprintf(". %sbin/activate && pip install -r %s", project.VenvDirServer, reqPath)
	term.RunLongCommand(con, cmd)
}

// copy project files excluding media directory
func copyProjectFiles(con *ssh.Client, project *config.Project) {
	fmt.Println("Copying project files...")

	cmd := fmt.Sprintf("mkdir -p %s", project.ProjectDirServer)
	term.RunLongCommand(con, cmd)

	ex := project.MediaDirLocal

	cmd = "rsync -avzh "
	for _, path := range project.Exclude {
		cmd += fmt.Sprintf(" --exclude %s", path)
	}

	cmd += fmt.Sprintf(" --exclude %s -e ssh %s %s@%s:%s", ex, project.ProjectDirLocal, project.Server.MainUser, project.Server.Ip, project.ProjectDirServer)
	fmt.Println(cmd)
	funcs.RunCommandSh(cmd)
}

// copy static files excluding media directory if media inside static
func copyStaticFiles(con *ssh.Client, project *config.Project) {
	fmt.Println("Copying static files...")

	cmd := fmt.Sprintf("mkdir -p %s", project.StaticDirServer)
	term.RunLongCommand(con, cmd)

	cmd = "rsync -avzh "

	if strings.Contains(project.MediaDirLocal, project.StaticDirLocal) {
		ex := strings.Replace(project.MediaDirLocal, project.StaticDirLocal, "", 1)
		cmd += fmt.Sprintf(" --exclude %s ", ex)
	}

	fmt.Println(project.StaticDirLocal)
	pathLocal := project.ProjectDirLocal + strings.TrimLeft(project.StaticDirLocal, "/")

	cmd += fmt.Sprintf(" -e ssh %s %s@%s:%s", pathLocal, project.Server.MainUser, project.Server.Ip, project.StaticDirServer)
	fmt.Println(cmd)
	funcs.RunCommandSh(cmd)
}

func copyMediaFiles(con *ssh.Client, project *config.Project) {
	fmt.Println("Copying media files...")
	cmd := fmt.Sprintf("mkdir -p %s", project.MediaDirServer)
	term.RunLongCommand(con, cmd)

	localPath := project.ProjectDirLocal + strings.TrimLeft(project.MediaDirLocal, "/")

	cmd = fmt.Sprintf("rsync -avzh -e ssh %s %s@%s:%s", localPath, project.Server.MainUser, project.Server.Ip, project.MediaDirServer)
	fmt.Println(cmd)
	funcs.RunCommandSh(cmd)
}

func get1stLevelDir(somePath string) string {
	ptn := `^/([A-Za-z0-9_-]+)/`
	re := regexp.MustCompile(ptn)
	return re.FindString(somePath)
}

func applyUserRights(con *ssh.Client, project *config.Project) {
	cmd := fmt.Sprintf("chown -R %s: %s", project.Server.MainUser, project.ProjectDirServer)
	term.RunLongCommand(con, cmd)

	cmd = fmt.Sprintf("chown -R %s:%s %s", project.Server.MainUser, project.Server.NginxName, project.StaticDirServer)
	term.RunLongCommand(con, cmd)

	if project.IsStaticDirSeparate {
		/*
			If static directory is separate, nginx should not have access to project
			folder. Project folder is accessed only by user.
		*/
		cmd := fmt.Sprintf("chmod 0700 %s", project.ProjectDirServer)
		term.RunLongCommand(con, cmd)

		// Nginx should have full access to static directory (as Group member)
		cmd = "chmod -R 0770 " + project.StaticDirServer
		term.RunLongCommand(con, cmd)

	} else {
		/*
			If static directory inside project - the 1st level directory
			inside main project directory should have access only for user and nobody else.
			We can take 1st level directory for example from "requirements" directory
			that is always on 2nd level.
		*/

		// "Others" need read access to project directory
		cmd := fmt.Sprintf("chmod 0755 %s", project.ProjectDirServer)
		term.RunLongCommand(con, cmd)

		// Nginx needs full access to static directory inside project (as Group member)
		cmd = fmt.Sprintf("chmod 0770 %s", project.StaticDirServer)
		term.RunLongCommand(con, cmd)

		// 1st level directory inside project should be accessible only by user.
		fDir := get1stLevelDir(project.RequirementsDir)
		fDir = project.ProjectDirServer + fDir

		cmd = fmt.Sprintf("chmod 0700 %s", fDir)
		term.RunLongCommand(con, cmd)
	}

}

func processNginx(con *ssh.Client, project *config.Project) {
	fmt.Println("Processing Nginx...")
	for _, conf := range project.NginxConfs {
		remotePath := "/etc/nginx/conf.d/" + conf.FileName
		copyFileToServer(project.Server, conf, "root", remotePath)
	}

	cmd := "sudo systemctl reload nginx"
	term.RunLongCommand(con, cmd)
}

func processSupervisor(con *ssh.Client, project *config.Project) {
	fmt.Println("Processing Supervisor...")

	for _, conf := range project.SupervisorConfs {
		remotePath := "/etc/supervisord/conf.d/" + conf.FileName
		copyFileToServer(project.Server, conf, "root", remotePath)
	}

	time.Sleep(3 * time.Second)
	term.RunLongCommand(con, "supervisorctl reread")
	term.RunLongCommand(con, "supervisorctl update")

	// get program names and restart them
	for _, conf := range project.SupervisorConfs {
		re := regexp.MustCompile(`\[program\:([a-z0-9_-]+)\]`)
		list := re.FindStringSubmatch(conf.Data)

		if len(list) > 1 {
			progName := list[1]
			term.RunLongCommand(con, fmt.Sprintf("supervisorctl restart %s", progName))
		}
	}

	time.Sleep(3 * time.Second)
	term.RunLongCommand(con, "supervisorctl status")
}

func touchReload(con *ssh.Client, project *config.Project) {
	fmt.Println("Reloading Django Project...")
	touchPath := project.ProjectDirServer + strings.TrimLeft(project.ReloadIniPath, "/")
	cmd := fmt.Sprintf("touch %s", touchPath)
	term.RunLongCommand(con, cmd)
}

func copyFileToServer(server *config.Server, conf *config.ProjectConf, remoteUser string, remotePath string) {

	data := conf.Data
	data = strings.Replace(data, "\r", "", -1) // important
	ioutil.WriteFile(conf.FileName, []byte(data), 0644)

	cmd := fmt.Sprintf("scp -v %s %s@%s:%s", conf.FileName, remoteUser, server.Ip, remotePath)
	funcs.RunCommand(cmd)
	os.Remove(conf.FileName)
}
