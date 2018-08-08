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
	"strings"
)

var influxCmd *exec.Cmd
var influxCmdIsRunning bool

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("/config called")

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
		writeInfluxConfigFile("", "lib\\config.conf")
		var telegrafArg = "-test -config lib\\config.conf"
		out, err := exec.Command("lib\\t.exe", strings.Split(telegrafArg, " ")...).Output()
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("The date is %s\n", out)
		outstr := fmt.Sprintf("%s", out)
		outstr = strings.Replace(outstr, "> ", "", -1)
		jsonStr := influxjson.LinesToJson(outstr)
		//w.Write([]byte(outstr + "\n\n" + jsonStr))
		respondWithJSONStr(w, 200, jsonStr)
		//to json
	})

	r.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		if influxCmdIsRunning == true {
			respondWithJSON(w, 200, map[string]string{"status": "already started"})
			return
		}
		writeInfluxConfigFile("", "lib\\config.conf")
		var telegrafArg = "-config lib\\config.conf"
		influxCmd = exec.Command("lib\\t.exe", strings.Split(telegrafArg, " ")...)
		err := influxCmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		influxCmdIsRunning = true
		respondWithJSON(w, 200, map[string]string{"status": "started"})
		//to json
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

	http.ListenAndServe(":8080", r)
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
