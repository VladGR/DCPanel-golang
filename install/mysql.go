package install

import (
	"djcontrol/config"
	"djcontrol/term"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func MySQL(con *ssh.Client, server *config.Server) {
	fmt.Println("Installing MySQL...")

	version := server.MySQL.Version
	password := server.MySQL.User.Password

	// ignore "Enter" manual password
	// password will be set automatically
	cmd := fmt.Sprintf("echo 'mysql-server-%s mysql-server/root_password password %s' | debconf-set-selections && ", version, password)
	cmd += fmt.Sprintf("echo 'mysql-server-%s mysql-server/root_password_again password %s' | debconf-set-selections && ", version, password)
	cmd += fmt.Sprintf("apt-get -y install mysql-server-%s", version)
	term.RunLongCommand(con, cmd)

	term.RunLongCommand(con, "mkdir -p /etc/mysql/conf.d/")

	CopyFileToServer(server, "mysql", "myconf.cnf", "root", "/etc/mysql/conf.d/myconf.cnf")

	term.RunLongCommand(con, "chown -R mysql: /etc/mysql/conf.d/")

	term.RunLongCommand(con, "systemctl enable mysql")
	term.RunLongCommand(con, "systemctl start mysql")

	// On running MySQL create remote access for "root" from any host

	cmd = fmt.Sprintf("mysql -u root -p%s ", password)
	cmd += fmt.Sprintf("-e \"GRANT ALL PRIVILEGES ON *.* TO 'root'@'%%' IDENTIFIED BY '%s' WITH GRANT OPTION;\"", password)
	term.RunLongCommand(con, cmd)

	// to make sure all configs to apply
	term.RunLongCommand(con, "systemctl restart mysql")
}
