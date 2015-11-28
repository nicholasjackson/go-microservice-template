package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

func requestName(scanner *bufio.Scanner) string {
	fmt.Println("What is the name of this microservice?")
	scanner.Scan()
	line := scanner.Text()

	fmt.Println("")
	return line
}

func requestNamespace(scanner *bufio.Scanner) string {
	fmt.Println("What is the namespace for this microservice e.g. github.com/nicholasjackson?")
	scanner.Scan()
	line := scanner.Text()

	fmt.Println("")
	return line
}

func main() {
	printHeader()

	nameSpace := requestNamespace(scanner)
	serviceName := requestName(scanner)

	if confirm(serviceName, nameSpace, scanner) {
		generateTemplate(serviceName, nameSpace)
	} else {
		fmt.Println("")
		fmt.Println("Fine I won't")
	}
}

func printHeader() {
	fmt.Println("___  ____                                    _          ")
	fmt.Println("|  \\/  (_)                                  (_)         ")
	fmt.Println("| .  . |_  ___ _ __ ___  ___  ___ _ ____   ___  ___ ___ ")
	fmt.Println("| |\\/| | |/ __| '__/ _ \\/ __|/ _ \\ '__\\ \\ / / |/ __/ _ \\")
	fmt.Println("| |  | | | (__| | | (_) \\__ \\  __/ |   \\ V /| | (_|  __/")
	fmt.Println("\\_|  |_/_|\\___|_|  \\___/|___/\\___|_|    \\_/ |_|\\___\\___|")
	fmt.Println(" _____                    _       _                     ")
	fmt.Println("|_   _|                  | |     | |                    ")
	fmt.Println("  | | ___ _ __ ___  _ __ | | __ _| |_ ___               ")
	fmt.Println("  | |/ _ \\ '_ ` _ \\| '_ \\| |/ _` | __/ _ \\              ")
	fmt.Println("  | |  __/ | | | | | |_) | | (_| | ||  __/              ")
	fmt.Println("  \\_/\\___|_| |_| |_| .__/|_|\\__,_|\\__\\___|              ")
	fmt.Println("                   | |                                  ")
	fmt.Println("                   |_|                                  ")
	fmt.Println("")
	fmt.Println("")
}

func confirm(serviceName, nameSpace string, scanner *bufio.Scanner) bool {
	fmt.Printf("Generating Microservice template: %s/%s in GOPATH\n", serviceName, nameSpace)
	fmt.Printf("Is this correct? (y|n)\n")
	scanner.Scan()
	line := scanner.Text()

	return line == "y"
}

func generateTemplate(serviceName string, nameSpace string) {
	fmt.Println("output path: ", os.Getenv("GOPATH"))
	destination := destinationFolder(serviceName, nameSpace)
	os.MkdirAll(destination, os.ModePerm)

	copyNonGitFiles(destination, serviceName)
	//renameInFiles(serviceName, nameSpace)
}

func copyNonGitFiles(destination string, serviceName string) {
	err := filepath.Walk("./template_files", func(path string, f os.FileInfo, err error) error {
		_ = copyNonGitFile(path, destination, serviceName)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func copyNonGitFile(path string, destination string, serviceName string) error {
	gitRegex := regexp.MustCompile("^.git.*")
	if gitRegex.MatchString(path) {
		fmt.Println("Skipping git path: ", path)
	} else {
		newFile := replaceDefaultNameInPath(path, serviceName)
		destinationFile := destination + "/" + newFile

		copyFile(path, destinationFile)
	}
	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer source.Close()
	dstDir := filepath.Dir(dst)
	sourceInfo, err := source.Stat()
	if sourceInfo.Mode().IsRegular() {
		fmt.Println("Copying path: ", src, " to: ", dst)
		os.MkdirAll(dstDir, os.ModePerm)
		destination, err := os.Create(dst)
		if err != nil {
			return err
		}
		if _, err := io.Copy(destination, source); err != nil {
			destination.Close()
			return err
		}
		return destination.Close()
	}
	return nil
}

/*
func renameInFiles(serviceName, destination string) {
	filesToEdit := [6]string{"go/src/github.com/nicholasjackson/microservice-template/server.go", "dockercompose/microservice-template/docker-compose.yml", "dockerfile/microservice-template/Dockerfile", "dockerfile/microservice-template/supervisord.conf", ".ruby-gemset", "Rakefile"}
	os.Chdir(destination)
	for _, path := range filesToEdit {
		filename := replaceMicroserviceTemplate(path)
		body, err := ioutil.ReadFile(filename)
		bodyString := string(body[:])
		if err != nil {
			panic(err)
		}
		newBodyString := replaceMicroserviceTemplate(bodyString)
		newBody := []byte(newBodyString)
		err = ioutil.WriteFile(filename, newBody, 0644)
		if err != nil {
			panic(err)
		}
	}
}


*/

func replaceDefaultNameInPath(path string, serviceName string) string {
	return strings.Replace(path, "microservice-template", serviceName, -1)
}

func destinationFolder(serviceName string, nameSpace string) string {
	return fmt.Sprintf("%s/src/%s/%s", os.Getenv("GOPATH"), nameSpace, serviceName)
}
