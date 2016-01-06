package main

import (
	"djcontrol/django"
	"djcontrol/funcs"
	"djcontrol/install"
	"fmt"
)

func main() {
	start()
}

func start() {
	funcs.RunCommand("clear")

	fmt.Println("Choose action:")
	fmt.Println("1: Server Installation")
	fmt.Println("2: Deploy Django Project")

	var input string
	fmt.Scanln(&input)

	switch input {
	case "1":
		install.Start()
		return
	case "2":
		django.Deploy()
		return
	default:
		fmt.Println("Incorrect number.")
	}
}
