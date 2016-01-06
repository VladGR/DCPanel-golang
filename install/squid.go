package install

import (
	"djcontrol/config"
	"djcontrol/term"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func Squid(con *ssh.Client, server *config.Server) {
	fmt.Println("Installing Squid...")

	CopyFileToServer(server, "squid", "interfaces", "root", "/etc/network/interfaces")
	term.RunLongCommand(con, "/etc/init.d/networking restart")

	cmd := "sudo apt-get install -y squid3 && systemctl enable squid3 && systemctl start squid3"
	term.RunLongCommand(con, cmd)

	CopyFileToServer(server, "squid", "squid.conf", "root", "/etc/squid3/squid.conf")
	term.RunLongCommand(con, "systemctl restart squid3")

	term.RunLongCommand(con, "ifconfig -a")
	term.RunLongCommand(con, "systemctl status squid3")
}
