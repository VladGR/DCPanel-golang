package install

import (
	"djcontrol/config"
	"djcontrol/funcs"
	"djcontrol/term"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func PHP(con *ssh.Client, server *config.Server) {
	fmt.Println("Installing PHP...")
	php01(con)
	php02(con, server)
	php03(con, server)
	php04(con, server)

}

func php01(con *ssh.Client) {
	cmd := "apt-get -y install php5 php5-fpm php5-common php5-cli php5-gd php5-mcrypt php5-mysql"
	term.RunLongCommand(con, cmd)
}

// Create "php" user
func php02(con *ssh.Client, server *config.Server) {
	fmt.Println("Creating user \"php\" ...")

	term.RunLongCommand(con, fmt.Sprintf("mkdir -p /home/php"))

	// ignore error if user exists
	cmd := fmt.Sprintf("id -u php &>/dev/null || useradd php -d /home/php")
	term.RunLongCommand(con, cmd)
	cmd = "chown php: /home/php"
	term.RunLongCommand(con, cmd)

	// create password
	newPassword := funcs.RandomString(10)
	cmd = fmt.Sprintf("echo php:%s | chpasswd", newPassword)
	term.RunLongCommand(con, cmd)

	// create ssh keys
	cmd = fmt.Sprintf("sshpass -p %s ssh-copy-id -i ~/.ssh/id_rsa.pub php@%s", newPassword, server.Ip)
	// ignore error if keys have been set earlier
	funcs.RunCommandShIgnoreError(cmd)
}

// Create directories
func php03(con *ssh.Client, server *config.Server) {
	fmt.Println("Creating directories ...")

	term.RunLongCommand(con, fmt.Sprintf("mkdir -p /home/php/sites"))

	cmd := fmt.Sprintf("chown php:%s /home/php/sites", server.NginxName)
	term.RunLongCommand(con, cmd)

	cmd = "chmod 0750 /home/php/sites"
	term.RunLongCommand(con, cmd)

	term.RunLongCommand(con, fmt.Sprintf("mkdir -p /home/php/logs/access"))
	term.RunLongCommand(con, fmt.Sprintf("mkdir -p /home/php/logs/error"))

}

// Copy configs and start php-fpm
func php04(con *ssh.Client, server *config.Server) {
	fmt.Println("Copying config files ...")
	CopyFileToServer(server, "php", "php-fpm.conf", "root", "/etc/php5/fpm/php-fpm.conf")
	CopyFileToServer(server, "php", "www.conf", "root", "/etc/php5/fpm/pool.d/www.conf")
	term.RunLongCommand(con, "systemctl start php5-fpm")
	term.RunLongCommand(con, "systemctl status php5-fpm")

}
