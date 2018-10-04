package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func StartJolokia(jolokiaPath string, jolokiaOption string, javaProcess string) (string, error) {

	var params = "-jar " + jolokiaPath + " " + jolokiaOption + " " + javaProcess
	java := filepath.Join(getJavaPath(), "/bin/java.exe")
	jolokiaCommand := exec.Command(java, strings.Split(params, " ")...)

	out, err := jolokiaCommand.Output()
	if err != nil {
		return "", err
	}
	outstr := fmt.Sprintf("%s", out)
	return outstr, nil
}

func getJavaPath() string {
	return os.Getenv("JAVA_HOME")
}

func getJavaVersion() string {
	java := filepath.Join(getJavaPath(), "/bin/java.exe")
	fmt.Println(java)
	javaVersion := exec.Command(java, strings.Split("-version", " ")...)
	out, err := javaVersion.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	outstr := fmt.Sprintf("%s", out)
	return outstr
}
