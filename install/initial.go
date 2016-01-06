package install

import (
	"djcontrol/config"
	"djcontrol/funcs"
	"djcontrol/term"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// Main User
func CreateUser(con *ssh.Client, server *config.Server, name string) {

	fmt.Printf("Creating user %q ...\n", name)

	term.RunLongCommand(con, fmt.Sprintf("mkdir -p /home/%s", name))

	// ignore error if user exists
	cmd := fmt.Sprintf("id -u %s &>/dev/null || useradd %s -d /home/%s", name, name, name)
	term.RunLongCommand(con, cmd)
	cmd = fmt.Sprintf("chown %s: /home/%s", name, name)
	term.RunLongCommand(con, cmd)

	// create password
	newPassword := funcs.RandomString(10)
	cmd = fmt.Sprintf("echo %s:%s | sudo chpasswd", name, newPassword)
	term.RunLongCommand(con, cmd)

	// create ssh keys
	cmd = fmt.Sprintf("sshpass -p %s ssh-copy-id -i ~/.ssh/id_rsa.pub %s@%s", newPassword, name, server.Ip)
	// ignore error if keys have been set earlier
	funcs.RunCommandShIgnoreError(cmd)
}

// Base initial install
func Base(con *ssh.Client) {
	fmt.Println("Installing base programs...")
	cmd := `
        apt-get update && \
        apt-get -y install \
        sudo build-essential libevent-dev libjpeg-dev nano htop \
        curl libcurl4-gnutls-dev libgnutls28-dev libghc-gnutls-dev libmysqlclient-dev \
        libxml2-dev libxslt1-dev libpq-dev \
        git-core debconf-utils
    `
	term.RunLongCommand(con, cmd)
	term.RunLongCommand(con, "mkdir -p /usr/lib/systemd/system/")
}

// create file and copy to remote server
// create symbolic link
func DropCache(con *ssh.Client, server *config.Server) {
	s := "#!/bin/bash\n"
	s += "sync; echo 3 > /proc/sys/vm/drop_caches"

	localFile := "dropcache"

	ioutil.WriteFile(localFile, []byte(s), 0644)
	cmd := fmt.Sprintf("scp -v %s root@%s:/root/dropcache", localFile, server.Ip)
	funcs.RunCommand(cmd)
	term.RunLongCommand(con, "chmod u+x /root/dropcache")
	term.RunLongCommand(con, "ln -sf /root/dropcache /usr/bin/dropcache")
	os.Remove(localFile)
}
