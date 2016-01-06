package install

import (
	"djcontrol/config"
	"djcontrol/term"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func IPTables(con *ssh.Client, server *config.Server) {
	fmt.Println("Installing IPTables...")
	cmd := "apt-get -y install iptables"
	term.RunLongCommand(con, cmd)

	term.RunLongCommand(con, "mkdir -p /etc/iptables/")

	CopyFileToServer(server, "iptables", "myiptables", "root", "/etc/iptables/myiptables")
	CopyFileToServer(server, "iptables", "myiptables-stop", "root", "/etc/iptables/myiptables-stop")

	CopyFileToServer(server, "iptables", "myip6tables", "root", "/etc/iptables/myip6tables")
	CopyFileToServer(server, "iptables", "myip6tables-stop", "root", "/etc/iptables/myip6tables-stop")

	CopyFileToServer(server, "iptables", "iptables.service", "root", "/usr/lib/systemd/system/iptables.service")
	CopyFileToServer(server, "iptables", "ip6tables.service", "root", "/usr/lib/systemd/system/ip6tables.service")

	term.RunLongCommand(con, "chmod +x /etc/iptables/myiptables")
	term.RunLongCommand(con, "chmod +x /etc/iptables/myiptables-stop")

	term.RunLongCommand(con, "chmod +x /etc/iptables/myip6tables")
	term.RunLongCommand(con, "chmod +x /etc/iptables/myip6tables-stop")

	term.RunLongCommand(con, "systemctl enable iptables")
	term.RunLongCommand(con, "systemctl start iptables")

	term.RunLongCommand(con, "systemctl enable ip6tables")
	term.RunLongCommand(con, "systemctl start ip6tables")

	fmt.Println("\n\n")
	fmt.Println("IPTables status:\n")
	term.RunLongCommand(con, "iptables --list-rules")

	fmt.Println("\n\n")
	fmt.Println("IP6Tables status:\n")
	term.RunLongCommand(con, "ip6tables --list-rules")

	// Delay to have time to view IPTables rules
	fmt.Println("\n\n")
	fmt.Println("Installation will be continued in 10 seconds...")
	time.Sleep(10 * time.Second)

}
