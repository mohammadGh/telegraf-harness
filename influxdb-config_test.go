package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMapToTomlString(t *testing.T) {
	var input map[string]interface{}
	input = make(map[string]interface{})
	input["id"] = 12
	input["interval"] = "10s"
	input["database"] = "telegraf"
	input["urls"] = [3]string{"A", "B", "C"}
	input["nums"] = [3]int{1, 2, 3}

	str := MapToTomlString(input)

	fmt.Println(str)

	assert.True(t, strings.Contains(str, "nums=[1,2,3]"))
	assert.True(t, strings.Contains(str, "urls=[\"A\",\"B\",\"C\"]"))
	assert.True(t, strings.Contains(str, "interval=\"10s\""))
}

func TestConfigTemplate(t *testing.T) {
	var input map[string]string
	input = make(map[string]string)
	input["interval"] = "8s"
	input["url"] = "http://195.146.59.40:8086"
	input["database"] = "myMetrics"

	str := fillConfigTemplate(input, "lib/config-template.conf")

	//fmt.Println(str)

	assert.True(t, strings.Contains(str, `interval = "8s"`))
	assert.True(t, strings.Contains(str, `urls = ["http://195.146.59.40:8086"]`))
	assert.True(t, strings.Contains(str, `database = "myMetrics"`))
}

func TestConfigTemplate2(t *testing.T) {
	var input map[string]string
	input = make(map[string]string)
	//input["interval"]="8s"
	//input["url"]="http://195.146.59.40:8086"
	//input["database"]="myMetrics"

	str := fillConfigTemplate(input, "lib/config-template.conf")

	//fmt.Println(str)

	assert.True(t, strings.Contains(str, `interval = "3s"`))
	assert.True(t, strings.Contains(str, `urls = ["http://127.0.0.1:8086"]`))
	assert.True(t, strings.Contains(str, `database = "metrics"`))
}

func TestDecodeArray(t *testing.T) {
	decodeArray("[1,2]")
}
