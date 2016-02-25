package install

import (
	"djcontrol/config"
	"djcontrol/term"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func Supervisor(con *ssh.Client, server *config.Server) {
	fmt.Println("Installing Supervisor...")

	term.RunLongCommand(con, "pip install supervisor")

	CopyFileToServer(server, "supervisor", "supervisord.conf", "root", "/etc/supervisord.conf")
	CopyFileToServer(server, "supervisor", "supervisord.service", "root", "/usr/lib/systemd/system/supervisord.service")

	term.RunLongCommand(con, "touch /tmp/supervisord.log")
	term.RunLongCommand(con, fmt.Sprintf("chown %s: /tmp/supervisord.log", server.MainUser))

	term.RunLongCommand(con, "mkdir -p /etc/supervisord/")
	term.RunLongCommand(con, "mkdir -p /etc/supervisord/conf.d")

	term.RunLongCommand(con, "sudo systemctl enable supervisord")
	term.RunLongCommand(con, "sudo systemctl start supervisord")

	// create directory for uwsgi logs
	term.RunLongCommand(con, "mkdir -p /var/log/uwsgi/sites")
	term.RunLongCommand(con, fmt.Sprintf("chown %s: /var/log/uwsgi/sites", server.MainUser))
	term.RunLongCommand(con, "chmod 0750 /var/log/uwsgi/sites")

}
