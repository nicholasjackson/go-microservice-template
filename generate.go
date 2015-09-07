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
)

var scanner = bufio.NewScanner(os.Stdin)
var serviceName = requestName(scanner)
var destination = requestDestination(scanner)

func main() {
    if confirm(serviceName, destination, scanner) {
      generateTemplate(serviceName, destination)
    } else {
      fmt.Println("Fine I won't")
    }
}

func requestName(scanner *bufio.Scanner) string {
    fmt.Println("What is the name of this microservice?")
    scanner.Scan()
    line := scanner.Text()
    return line
}

func requestDestination(scanner *bufio.Scanner) string {
    fmt.Println("Where shall I save the template?")
    scanner.Scan()
    line := scanner.Text()
    return line
}

func confirm(serviceName, destination string, scanner *bufio.Scanner) bool {
    fmt.Printf("Generating Microservice template: %s in %s\n", serviceName, destination)
    fmt.Printf("Is this correct? (y|n)\n")
    scanner.Scan()
    line := scanner.Text()

    return line == "y"
}

func generateTemplate(serviceName, destination string) {
    copyNonGitFiles(serviceName, destination)
    renameInFiles(serviceName, destination)
}

func copyNonGitFiles(serviceName, destination string) {
    os.MkdirAll(destination, os.ModePerm)
    err := filepath.Walk(".", copyNonGitFile)
    if err != nil {
        panic(err)
    }
    return
}

func copyNonGitFile(path string, f os.FileInfo, err error) error {
    gitRegex := regexp.MustCompile("^.git.*")
    if gitRegex.MatchString(path) {
        fmt.Println("Skipping git path: ", path)
    } else {
        newFile := replaceMicroserviceTemplate(path)
        destinationFile := destination + "/" + newFile

        copyFile(path, destinationFile)
    }
    return nil
}

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
        fmt.Println("Copying path: ", src)
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

func replaceMicroserviceTemplate(string string) string {
    return strings.Replace(string, "microservice-template", serviceName, -1)
}
