package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func writeInfluxConfigFile(txt string, path string) {
	var sampleConfig = `
[[outputs.influxdb]]
	urls = ["http://127.0.0.1:8086"] 
	database = "telegraf"

[[inputs.win_perf_counters.object]]
	ObjectName = "Processor"
	Instances = ["*"]
	Counters = ["% Idle Time"]
	Measurement = "cpu"
	IncludeTotal=true
`
	bytes := []byte(sampleConfig)
	// write the whole body at once
	err := ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		panic(err)
	}
}
func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func writeConfigFile(txt string, path string, configMap map[string]interface{}) error {
	// parse key:values to valid influxdb config statements

	for key, value := range configMap {

	}

	return nil
}
func MapToString(m map[string]string) string {
	bytes := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(bytes, "%s=\"%s\"\n", key, value)
	}
	return bytes.String()
}
