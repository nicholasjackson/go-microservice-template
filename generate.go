package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var templateRegex = regexp.MustCompile("^.*\\.tmpl")
var gitRegex = regexp.MustCompile("^.git.*")

var scanner = bufio.NewScanner(os.Stdin)

type templateData struct {
	ServiceName string
	Namespace   string
	StatsD      bool
}

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

func includeStatsD(scanner *bufio.Scanner) bool {
	fmt.Printf("Include StatsD? (y|n)\n")
	scanner.Scan()
	line := scanner.Text()
	fmt.Println("")

	return line == "y"
}

func main() {
	printHeader()

	nameSpace := requestNamespace(scanner)
	serviceName := requestName(scanner)
	statsD := includeStatsD(scanner)

	if confirm(serviceName, nameSpace, scanner) {
		data := templateData{ServiceName: serviceName, Namespace: nameSpace, StatsD: statsD}
		generateTemplate(data)
	} else {
		fmt.Println("")
		fmt.Println("Fine I won't")
	}
}

func printHeader() {
	header := `
	___  ____                                    _
	|  \/  (_)                                  (_)
	| .  . |_  ___ _ __ ___  ___  ___ _ ____   ___  ___ ___
	| |\/| | |/ __| '__/ _ \/ __|/ _ \ '__\ \ / / |/ __/ _ \
	| |  | | | (__| | | (_) \__ \  __/ |   \ V /| | (_|  __/
	\_|  |_/_|\___|_|  \___/|___/\___|_|    \_/ |_|\___\___|
	 _____                    _       _
	|_   _|                  | |     | |
	  | | ___ _ __ ___  _ __ | | __ _| |_ ___
	  | |/ _ \ '_ ' _ \| '_ \| |/ _' | __/ _ \
	  | |  __/ | | | | | |_) | | (_| | ||  __/
	  \_/\___|_| |_| |_| .__/|_|\__,_|\__\___|
	                   | |
	                   |_|
	`
	fmt.Println(header)
	fmt.Println("")
}

func confirm(serviceName, nameSpace string, scanner *bufio.Scanner) bool {
	fmt.Printf("Generating Microservice template: %s/%s in GOPATH\n", nameSpace, serviceName)
	fmt.Printf("Is this correct? (y|n)\n")
	scanner.Scan()
	line := scanner.Text()

	return line == "y"
}

func generateTemplate(data templateData) {
	fmt.Println("output path: ", os.Getenv("GOPATH"))
	destination := destinationFolder(data.ServiceName, data.Namespace)
	os.MkdirAll(destination, os.ModePerm)

	copyNonGitFiles(destination, data)
	//renameInFiles(serviceName, nameSpace)
}

func copyNonGitFiles(destination string, data templateData) {
	err := filepath.Walk("./template_files", func(path string, f os.FileInfo, err error) error {
		_ = processNonGitFile(path, destination, data)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func processNonGitFile(path string, destination string, data templateData) error {
	newFile := replaceDefaultNameInPath(path, data.ServiceName)
	newFile = replaceTemplateExtInPath(newFile)
	destinationFile := destination + "/" + newFile

	switch {
	case gitRegex.MatchString(path):
		fmt.Println("Skipping git path: ", path)
	case templateRegex.MatchString(path):
		err := saveAndProcessTemplate(path, destinationFile, data)
		if err != nil {
			fmt.Println("Unable to process template:", err)
			return err
		}
	default:
		copyFile(path, destinationFile)
	}

	return nil
}

func saveAndProcessTemplate(src string, dst string, data templateData) error {
	fmt.Println("Process template: ", src)

	f, err := ioutil.ReadFile(src)
	templateString := string(f[:])

	tmpl, err := template.New("Template").Parse(templateString)
	if err != nil {
		return err
	}

	createFolder(dst)

	output, err := os.Create(dst)
	defer output.Close()

	return tmpl.Execute(output, data)
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer source.Close()

	sourceInfo, err := source.Stat()
	if sourceInfo.Mode().IsRegular() {
		fmt.Println("Copying path: ", src)
		createFolder(dst)

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

func createFolder(dst string) {
	dstDir := filepath.Dir(dst)
	os.MkdirAll(dstDir, os.ModePerm)
}

func replaceDefaultNameInPath(path string, serviceName string) string {
	newpath := strings.Replace(path, "microservice-template", serviceName, -1)
	newpath = strings.Replace(newpath, "template_files/", "", -1)

	return newpath
}

func replaceTemplateExtInPath(path string) string {
	return strings.Replace(path, ".tmpl", "", -1)
}

func destinationFolder(serviceName string, nameSpace string) string {
	return fmt.Sprintf("%s/src/%s/%s", os.Getenv("GOPATH"), nameSpace, serviceName)
}
