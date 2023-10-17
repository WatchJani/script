package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

//get the go version
//if not exist install
//if exist check this version

const (
	DOWNLOAD_DIRECTORY = "/usr/local" //
	BASH_PATH          = "/home/janko/.bashrc"
	LINE_TO_ADD        = "export PATH=$PATH:/usr/local/go/bin"
)

func main() {
	//get link from arg
	arg, err := ArgParser()
	if err != nil {
		log.Println(err)
		return
	}

	//check link
	if isValidURL(arg) {
		log.Println("Please check argument, need to be path to Golang source")
		return
	}

	// //check if exist go on your machine
	// version, err := CurrentVersion()

	// if err != nil {
	// 	if !strings.Contains(err.Error(), "Error executing 'go version'") {
	// 		fmt.Printf("Error while obtaining the Go language version: %v\n", err)
	// 		return
	// 	}
	// }

	// if version == -1 {
	// 	return
	// }

	//popraviti greske
	//install
	err = Install(arg)
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Error while obtaining the Go language version: %v\n", err)
			return
		}

		fmt.Printf("Error while obtaining the Go language version: %v\n", err)
	}

	//get installed file name
	file := GetFile(arg)

	//extract
	err = Extract(file)
	if err != nil {
		log.Println(err)
		return
	}

	//delate compress dir
	err = DeleteDir(file)
	if err != nil {
		log.Println(err)
		return
	}

	//update file configuration
	if !IsLinePresent() {
		err = AddLineToBashrc()
		if err != nil {
			log.Println(err)
			return
		}
	}

	//reload file configuration
	err = ReloadBash()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("GO IS INSTALLED!!!")
}

func IsLinePresent() bool {
	file, err := os.Open(BASH_PATH)
	if err != nil {
		return false
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == LINE_TO_ADD {
			return true
		}
	}

	return false
}

func ReloadBash() error {
	cmd := exec.Command("bash", "-c", "source "+BASH_PATH)

	return cmd.Run()
}

func AddLineToBashrc() error {
	file, err := os.OpenFile(BASH_PATH, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(LINE_TO_ADD + "\n")
	return err
}

func GetFile(path string) string {
	name := strings.Split(path, "/")

	return name[len(name)-1]
}

func Extract(dir string) error {
	cmd := exec.Command("tar", "-xzvf", DOWNLOAD_DIRECTORY+dir)

	cmd.Dir = DOWNLOAD_DIRECTORY

	return cmd.Run()
}

func DeleteDir(dir string) error {
	return os.Remove(DOWNLOAD_DIRECTORY + dir)
}

func ArgParser() (string, error) {
	args := os.Args[1:]
	if len(args) != 1 {
		return "", fmt.Errorf("This command requires exactly one argument")
	}
	return args[0], nil
}

func Install(path string) error {
	cmd := exec.Command("wget", path)

	cmd.Dir = DOWNLOAD_DIRECTORY

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func isValidURL(url string) bool {
	re := regexp.MustCompile(`^https://go\.dev/dl/go\d+\.\d+\.\d+\.linux-amd64\.tar\.gz$`)
	return !re.MatchString(url)
}

func CurrentVersion() (int, error) {
	cmd := exec.Command("go", "version")

	output, err := cmd.Output()
	if err != nil {
		return -1, fmt.Errorf("Error executing 'go version': %v", err)
	}

	versionString := string(output)
	versionParts := strings.Fields(versionString)

	if len(versionParts) < 3 {
		return -1, fmt.Errorf("Unable to extract the Go language version")
	}

	goVersion := versionParts[2]

	return Remove(strings.TrimPrefix(goVersion, "go"))
}

func Remove(version string) (int, error) {
	version = strings.ReplaceAll(version, ".", "")

	return IntConverter(version)
}

func IntConverter(value string) (int, error) {
	num, err := strconv.Atoi(value)

	if err != nil {
		return -1, fmt.Errorf("Error converting to an integer: %v", err)
	}
	return num, nil
}
