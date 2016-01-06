package install

import (
	"djcontrol/config"
	"djcontrol/term"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func Nginx(con *ssh.Client, server *config.Server) {
	fmt.Println("Installing Nginx...")
	cmd := "apt-get -y install nginx && "
	cmd += "systemctl enable nginx && "
	cmd += "systemctl start nginx"
	term.RunLongCommand(con, cmd)

	cmd = "mkdir -p /var/log/nginx/sites/access"
	term.RunLongCommand(con, cmd)

	cmd = "mkdir -p /var/log/nginx/sites/error"
	term.RunLongCommand(con, cmd)

	cmd = fmt.Sprintf("chown -R %s:%s /var/log/nginx/sites", server.NginxName, server.NginxName)
	term.RunLongCommand(con, cmd)

}
