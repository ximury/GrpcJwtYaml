package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//type GoWeb struct {
//	Port string `yaml:"port"`
//}
type Test struct {
	Port string `yaml:"port"`
}
type GoWeb struct {
	Test Test `yaml:"Test"`
}
type Config struct {
	//成员名称首字母大写
	GrpcSerType string `yaml:"grpcSerType"`
	GrpcSerURL string `yaml:"grpcSerURL"`
	GrpcCliURL string `yaml:"grpcCliURL"`
	GoWeb GoWeb `yaml:"go_web"`
}
func main() {
	var setting Config
	config, err := ioutil.ReadFile("./yaml2/config.yaml")
	if err != nil {
		fmt.Print(err)
	}
	err = yaml.Unmarshal(config, &setting)
	if err != nil {


		fmt.Print(err)
	}
	fmt.Println(setting.GrpcSerType)
	fmt.Println(setting.GrpcSerURL)
	fmt.Println(setting.GrpcCliURL)
	//fmt.Println(setting.GoWeb.Port)
	fmt.Println(setting.GoWeb.Test.Port)
}
