package main

import (
	encodeJson "encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"resas-api/env"
	"resas-api/structs"
	"strconv"
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

	json, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Println("Error input:", input)
		return err
	}

	var sql string

	switch table {
	case "prefectures":
		sql, err = getSqlPrefectures(&json, cols)
	}
	if err != nil {
		fmt.Println("Error create sql:", input)
		return err
	}

	if output == "" {
		fmt.Println(sql)
		return nil
	}

	err = ioutil.WriteFile(output, []byte(sql), 0644)
	if err != nil {
		fmt.Println("Error output:", output)
		return err
	}
	fmt.Println("Success output:", output)

	return nil
}
func getSqlPrefectures(json *[]byte, cols string) (string, error) {
	var prefectures structs.Prefectures
	err := encodeJson.Unmarshal(*json, &prefectures)
	if err != nil {
		fmt.Println("Error read json:")
		return "", err
	}

	// fmt.Println(getSqlColumns(cols, structs.Prefecture{}))
	tags := getSqlColumns(cols, structs.Prefecture{})
	var names, values []string
	// struct
	for index, prefecture := range prefectures.Result {
		t := reflect.TypeOf(prefecture)
		// field
		colVals := make([]string, 0)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("json")
			if tags[tag] == "" {
				continue
			}

			if index == 0 {
				names = append(names, tags[tag])
			}

			switch field.Name {
			case "PrefName":
				colVals = append(colVals, "'"+prefecture.PrefName+"'")
			case "PrefCode":
				colVals = append(colVals, strconv.Itoa(prefecture.PrefCode))
			}
		}
		value := "(" + strings.Join(colVals, ",") + ")"
		values = append(values, value)
	}
	sqlf := "INSERT INTO prefectures(%s) VALUES %s"
	sql := fmt.Sprintf(sqlf, strings.Join(names, ","), strings.Join(values, ","))

	return sql, nil
}
func getSqlColumns(cols string, jsonSchema interface{}) map[string]string {
	result := make(map[string]string)

	if cols == "" {
		t := reflect.TypeOf(jsonSchema)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("json")
			result[tag] = tag
		}
		return result
	}

	cols = strings.ReplaceAll(cols, " ", "")
	for _, v := range strings.Split(cols, ",") {
		if strings.Contains(v, ">") {
			s := strings.Split(v, ">")
			result[s[0]] = s[1]
			continue
		}
		result[v] = v
	}
	return result
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
