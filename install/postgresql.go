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

	// TODO create password
}
