package main

import "io/ioutil"

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
