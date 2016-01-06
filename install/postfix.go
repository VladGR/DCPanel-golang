package install

import (
	"djcontrol/config"
	"djcontrol/funcs"
	"djcontrol/term"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func Postfix(con *ssh.Client, server *config.Server) {
	Postfix01(con, server)
	Postfix02(con, server)
	Postfix03(con, server)
	Postfix04(con, server)
	Postfix05(con, server)
	Postfix06(con, server)
	Postfix07(con, server)
	Postfix08(con, server)
	Postfix09(con, server)
	Postfix10(con, server)
	Postfix11(con, server)
	Postfix12(con, server)
	Postfix13(con, server)
	Postfix14(con, server)
}

// install postfix
func Postfix01(con *ssh.Client, server *config.Server) {
	fmt.Println("Installing Postfix...")

	// remove previous installation
	cmd := "apt-get purge -y postfix dovecot-core dovecot-imapd dovecot-lmtpd dovecot-mysql; "
	cmd += "rm -rf /etc/postfix; rm -rf /etc/dovecot; "
	cmd += "apt-get purge -y spamassassin spamc; "
	cmd += "userdel vmail; userdel spamd; rm -rf /var/mail; rm -rf /home/spamd; "
	cmd += "apt-get purge -y opendkim opendkim-tools; "
	cmd += "rm -rf /etc/opendkim"
	term.RunLongCommandIgnoreError(con, cmd)

	cmd = fmt.Sprintf("sudo hostnamectl set-hostname %s", server.Postfix.Hostname)
	term.RunLongCommand(con, cmd)

	cmd = "/etc/init.d/networking restart"
	term.RunLongCommand(con, cmd)

	// installing myhostname will be overriden in future steps
	cmd = fmt.Sprintf("echo 'postfix postfix/mailname string %s' | debconf-set-selections && ", server.Postfix.Hostname)
	cmd += "echo 'postfix postfix/main_mailer_type string \"Internet Site\"' | debconf-set-selections && "
	cmd += "apt-get install -y postfix postfix-mysql dovecot-core dovecot-imapd dovecot-lmtpd dovecot-mysql"
	term.RunLongCommand(con, cmd)
}

// configure MySQL
func Postfix02(con *ssh.Client, server *config.Server) {
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8;", server.Postfix.MySQLDb)
	sql += fmt.Sprintf("GRANT SELECT ON %s.* TO '%s'@'127.0.0.1' IDENTIFIED BY '%s';", server.Postfix.MySQLDb, server.Postfix.MySQLUser, server.Postfix.MySQLPassword)
	sql += "FLUSH PRIVILEGES;"
	sql += "USE servermail;"

	sql += `
        CREATE TABLE IF NOT EXISTS virtual_domains (
            id  INT NOT NULL AUTO_INCREMENT,
            name VARCHAR(50) NOT NULL,
            PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	sql += `
        CREATE TABLE IF NOT EXISTS virtual_users (
            id INT NOT NULL AUTO_INCREMENT,
            domain_id INT NOT NULL,
            password VARCHAR(106) NOT NULL,
            email VARCHAR(120) NOT NULL,
            PRIMARY KEY (id),
            UNIQUE KEY email (email),
            FOREIGN KEY (domain_id) REFERENCES virtual_domains(id) ON DELETE CASCADE
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	sql += `
        CREATE TABLE IF NOT EXISTS virtual_aliases (
            id INT NOT NULL AUTO_INCREMENT,
            domain_id INT NOT NULL,
            source varchar(100) NOT NULL,
            destination varchar(100) NOT NULL,
            PRIMARY KEY (id),
            FOREIGN KEY (domain_id) REFERENCES virtual_domains(id) ON DELETE CASCADE
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	cmd := fmt.Sprintf("mysql -u root -p%s -h %s -e %q", server.MySQL.User.Password, server.Ip, sql)
	funcs.RunCommandSh(cmd)

	sql = fmt.Sprintf("REPLACE INTO %s.virtual_domains (id ,name) VALUES ", server.Postfix.MySQLDb)
	sql += fmt.Sprintf("('1', '%s');", server.Postfix.Hostname)

	cmd = fmt.Sprintf("mysql -u root -p%s -h %s -e %q", server.MySQL.User.Password, server.Ip, sql)
	funcs.RunCommandSh(cmd)

	for _, email := range server.Postfix.Emails {
		sql = fmt.Sprintf("REPLACE INTO %s.virtual_users (id, domain_id, password, email) VALUES ", server.Postfix.MySQLDb)
		sql += fmt.Sprintf("('1', '1', ENCRYPT('%s', '%s'), '%s');", email.Password, server.Postfix.MySQLSalt, email.Email)

		cmd = fmt.Sprintf("mysql -u root -p%s -h %s -e %q", server.MySQL.User.Password, server.Ip, sql)
		funcs.RunCommandSh(cmd)

		if len(email.Alias) > 0 {
			sql = fmt.Sprintf("REPLACE INTO %s.virtual_aliases (id, domain_id, source, destination) VALUES", server.Postfix.MySQLDb)
			sql += fmt.Sprintf("('1', '1', '%s', '%s');", email.Alias, email.Email)

			cmd = fmt.Sprintf("mysql -u root -p%s -h %s -e %q", server.MySQL.User.Password, server.Ip, sql)
			funcs.RunCommandSh(cmd)
		}

	}

}

func Postfix03(con *ssh.Client, server *config.Server) {
	CopyFileToServer(server, "postfix", "main.cf", "root", "/etc/postfix/main.cf")

	// Add to the end of "main.cf" hostname
	// instead hostname: mail.hostname
	// in other places hostname = domain
	cmd := fmt.Sprintf("echo '\nmyhostname = mail.%s' >> /etc/postfix/main.cf", server.Postfix.Hostname)
	term.RunLongCommand(con, cmd)
}

// /etc/postfix/mysql-virtual-mailbox-domains.cf
func Postfix04(con *ssh.Client, server *config.Server) {
	// create file
	data := fmt.Sprintf("user = %s\n", server.Postfix.MySQLUser)
	data += fmt.Sprintf("password = %s\n", server.Postfix.MySQLPassword)
	data += "hosts = 127.0.0.1\n"
	data += fmt.Sprintf("dbname = %s\n", server.Postfix.MySQLDb)
	data += "query = SELECT 1 FROM virtual_domains WHERE name='%s'\n"

	CopyTempFileToServer(server, data, "root", "/etc/postfix/mysql-virtual-mailbox-domains.cf")
}

// /etc/postfix/mysql-virtual-mailbox-maps.cf
func Postfix05(con *ssh.Client, server *config.Server) {
	// создаем файл
	data := fmt.Sprintf("user = %s\n", server.Postfix.MySQLUser)
	data += fmt.Sprintf("password = %s\n", server.Postfix.MySQLPassword)
	data += "hosts = 127.0.0.1\n"
	data += fmt.Sprintf("dbname = %s\n", server.Postfix.MySQLDb)
	data += "query = SELECT 1 FROM virtual_users WHERE email='%s'\n"

	CopyTempFileToServer(server, data, "root", "/etc/postfix/mysql-virtual-mailbox-maps.cf")
}

// /etc/postfix/mysql-virtual-alias-maps.cf
func Postfix06(con *ssh.Client, server *config.Server) {
	// создаем файл
	data := fmt.Sprintf("user = %s\n", server.Postfix.MySQLUser)
	data += fmt.Sprintf("password = %s\n", server.Postfix.MySQLPassword)
	data += "hosts = 127.0.0.1\n"
	data += fmt.Sprintf("dbname = %s\n", server.Postfix.MySQLDb)
	data += "query = SELECT 1 FROM virtual_aliases WHERE source='%s'\n"

	CopyTempFileToServer(server, data, "root", "/etc/postfix/mysql-virtual-alias-maps.cf")
}

func Postfix07(con *ssh.Client, server *config.Server) {
	CopyFileToServer(server, "postfix", "master.cf", "root", "/etc/postfix/master.cf")

	CopyFileToServer(server, "postfix", "dovecot.conf", "root", "/etc/dovecot/dovecot.conf")
	CopyFileToServer(server, "postfix", "10-mail.conf", "root", "/etc/dovecot/conf.d/10-mail.conf")
	CopyFileToServer(server, "postfix", "10-auth.conf", "root", "/etc/dovecot/conf.d/10-auth.conf")
	CopyFileToServer(server, "postfix", "auth-sql.conf.ext", "root", "/etc/dovecot/conf.d/auth-sql.conf.ext")
	CopyFileToServer(server, "postfix", "10-master.conf", "root", "/etc/dovecot/conf.d/10-master.conf")
	CopyFileToServer(server, "postfix", "10-ssl.conf", "root", "/etc/dovecot/conf.d/10-ssl.conf")
}

func Postfix08(con *ssh.Client, server *config.Server) {
	cmd := fmt.Sprintf("mkdir -p /var/mail/vhosts/%s", server.Postfix.Hostname)
	term.RunLongCommand(con, cmd)

	term.RunLongCommand(con, "groupadd -g 5000 vmail")
	term.RunLongCommand(con, "useradd -g vmail -u 5000 vmail -d /var/mail")
	term.RunLongCommand(con, "chown -R vmail:vmail /var/mail")
}

// /etc/dovecot/dovecot-sql.conf.ext
func Postfix09(con *ssh.Client, server *config.Server) {
	data := "driver = mysql\n"
	data += fmt.Sprintf("connect = host=127.0.0.1 dbname=%s user=%s password=%s\n", server.Postfix.MySQLDb, server.Postfix.MySQLUser, server.Postfix.MySQLPassword)
	data += "default_pass_scheme = SHA512-CRYPT\n"
	data += "password_query = SELECT email as user, password FROM virtual_users WHERE email='%u';\n"

	CopyTempFileToServer(server, data, "root", "/etc/dovecot/dovecot-sql.conf.ext")
}

// self-sertificate for Dovecot
func Postfix10(con *ssh.Client, server *config.Server) {
	cmd := "openssl req -new -x509 -days 1500 -nodes -out \"/etc/ssl/certs/dovecot.pem\" -keyout \"/etc/ssl/private/dovecot.pem\" "
	cmd += "-subj \"/C=RU/ST=./L=./O=./OU=./CN=.\" "
	term.RunLongCommand(con, cmd)
	term.RunLongCommand(con, "chown -R vmail:dovecot /etc/dovecot && chmod -R o-rwx /etc/dovecot")
}

// spamassassin and it's configs
func Postfix11(con *ssh.Client, server *config.Server) {
	term.RunLongCommand(con, "apt-get install -y spamassassin spamc")
	term.RunLongCommand(con, "useradd spamd -d /home/spamd/ -m -s /bin/false")

	CopyFileToServer(server, "postfix", "spamassassin", "root", "/etc/default/spamassassin")
	CopyFileToServer(server, "postfix", "local.cf", "root", "/etc/spamassassin/local.cf")
}

// opendkim
func Postfix12(con *ssh.Client, server *config.Server) {
	term.RunLongCommand(con, "apt-get install -y opendkim opendkim-tools")
	term.RunLongCommand(con, "mkdir -p /etc/opendkim && mkdir -p /etc/opendkim/keys")

	CopyFileToServer(server, "postfix", "opendkim.conf", "root", "/etc/opendkim.conf")
	CopyFileToServer(server, "postfix", "opendkim", "root", "/etc/default/opendkim")
	CopyFileToServer(server, "postfix", "TrustedHosts", "root", "/etc/opendkim/TrustedHosts")
}

// /etc/opendkim/KeyTable and /etc/opendkim/SigningTable
func Postfix13(con *ssh.Client, server *config.Server) {
	s := fmt.Sprintf("mail._domainkey.%s %s:mail:/etc/opendkim/keys/%s/mail.private", server.Postfix.Hostname, server.Postfix.Hostname, server.Postfix.Hostname)
	cmd := fmt.Sprintf("echo '%s' > /etc/opendkim/KeyTable", s)
	term.RunLongCommand(con, cmd)

	s = fmt.Sprintf("*@%s mail._domainkey.%s", server.Postfix.Hostname, server.Postfix.Hostname)
	cmd = fmt.Sprintf("echo '%s' > /etc/opendkim/SigningTable", s)
	term.RunLongCommand(con, cmd)

	cmd = fmt.Sprintf("mkdir -p /etc/opendkim/keys/%s", server.Postfix.Hostname)
	term.RunLongCommand(con, cmd)

	cmd = fmt.Sprintf("opendkim-genkey -s /etc/opendkim/keys/%s/mail -d /etc/opendkim/keys/%s", server.Postfix.Hostname, server.Postfix.Hostname)
	term.RunLongCommand(con, cmd)

	cmd = fmt.Sprintf("chown opendkim:opendkim /etc/opendkim/keys/%s/mail.private", server.Postfix.Hostname)
	term.RunLongCommand(con, cmd)
}

func Postfix14(con *ssh.Client, server *config.Server) {
	term.RunLongCommand(con, "systemctl restart spamassassin")
	term.RunLongCommand(con, "systemctl restart postfix")
	term.RunLongCommand(con, "systemctl restart dovecot")
	term.RunLongCommand(con, "systemctl restart opendkim")

}
