package main

import (
	"encoding/json"
	"io/ioutil"
)

var CONFIG *Config

type Config struct {
	Bin 	 string
	ChainDir string
	ChainCount int
	NodesPerChain int
}

func NewConfig(path string) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 
	}
	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil { return }
	CONFIG = config
	return
}

func Setup() (err error) {
	return
}