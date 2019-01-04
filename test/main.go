package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"gopkg.in/yaml.v2"
)

type socketClient struct {
	Server string   `yaml:"server"`
	Port   string   `yaml:"port"`
	Cmds   []string `yaml:"cmds"`
}

func main() {
	yamlFile, ioErr := ioutil.ReadFile("config/socket_client.yaml") // open file and read
	var scData socketClient
	if ioErr != nil {
		errorStr := fmt.Sprintf("Read config file error! %v", ioErr)
		panic(errorStr)
	} // fi

	yamlError := yaml.Unmarshal(yamlFile, &scData)
	if yamlError != nil {
		errorStr := fmt.Sprintf("error while unmarshal from db config: %v", yamlError)
		panic(errorStr)
	}
	serverAddr := fmt.Sprintf("%s:%s", scData.Server, scData.Port)
	var queryStr string
	for _, v := range scData.Cmds {
		conn, err := net.Dial("tcp", serverAddr) // TODO put connection type to config
		if err != nil {
			fmt.Println(err)
		} else {
			start := time.Now()
			readBuf := make([]byte, 1024)
			queryStr = fmt.Sprintf("cmd=query&&audioPath=%s", v)
			fmt.Println(queryStr)
			conn.Write([]byte(queryStr))
			conn.Read(readBuf)
			t := time.Now()
			elapsed := t.Sub(start)
			fmt.Println(string(readBuf[:]), elapsed)

		} // fi
		conn.Close()
	} // for
}
