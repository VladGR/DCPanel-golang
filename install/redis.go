package install

import (
	"djcontrol/term"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func Redis(con *ssh.Client) {
	fmt.Println("Installing Redis...")
	cmd := "apt-get -y install redis-server && "
	cmd += "systemctl enable redis-server && "
	cmd += "systemctl start redis-server"

	term.RunLongCommand(con, cmd)
}
