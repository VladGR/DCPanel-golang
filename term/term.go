/*
RunCommand      short command that return string
RunLongCommand  command that works in real-time
*/

package term

import (
	"bytes"
	"djcontrol/config"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func GetConnection(ip string, userName string) *ssh.Client {
	cfg := config.GetConfig()
	con, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), getSSHConfig(userName, cfg.LocalLinuxUser))
	if err != nil {
		checkErr(err)
	}
	return con
}

func RunCommand(con *ssh.Client, cmd string) string {
	var b bytes.Buffer

	session := GetSession(con)
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		checkErr(err)
	}
	return b.String()
}

func RunLongCommand(con *ssh.Client, cmd string) {
	session := GetRealTimeSession(con)
	if err := session.Run(cmd); err != nil {
		checkErr(err)
	}
}

func RunLongCommandIgnoreError(con *ssh.Client, cmd string) {
	session := GetRealTimeSession(con)
	session.Run(cmd)
}

func GetSession(con *ssh.Client) *ssh.Session {
	session, err := con.NewSession()
	if err != nil {
		checkErr(err)
	}
	return session
}

func GetRealTimeSession(con *ssh.Client) *ssh.Session {
	session, err := con.NewSession()
	if err != nil {
		checkErr(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		checkErr(err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		checkErr(err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		checkErr(err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		checkErr(err)
	}
	go io.Copy(os.Stderr, stderr)

	return session
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func getSSHConfig(remoteUserName string, localUserName string) *ssh.ClientConfig {
	p := fmt.Sprintf("/home/%s/.ssh/id_rsa", localUserName)

	sshConfig := &ssh.ClientConfig{
		User: remoteUserName,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(p),
		},
	}
	return sshConfig
}

func checkErr(err error) {
	if err != nil {
		log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
		log.Fatal(err)
	}
}
