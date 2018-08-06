package main

import (
	"fmt"
	"net/http"

	"encoding/json"
	"github.com/gorilla/mux"
	"influxdb-line-protocol-to-json"
	"log"
	"os/exec"
	"strings"
)

var metrics = [1]string{"win_cpu"}

type Agent struct {
	db string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
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
		jsonStr := influxdb_line_protocol_to_json.LinesToJson(outstr)
		//w.Write([]byte(outstr + "\n\n" + jsonStr))
		respondWithJSONStr(w, 200, jsonStr)
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
