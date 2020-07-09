package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//Nginx nginx  配置
type Nginx struct {
	Port int `yaml:"Port"`
	LogPath string `yaml:"LogPath"`
	Path string `yaml:"Path"`
}
//Config   系统配置配置
type Config struct{
	Name string `yaml:"SiteName"`
	Addr string `yaml:"SiteAddr"`
	HTTPS bool `yaml:"Https"`
	SiteNginx  Nginx `yaml:"Nginx"`
}

func main() {

	var setting Config
	config, err := ioutil.ReadFile("./yaml/test.yaml")
	if err != nil {
		fmt.Print(err)
	}
	_ = yaml.Unmarshal(config, &setting)

	fmt.Println(setting.Name)
	fmt.Println(setting.Addr)
	fmt.Println(setting.HTTPS)
	fmt.Println(setting.SiteNginx.Port)
	fmt.Println(setting.SiteNginx.LogPath)
	fmt.Println(setting.SiteNginx.Path)

}