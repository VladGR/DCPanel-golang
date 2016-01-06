package install

import (
	"djcontrol/config"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func Bash(con *ssh.Client, server *config.Server) {
	fmt.Println("Processing bash...")
	BashAliases(con, server)
}

func BashAliases(con *ssh.Client, server *config.Server) {
	c := config.GetConfig()

	CopyFileToServer(server, "bash", "bash_aliases", "root", "~/.bash_aliases")
	CopyFileToServer(server, "bash", "bash_aliases", c.LocalLinuxUser, "~/.bash_aliases")

	CopyFileToServer(server, "bash", "bashrc", "root", "~/.bashrc")
	CopyFileToServer(server, "bash", "bashrc", c.LocalLinuxUser, "~/.bashrc")
}
