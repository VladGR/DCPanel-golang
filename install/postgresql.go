package install

import (
	"djcontrol/config"
	"djcontrol/term"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func PostgreSQL(con *ssh.Client, server *config.Server) {
	fmt.Println("Installing PostgeSQL...")

	cmd := fmt.Sprintf("apt-get -y install postgresql-%s && ", server.PostgreSQL.Version)
	cmd += "apt-get -y install pgadmin3 && "
	cmd += "apt-get -y install postgresql-contrib && "
	cmd += "systemctl enable postgresql && "
	cmd += "systemctl start postgresql"
	term.RunLongCommand(con, cmd)

	pgHBA := fmt.Sprintf("/etc/postgresql/%s/main/pg_hba.conf", server.PostgreSQL.Version)
	CopyFileToServer(server, "postgresql", "pg_hba.conf", "root", pgHBA)

	cmd = fmt.Sprintf("chown postgres: %s", pgHBA)
	term.RunLongCommand(con, cmd)

	cmd = "cd /tmp && "
	cmd += fmt.Sprintf("cd /tmp && sudo -u postgres psql -U postgres -d postgres -c \"alter user postgres with password '%s';\"", server.PostgreSQL.User.Password)
	cmd += " && cd /root "
	term.RunLongCommand(con, cmd)

	// Access from any IPs
	cmd = fmt.Sprintf("echo \"listen_addresses = '*'\" >> /etc/postgresql/%s/main/postgresql.conf", server.PostgreSQL.Version)
	term.RunLongCommand(con, cmd)

	term.RunLongCommand(con, "systemctl restart postgresql")

}
