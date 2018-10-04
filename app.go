package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	influxjson "github.com/mohammadGh/influxdb-line-protocol-to-json"
	"github.com/tevino/abool"
)

var influxCmd *exec.Cmd
var influxRunningMutexFlag *abool.AtomicBool = abool.New()
var javaProcess string = ""
var Info *log.Logger

func main() {
	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime)

	r := mux.NewRouter()

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if len(params) == 0 {
			params = r.Form
		}
		formAsKeyMaps := httpFormToKeyMap(params)

		writeInfluxConfigFile(formAsKeyMaps, "lib/config-template.conf", "lib/config.conf")
		var telegrafArg = "-test -config lib/config.conf"
		cmd := exec.Command("lib/t.exe", strings.Split(telegrafArg, " ")...)
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		outstr := fmt.Sprintf("%s", out)
		outstr = strings.Replace(outstr, "> ", "", -1)

		jsonStr := influxjson.LinesToJson(outstr)

		// info telegraf test output
		outputLines := strings.Split(outstr, "\n")
		numLines := strconv.Itoa(len(outputLines))
		Info.Println("Telegraf test command extracted " + numLines + " metrics")

		respondWithJSONStr(w, 200, jsonStr)
	})

	r.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		var result map[string]string = make(map[string]string)
		if influxRunningMutexFlag.SetToIf(false, true) {

			params := r.URL.Query()
			if len(params) == 0 {
				params = r.Form
			}
			formAsKeyMaps := httpFormToKeyMap(params)

			//start jolokia
			if value, ok := formAsKeyMaps["java"]; ok {
				javaProcess = value
				_, err := StartJolokia("lib/jolokia-jvm-1.6.0-agent.jar", "start", javaProcess)
				if err != nil {
					result["extra_agents"] = "Java(Jolokia) on process " + javaProcess
					Info.Println("Jolokia agent started on java process " + javaProcess)
				} else {
					Info.Println("Error in starting Jolokia agent on java process " + javaProcess)
				}
			}

			writeInfluxConfigFile(formAsKeyMaps, "lib/config-template.conf", "lib/config.conf")
			var telegrafArg = "-config lib/config.conf"
			influxCmd = exec.Command("lib/t.exe", strings.Split(telegrafArg, " ")...)
			err := influxCmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			Info.Println("Telegraf is Started")
			result["status"] = "started"
			respondWithJSON(w, 200, result)
			return
		}
		Info.Println("Telegraf is already Started")
		respondWithJSON(w, 200, map[string]string{"status": "already started"})
	})

	r.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if influxRunningMutexFlag.SetToIf(true, false) {
			//Stop jolokia
			if javaProcess != "" {
				StartJolokia("lib/jolokia-jvm-1.6.0-agent.jar", "stop", javaProcess)
			}

			err := influxCmd.Process.Kill()
			if err != nil {
				log.Fatal(err)
			}
			influxCmd = nil
			javaProcess = ""
			Info.Println("Telegraf is Stopped")
			respondWithJSON(w, 200, map[string]string{"status": "stopped"})
			return
		}

		Info.Println("Telegraf is already Stopped")
		respondWithJSON(w, 200, map[string]string{"status": "already stopped"})
	})

	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		//it should allow to upload a custom configuaratio template file
		//TODO: mgh, 2017: get the uploaded file from multi-form value and replace it with "lib/config-template.conf"
		respondWithError(w, 500, "Not implemented yet")
		Info.Println("CONFIG Command => Not implemented yet")
	})

	r.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		//a reset factory function which reset the configuration template file
		//TODO: mgh, 2017: read the original template file from "lib/config-template-backup.conf" and write to "lib/config-template"
		respondWithError(w, 500, "Not implemented yet")
		Info.Println("reset Command => Not implemented yet")
	})

	//start emmbed http server on port 6663
	//TODO: mgh, 2017: the port number must read from configuarion file
	Info.Println("Telegraf Http Harness Version 0.3 Started")
	Info.Println("Listen on port 6663 for incoming http commands")
	err := http.ListenAndServe(":6663", r)
	log.Fatal(err)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func respondWithJSONStr(w http.ResponseWriter, code int, payload string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(payload))
}

func httpFormToKeyMap(form map[string][]string) map[string]string {
	var result map[string]string
	result = make(map[string]string)
	for key, value := range form {
		if len(value) > 0 {
			result[key] = value[0]
		}
	}
	return result
}

/*	r.HandleFunc("/test-form", func(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	var jsonMap map[string]interface{}
	jsonMap = make(map[string]interface{})
	for key, value := range r.Form {
		jsonMap[key] = value
	}
	_, header, _ := r.FormFile("file")
	if header != nil {
		jsonMap["filename"] = header.Filename
	}
	respondWithJSON(w, 200, jsonMap)
})*/

/*	r.HandleFunc("/ddd", func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte ( r.FormValue("isna")))
		r.ParseMultipartForm(32 << 20)
		//r.ParseMultipartForm(32 << 20)
		var jsonMap map[string]interface{}
		jsonMap = make(map[string]interface{})
		for key, value := range r.Form {
			if len(value) > 0 {
				//decode as int
				intVal, err := strconv.Atoi(value[0])
				if err == nil {
					jsonMap[key] = intVal
					continue
				}

				//decode as float
				floatVal, err := strconv.ParseFloat(value[0], 64)
				if err == nil {
					jsonMap[key] = floatVal
					continue
				}

				//decode as array
				arrayVal, err := decodeArray(value[0])
				if err == nil {
					jsonMap[key] = arrayVal
					continue
				}
				//decode as string
				jsonMap[key] = value[0]
			}
		}
		str := MapToTomlString(jsonMap)
		w.Write([]byte(str))
})*/

/*r.HandleFunc("/ccc", func(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 100)
	var jsonMap map[string]interface{}
	jsonMap = make(map[string]interface{})
	jsonMap["agent.interval"] = r.FormValue("agent.interval")
	jsonMap["agent.hostname"] = r.FormValue("agent.hostname")
	jsonMap["outputs.influxdb.urls"] = r.FormValue("outputs.influxdb.urls")
	jsonMap["outputs.influxdb.database"] = r.FormValue("outputs.influxdb.database")
	jsonMap["interval"] = r.FormValue("outputs.influxdb.precision") //default=s
	_, header, _ := r.FormFile("file")
	jsonMap["filename"] = header.Filename
	respondWithJSON(w, 200, jsonMap)
})*/

/*r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("/config called")
	vars := mux.Vars(r)
	println(vars["interval"])
	var Buf bytes.Buffer
	// in your case file would be fileupload
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	ioutil.WriteFile("lib/config-template.conf",file,)
	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile("/tmp/dat1", d1, 0644)
	check(err)
	contents := Buf.String()
	writeInfluxTemplateConfigFile(contents)
	Buf.Reset()
	return

})
*/
