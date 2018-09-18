package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"text/template"
)

/*func writeInfluxConfigFile(url string,database string,interval string, path string) {*/
func writeInfluxConfigFile(params map[string]string, pathToTemplateFile string, pathToConfigFile string) {
	configStr := fillConfigTemplate(params, pathToTemplateFile)
	bytes := []byte(configStr)
	// write the whole body at once
	err := ioutil.WriteFile(pathToConfigFile, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func fillConfigTemplate(data map[string]string, path string) string {
	templStr, err := readFile(path)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("test").Parse(templStr)
	if err != nil {
		panic(err)
	}
	builder := &strings.Builder{}
	if err := tmpl.Execute(builder, data); err != nil {
		panic(err)
	}
	return builder.String()
}
func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

//map to string using TOML format
func MapToTomlString(m map[string]interface{}) string {
	var result strings.Builder
	for key, val := range m {
		v := reflect.ValueOf(val)
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			var arrayStr strings.Builder
			n := v.Len()
			var strs []string
			for i := 0; i < n; i++ {
				strs = append(strs, reflectToString(v.Index(i)))
			}
			arrayStr.WriteString("[")
			arrayStr.WriteString(strings.Join(strs, ","))
			arrayStr.WriteString("]")
			fmt.Fprintf(&result, "%s=%s\n", key, arrayStr.String())
		default:
			fmt.Fprintf(&result, "%s=%s\n", key, reflectToString(v))
		}
	}
	return result.String()
}

func reflectToString(obj reflect.Value) string {
	var result strings.Builder
	switch obj.Kind() {
	case reflect.Int:
		fmt.Fprintf(&result, "%d", obj.Int())
	case reflect.Float64:
		fmt.Fprintf(&result, "%f", obj.Float())
	case reflect.String:
		fmt.Fprintf(&result, "\"%s\"", obj.String())
	}
	return result.String()
}

func decodeArray(str string) ([]interface{}, error) {
	var restult []interface{}
	dec := json.NewDecoder(strings.NewReader(str))
	err := dec.Decode(&restult)
	if err != nil {
		return nil, err
	}

	return restult, nil
}
