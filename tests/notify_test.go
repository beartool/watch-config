package tests

import (
	"log"
	"testing"
	"watch-config/cmd"
	"watch-config/configs"
)

func init() {
	config := "../configs/file.toml"
	_, err := configs.ReadConf(config)
	if err != nil {
		log.Fatal("读取配置文件失败")
		return
	}
}

func TestAfterOperation(t *testing.T) {
	notify := cmd.NewNotifyToSync()

	notify.AfterOperation()
}

func TestFileSync(t *testing.T) {
	notify := cmd.NewNotifyToSync()

	err := notify.FileSync("/Users/hongxue.cao/work/www/go/fsevents/example/a.text")
	if err != nil {
		return
	}
}
