package service

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/guoruibiao/rediscustomsync/models/dao"
	"github.com/guoruibiao/rediscustomsync/models/dao/redisdao"
	"os"
	"testing"
)

var appConfig dao.AppConfig

func init() {

	metaData, err := toml.DecodeFile("../../conf/app.toml", &appConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}else{
		fmt.Printf("%#v\n", metaData)
		fmt.Printf("%#v\n", appConfig)
	}
	// 初始化读写的 redisdao 配置
	redisdao.Init(appConfig)
}

func TestTransfer(t *testing.T) {
	success, failed, err := TransferKeysfile(appConfig)
	if err != nil {
		t.Error(err)
	}else{
		t.Log("Success=", success)
		t.Log("Failed=", failed)
	}
}
