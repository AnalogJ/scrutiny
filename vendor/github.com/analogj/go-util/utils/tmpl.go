package utils

import (
	"bytes"
	"encoding/json"
	"hash/fnv"
	"text/template"
)


func PopulatePathTemplate(pathTmplContent string, data interface{}) (string, error){
	tmplFilepath, err := PopulateTemplate(pathTmplContent, data)
	if err != nil {
		return "", nil
	}
	tmplFilepath, err = ExpandPath(tmplFilepath)
	if err != nil {
		return "", nil
	}
	return tmplFilepath, nil
}


func PopulateTemplate(tmplContent string, data interface{}) (string, error) {
	//set functions
	fns := template.FuncMap{
		"uniquePort": UniquePort,
		"expandPath": ExpandPath,
	}

	// prep the template, set the option
	tmpl, err := template.New("populate").Option("missingkey=error").Funcs(fns).Parse(tmplContent)
	if err != nil {
		return "", err
	}

	//specify that any missing keys in the template will throw an error
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	//convert buffered content to string
	return buf.String(), nil
}

// https://play.golang.org/p/k8bws03uid
func UniquePort(data interface{}) (int, error) {

	var contentData []byte
	switch in := data.(type) {
	case string:
		contentData = []byte(in)
	default:
		jsonData, err := json.Marshal(StringifyYAMLMapKeys(in))
		if err != nil {
			return 0, err
		}
		contentData = jsonData
	}

	hash := fnv.New32a()
	hash.Write(contentData)

	//last port - last privileged port.
	portRange := 65535 - 1023

	uniquePort := (hash.Sum32() % uint32(portRange)) + 1023
	return int(uniquePort), nil
}
