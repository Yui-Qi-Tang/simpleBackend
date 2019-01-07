package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type socketClient struct {
	Server string   `yaml:"server"`
	Port   string   `yaml:"port"`
	Cmds   []string `yaml:"cmds"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// Read config file
	yamlFile, ioErr := ioutil.ReadFile("config/socket_client.yaml") // open file and read
	var scData socketClient
	if ioErr != nil {
		errorStr := fmt.Sprintf("Read config file error! %v", ioErr)
		panic(errorStr)
	} // fi

	// yaml Unmarshal
	yamlError := yaml.Unmarshal(yamlFile, &scData)
	if yamlError != nil {
		errorStr := fmt.Sprintf("error while unmarshal from db config: %v", yamlError)
		panic(errorStr)
	}
	// combine IP:PORT
	serverAddr := fmt.Sprintf("%s:%s", scData.Server, scData.Port)
	// string for query
	var queryStr string
	// log file stream
	var fileStreamErr error
	var f *os.File
	f, fileStreamErr = os.Create("./result.log")
	check(fileStreamErr)
	defer f.Close() // close file stream

	// send cmd to server via socket
	for _, v := range scData.Cmds {
		conn, err := net.Dial("tcp", serverAddr) // TODO put connection type to config
		if err != nil {
			fmt.Println(err)
		} else {
			start := time.Now()

			// logging query string
			queryFile := []byte(v)
			queryFile = append(queryFile, '\t')
			_, fileStreamErr = f.Write(queryFile)
			check(err)
			f.Sync()

			// send query for socket server
			readBuf := make([]byte, 1024)
			queryStr = fmt.Sprintf("cmd=query&&audioPath=%s", v)
			fmt.Println(queryStr)
			conn.Write([]byte(queryStr))

			// receive result from server
			conn.Read(readBuf)
			readBuf = append(readBuf, byte('\n'))
			_, fileStreamErr = f.Write(readBuf)
			check(err)
			f.Sync()

			// show result
			t := time.Now()
			elapsed := t.Sub(start)
			fmt.Println(string(readBuf[:]), elapsed)

		} // fi
		conn.Close()
	} // for
}
