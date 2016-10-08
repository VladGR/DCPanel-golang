package install

import (
	"djcontrol/config"
	"djcontrol/funcs"
	"djcontrol/term"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Start() {
	server, err := config.InputServer()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	config.PrepareServerToInstall(server)

	con := term.GetConnection(server.Ip, "root")

	CreateUser(con, server)
	Base(con)
	DropCache(con, server)

	// if funcs.IsSliceContainsString(server.Installs, "iptables") {
	// 	IPTables(con, server)
	// }

	// if funcs.IsSliceContainsString(server.Installs, "bash") {
	// 	Bash(con, server)
	// }

	// if funcs.IsSliceContainsString(server.Installs, "python2") {
	// 	Python2(con)
	// }

	// if funcs.IsSliceContainsString(server.Installs, "python3") {
	// 	Python3(con)
	// }

	// if funcs.IsSliceContainsString(server.Installs, "nginx") {
	// 	Nginx(con, server)
	// }

	// if funcs.IsSliceContainsString(server.Installs, "redis") {
	// 	Redis(con)
	// }

	if funcs.IsSliceContainsString(server.Installs, "postgresql") {
		PostgreSQL(con, server)
	}

	// if funcs.IsSliceContainsString(server.Installs, "mysql") {
	// 	MySQL(con, server)
	// }

	// if funcs.IsSliceContainsString(server.Installs, "supervisor") {
	// 	Supervisor(con, server)
	// }

	// if funcs.IsSliceContainsString(server.Installs, "squid") {
	// 	Squid(con, server)
	// }

	// if funcs.IsSliceContainsString(server.Installs, "postfix") {
	// 	Postfix(con, server)
	// }

	// if funcs.IsSliceContainsString(server.Installs, "php") {
	// 	PHP(con, server)
	// }

	fmt.Println("Installation complete!!!")
}

// Get file by server, item and filename from api
func CopyFileToServer(server *config.Server, item string, filename string, remoteUser string, remotePath string) {
	conf := config.GetServerConf(server.Id, item, filename)

	data := conf.Data
	data = strings.Replace(data, "\r", "", -1) // important
	ioutil.WriteFile(filename, []byte(data), 0644)

	cmd := fmt.Sprintf("scp -v %s %s@%s:%s", filename, remoteUser, server.Ip, remotePath)
	funcs.RunCommand(cmd)
	os.Remove(filename)
}

// Create temp file from data passed to function
func CopyTempFileToServer(server *config.Server, data string, remoteUser string, remotePath string) {
	filename := "temp.txt"
	ioutil.WriteFile(filename, []byte(data), 0644)

	cmd := fmt.Sprintf("scp -v %s %s@%s:%s", filename, remoteUser, server.Ip, remotePath)
	funcs.RunCommand(cmd)
	os.Remove(filename)
}
