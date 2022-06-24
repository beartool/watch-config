package tests

import (
	"fmt"
	"log"
	"regexp"
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

	err := notify.FileSync("/home/www/watch-config/fsevents/example/text")
	if err != nil {
		return
	}
}

func TestRemove(t *testing.T) {
	notify := cmd.NewNotifyToSync()

	notify.RemoveNotify("/home/www/watch-config/fsevents/example/text")
}

func TestMatchString(t *testing.T) {
	matched, err := regexp.MatchString(`\w+\.(swp|swx|\w+~)$`, "a.test~a")
	fmt.Println(matched, err) //true <nil>
}
