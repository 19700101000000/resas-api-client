package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	// "reflect"
	encodeJson "encoding/json"
	"resas-api/env"
	"resas-api/structs"
	"strings"
)

const (
	messageHelp = "If you need help, enter `resas-api -h`."
)

func main() {
	var (
		mode   = flag.String("mode", "get", "this exec modes:")
		apiKey = flag.String("key", "", "this is RESAS-API's API-KEY.")
		path   = flag.String("path", "", "this is RESAS-API's GET-PATH.")
		output = flag.String("out", "", "this is output file name.")
		input  = flag.String("in", "", "this is input file name.")
		table  = flag.String("table", "", "this is SQL's parse type.")
		cols   = flag.String("cols", "", "this is SQL's columns.")
		params = flag.String("params", "", "this is GET parameters.")
		err    error
	)
	flag.Parse()

	fmt.Println("mode:", *mode)
	switch *mode {
	case "get":
		err = get(*apiKey, *path, *params, *output)
	case "get_cities":
		err = getCities(*apiKey, *output)
	case "sql":
		err = sql(*table, *cols, *input, *output)
	}

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func sql(table, cols, input, output string) error {
	if table == "" {
		fmt.Println("You must set table.", messageHelp)
		return nil
	}
	if input == "" {
		fmt.Println("You must set in.", messageHelp)
		return nil
	}
	if output == "" {
		fmt.Println("You must set out.", messageHelp)
		return nil
	}

	json, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Println("Error input:", input)
		return err
	}

	switch table {
	case "prefectures":
		return sqlPrefectures(&json, cols)
	}

	return nil
}
func sqlPrefectures(json *[]byte, cols string) error {
	var prefectures structs.Prefectures
	err := encodeJson.Unmarshal(*json, &prefectures)
	if err != nil {
		fmt.Println("Error read json:")
		return err
	}

	getSqlColumns(cols, structs.Prefecture{})

	return nil
}
func getSqlColumns(cols string, jsonSchema interface{}) {
	fmt.Println(jsonSchema)
}

func get(apiKey, path, params, output string) error {
	if apiKey == "" {
		fmt.Println("You must set key.", messageHelp)
		return nil
	}
	if path == "" {
		fmt.Println("You must set path.", messageHelp)
		return nil
	}

	fullPath := env.Endpoint + path
	if params != "" {
		params = strings.ReplaceAll(params, " ", "")
		params = strings.ReplaceAll(params, ",", "&")
		fullPath += "?" + params
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", fullPath, nil)
	if err != nil {
		fmt.Println("You Don't create Request:")
		return err
	}

	req.Header.Add("X-API-KEY", apiKey)

	fmt.Println("start connection...", fullPath)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error connect: %s", fullPath)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error io:")
		return err
	}

	if output == "" {
		fmt.Printf("%s", body)
		return nil
	}

	err = ioutil.WriteFile(output, body, 0644)
	if err != nil {
		fmt.Println("Error output:", output)
		return err
	}
	fmt.Println("Success output:", output)

	return nil
}

func getCities(apiKey, out string) error {
	var (
		start  = 1
		end    = 47
		finish = make(chan bool, end)
		count  int8
	)
	for i := start; i <= end; i++ {

		path := "api/v1/cities"
		path += fmt.Sprintf("?prefCode=%d", i)
		output := fmt.Sprintf("%scities_%d.json", out, i)
		go getCity(apiKey, path, output, finish)
	}

	for i := start; i <= end; i++ {
		if <-finish {
			count++
		}
	}
	fmt.Println("Result get cities:", count)
	return nil
}
func getCity(apiKey, path, out string, finish chan<- bool) {
	fmt.Println(out)
	err := get(apiKey, path, "", out)
	if err != nil {
		fmt.Println("Error:", out)
		finish <- false
	}
	finish <- true
}
