package main

import (
	"flag"
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

func check(e error, errMsg string) {
	if e != nil {
		errorStr := fmt.Sprintf("%s: %v", errMsg, e)
		panic(errorStr)
	}
}

func main() {

	// Usage
	logPath := flag.String("log", "./result.log", "log file path")
	configPath := flag.String("config", "config/config.yaml", "config file path")
	flag.Parse() // pare variables from commnad line

	// Read config yaml file
	yamlFile, ioErr := ioutil.ReadFile(*configPath) // open file and read
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
	f, fileStreamErr = os.Create(*logPath)
	check(fileStreamErr, "File stream error!")
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
			check(err, "Write file error")
			f.Sync()

			// send query for socket server
			readBuf := make([]byte, 1024)
			queryStr = fmt.Sprintf("cmd=query&&audioPath=%s", v)
			fmt.Println(queryStr)
			conn.Write([]byte(queryStr))

			// receive result from server
			conn.Read(readBuf)
			readBuf = append(readBuf, byte('\n'))
			_, fileStreamErr = f.Write(readBuf) // TODO: write json string
			check(err, "write file error")
			f.Sync()

			// show result
			t := time.Now()
			elapsed := t.Sub(start)
			fmt.Println(string(readBuf[:]), elapsed)

		} // fi
		conn.Close()
	} // for
}
