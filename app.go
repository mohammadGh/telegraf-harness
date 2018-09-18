package main

import (
	"fmt"
	"net/http"

	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	influxjson "github.com/mohammadGh/influxdb-line-protocol-to-json"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var influxCmd *exec.Cmd
var influxCmdIsRunning bool

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/test-form", func(w http.ResponseWriter, r *http.Request) {
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
	})
	r.HandleFunc("/ddd", func(w http.ResponseWriter, r *http.Request) {
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
		/*
			_, header, _ := r.FormFile("file")
			jsonMap["filename"] = header.Filename
			respondWithJSON(w, 200, jsonMap)*/
	})
	r.HandleFunc("/ccc", func(w http.ResponseWriter, r *http.Request) {
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
	})

	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
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
		// do something with the contents...
		// I normally have a struct defined and unmarshal into a struct, but this will
		// work as an example
		contents := Buf.String()
		fmt.Println(contents)
		// I reset the buffer in case I want to use it again
		// reduces memory allocations in more intense projects
		Buf.Reset()
		// do something else
		// etc write header
		return

	})

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

		respondWithJSONStr(w, 200, jsonStr)
	})

	r.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {

		if influxCmdIsRunning == true {
			respondWithJSON(w, 200, map[string]string{"status": "already started"})
			return
		}
		params := r.URL.Query()
		if len(params) == 0 {
			params = r.Form
		}
		formAsKeyMaps := httpFormToKeyMap(params)
		writeInfluxConfigFile(formAsKeyMaps, "lib/config-template.conf", "lib/config.conf")
		var telegrafArg = "-config lib/config.conf"
		influxCmd = exec.Command("lib/t.exe", strings.Split(telegrafArg, " ")...)
		err := influxCmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		influxCmdIsRunning = true
		respondWithJSON(w, 200, map[string]string{"status": "started"})
	})

	r.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if influxCmd == nil {
			respondWithJSON(w, 200, map[string]string{"status": "already stopped"})
			return
		}
		err := influxCmd.Process.Kill()
		if err != nil {
			log.Fatal(err)
		}
		influxCmd = nil
		influxCmdIsRunning = false
		respondWithJSON(w, 200, map[string]string{"status": "stopped"})
		//to json
	})

	r.HandleFunc("/Config", func(w http.ResponseWriter, r *http.Request) {

		if influxCmdIsRunning == true {
			respondWithJSON(w, 200, map[string]string{"status": "already started"})
			return
		}
		formAsKeyMaps := httpFormToKeyMap(r.Form)
		writeInfluxConfigFile(formAsKeyMaps, "lib/config-template.conf", "lib/config.conf")
		var telegrafArg = "-config lib/config.conf"
		influxCmd = exec.Command("lib/t.exe", strings.Split(telegrafArg, " ")...)
		err := influxCmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		influxCmdIsRunning = true
		respondWithJSON(w, 200, map[string]string{"status": "started"})
		//to json
	})

	http.ListenAndServe(":6663", r)
}

/*func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}*/

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
