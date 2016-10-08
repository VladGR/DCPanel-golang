package install

import (
	"djcontrol/config"
	"djcontrol/term"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func MySQL(con *ssh.Client, server *config.Server) {
	fmt.Println("Installing MySQL...")

	// Перед установкой MySQL, чтобы вместо 5.5 установить5.7
	// wget http://dev.mysql.com/get/mysql-apt-config_0.6.0-1_all.deb
	// sudo dpkg -i mysql-apt-config_0.6.0-1_all.deb

	// 1 экран - выбираем MySQL Server (mysql 5.7) - нажимаем Enter
	// 2 экран - выбираем mysql-5.7 - нажимаем Enter
	// 3 экран - стрелкой переходим на Apply

	// sudo apt-get update

	// version := server.MySQL.Version
	password := server.MySQL.User.Password

	cmd := fmt.Sprintf("export DEBIAN_FRONTEND=noninteractive && ")
	cmd += fmt.Sprintf("echo 'mysql-apt-config mysql-apt-config/repo-codename select trusty' | debconf-set-selections && ")
	cmd += fmt.Sprintf("echo 'mysql-apt-config mysql-apt-config/repo-distro select ubuntu' | debconf-set-selections && ")
	cmd += fmt.Sprintf("echo 'mysql-apt-config mysql-apt-config/repo-url string http://repo.mysql.com/apt/' | debconf-set-selections && ")
	cmd += fmt.Sprintf("echo 'mysql-apt-config mysql-apt-config/select-preview select ' | debconf-set-selections && ")
	cmd += fmt.Sprintf("echo 'mysql-apt-config mysql-apt-config/select-product select Ok' | debconf-set-selections && ")
	cmd += fmt.Sprintf("echo 'mysql-apt-config mysql-apt-config/select-server select mysql-5.7' | debconf-set-selections && ")
	cmd += fmt.Sprintf("echo 'mysql-apt-config mysql-apt-config/select-tools select ' | debconf-set-selections && ")
	cmd += fmt.Sprintf("echo 'mysql-apt-config mysql-apt-config/unsupported-platform select abort' | debconf-set-selections && ")
	cmd += fmt.Sprintf("wget http://dev.mysql.com/get/mysql-apt-config_0.7.2-1_all.deb && ")
	cmd += fmt.Sprintf("dpkg -i mysql-apt-config_0.7.2-1_all.deb && ")
	cmd += fmt.Sprintf("apt-get update && ")
	cmd += fmt.Sprintf("apt-get -y install mysql-server")
	term.RunLongCommand(con, cmd)

	term.RunLongCommand(con, "mkdir -p /etc/mysql/conf.d/")
	term.RunLongCommand(con, "mkdir -p /etc/mysql/mysql.conf.d/")

	CopyFileToServer(server, "mysql", "mysqld.cnf", "root", "/etc/mysql/mysql.conf.d/mysqld.cnf")

	term.RunLongCommand(con, "chown -R mysql: /etc/mysql/conf.d/")
	term.RunLongCommand(con, "chown -R mysql: /etc/mysql/mysql.conf.d/")

	term.RunLongCommand(con, "systemctl enable mysql")
	term.RunLongCommand(con, "systemctl start mysql")

	// On running MySQL create remote access for "root" from any host

	cmd = fmt.Sprintf("mysql -u root -p%s ", password)
	cmd += fmt.Sprintf("-e \"GRANT ALL PRIVILEGES ON *.* TO 'root'@'%%' IDENTIFIED BY '%s' WITH GRANT OPTION;\"", password)
	term.RunLongCommand(con, cmd)

	// to make sure all configs to apply
	term.RunLongCommand(con, "systemctl restart mysql")
}
