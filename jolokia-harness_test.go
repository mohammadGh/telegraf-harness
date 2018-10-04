package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestRunJolokia(t *testing.T) {
	StartJolokia("lib/jolokia-jvm-1.6.0-agent.jar", "stop", "4916")
}

func TestGetJavaHome(t *testing.T) {
	str := getJavaPath()
	fmt.Println(str)
	//assert.True(t, strings.Contains(str, "Java"))
}

func TestGetJavaVersion(t *testing.T) {
	str := getJavaVersion()
	fmt.Println(str)
	//assert.True(t, strings.Contains(str, "Java"))
}

func TestSplit(t *testing.T) {
	str := strings.Split("-version", " ")
	fmt.Println(len(str))
	fmt.Println(str[0])
	//assert.True(t, strings.Contains(str, "Java"))
}
