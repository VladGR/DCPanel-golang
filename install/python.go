package install

import (
	"djcontrol/term"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func Python2(con *ssh.Client) {
	fmt.Println("Installing Python2...")
	cmd := "apt-get -y install python python-dev python-virtualenv python-setuptools python-pip python-psycopg2 && "
	cmd += "pip install pycurl && "
	cmd += "pip install django && "
	cmd += "pip install pillow && "
	cmd += "pip install mysql-python && "
	cmd += "pip install xlrd && "
	cmd += "pip install xlwt && "
	cmd += "pip install beautifulsoup && "
	cmd += "pip install beautifulsoup4 && "
	// cmd += "pip install lxml && "
	cmd += "pip install redis && "
	cmd += "pip install reportlab && "
	cmd += "pip install pycrypto"
	term.RunLongCommand(con, cmd)

}

func Python3(con *ssh.Client) {
	fmt.Println("Installing Python3...")
	cmd := "apt-get -y install python3 python3-dev python-virtualenv python3-setuptools python3-pip python3-psycopg2 && "
	cmd += "pip3 install pycurl && "
	cmd += "pip3 install django  && "
	cmd += "pip3 install pillow && "
	cmd += "pip3 install mysqlclient && "
	cmd += "pip3 install xlrd && "
	cmd += "pip3 install xlwt-future && "
	cmd += "pip3 install beautifulsoup4 && "
	// cmd += "pip3 install lxml && "
	cmd += "pip3 install redis && "
	cmd += "pip3 install reportlab && "
	cmd += "pip3 install pycrypto"
	term.RunLongCommand(con, cmd)

}
