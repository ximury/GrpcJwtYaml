package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type configuration struct {
	Enabled bool
	Path    string
}

func main() {
	//file, er := os.Open("./json2/conf.json")
	file, er := os.Open("./json/conf.json")
	fmt.Println(er)

	defer file.Close()

	decoder := json.NewDecoder(file)
	var conf = configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(conf.Enabled, conf.Path)
}