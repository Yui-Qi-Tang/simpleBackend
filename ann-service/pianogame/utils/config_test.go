package utils_test

import (
	"simpleBackend/ann-service/pianogame/utils"
	"testing"
)

type configData struct {
	filePath   string
	successMsg string
	failMsg    string
	data       interface{}
}

func configList() []configData {
	return []configData{
		configData{
			filePath:   "../../config/website/config.yaml",
			successMsg: "ok",
			failMsg:    "fail",
		},
		configData{
			filePath:   "../../config/api/config.yaml",
			successMsg: "ok",
			failMsg:    "fail",
		},
		configData{
			filePath:   "../../config/ssl/config.yaml",
			successMsg: "ok",
			failMsg:    "fail",
		},
		configData{
			filePath:   "../../config/auth/config.yaml",
			successMsg: "ok",
			failMsg:    "fail",
		},
		configData{
			filePath:   "../../config/database/mongo/config.yaml",
			successMsg: "ok",
			failMsg:    "fail",
		},
		configData{
			filePath:   "../../config/grpc/config.yaml",
			successMsg: "ok",
			failMsg:    "fail",
		},
	}
}

func TestLoadConfig(t *testing.T) {
	// Just test read file function and yaml Unmarshal function are correctly
	configList := configList()
	structNothing := &struct{}{}
	for _, v := range configList {
		utils.LoadYAMLConfig(v.filePath, v.failMsg, v.successMsg, structNothing)
	}
	t.Log("Load config files pass")

}
