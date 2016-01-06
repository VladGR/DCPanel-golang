package funcs

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// fast command
func RunCommandSh(s string) {

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("sh", "-c", s)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		fmt.Println(out.String())
		log.Fatal(err.Error())
	}

	fmt.Println(out.String())
}

// fast command ignores errors
func RunCommandShIgnoreError(s string) {

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("sh", "-c", s)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
	}
	fmt.Println(out.String())
}

// fast command
func RunCommand(s string) {

	var out bytes.Buffer
	var stderr bytes.Buffer

	parts := strings.Fields(s)
	if len(parts) == 1 {
		cmd := exec.Command(s)
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(stderr.String())
			log.Fatal(err)
		}
	} else {
		head := parts[0]
		parts = parts[1:len(parts)]

		cmd := exec.Command(head, parts...)

		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println(stderr.String())
			log.Fatal(err)
		}
	}

	fmt.Println(out.String())
}

// long running real-time command
func RunLongCommand(s string) {
	parts := strings.Fields(s)
	head := parts[0]
	parts = parts[1:len(parts)]

	cmd := exec.Command(head, parts...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Ошибка при создании StdoutPipe")
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal("Ошибка запуска команды")
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка при ожидании результата", err)
		os.Exit(1)
	}
}

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func IsSliceContainsString(slice []string, s string) bool {
	for _, x := range slice {
		if s == x {
			return true
		}
	}
	return false
}

func CheckErr(err error) {
	if err != nil {
		log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
		log.Fatal(err)
	}
}
