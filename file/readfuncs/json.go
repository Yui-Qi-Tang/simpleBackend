package readfilefunc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func getFilePath(name string) string {
	return fmt.Sprintf("data/%s", name)
}

func readDataBytes(fileName string) []byte {
	dataJSON, err := os.Open(getFilePath(fileName))
	ErrorCheck(err)
	defer dataJSON.Close()

	bytesData, err := ioutil.ReadAll(dataJSON)
	ErrorCheck(err)
	return bytesData
}

func getMainTupleData() []MainTuple {
	dataBytes := readDataBytes("data1.json")
	var mainTuples []MainTuple
	json.Unmarshal(dataBytes, &mainTuples) // get data1.json
	return mainTuples
}

func getKeyCell8() []KeyCell8 {
	dataBytes := readDataBytes("data2.json")
	var result []KeyCell8
	json.Unmarshal(dataBytes, &result) // get data1.json
	return result
}

func getCell49() []Cell49 {
	dataBytes := readDataBytes("data3.json")
	var f interface{}
	json.Unmarshal(dataBytes, &f)
	itemsMap := f.(map[string]interface{})
	var result []Cell49

	for _, v := range itemsMap {
		t := v.(map[string]interface{}) // decode
		result = append(result, Cell49{
			Cell4: t["cell4"].(string),
			Cell9: t["cell9"].(string),
		})
	}
	return result
}
