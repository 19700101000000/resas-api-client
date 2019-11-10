package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"resas-api/env"
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
		err    error
	)
	flag.Parse()

	fmt.Println("mode:", *mode)
	switch *mode {
	case "get":
		err = get(*apiKey, *path, *output)
	case "sql":
		err = sql(*table, *cols, *input, *output)
	}

	if err != nil {
		fmt.Printf("%v", err)
	}
}

func sql(table, cols, input, output string) error {
	if table == "" {
		fmt.Println("You must set table.", messageHelp)
		return nil
	}
	if input == "" {
		fmt.Println("You must set input.", messageHelp)
		return nil
	}
	if output == "" {
		fmt.Println("You must set output.", messageHelp)
		return nil
	}

	_, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Println("Error input:", input)
		return err
	}

	
	return nil
}

func get(apiKey, path, output string) error {
	if apiKey == "" {
		fmt.Println("You must set key.", messageHelp)
		return nil
	}
	if path == "" {
		fmt.Println("You must set path.", messageHelp)
		return nil
	}

	fullPath := env.Endpoint + path

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
