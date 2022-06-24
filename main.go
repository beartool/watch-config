//+build linux darwin windows

package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"watch-config/cmd"
	"watch-config/configs"
)

func main() {
	app := &cli.App{
		Name:  "watcher-file",
		Usage: "sync config to other environment",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "sync.toml",
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: func(c *cli.Context) error {
			tomlPath := c.String("file")
			tomlFile, err := os.Stat(tomlPath)
			if err != nil || tomlFile.IsDir() {
				log.Println("toml file not found !")
				return nil
			}
			exec(tomlPath)
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// exec
// @Description: 执行监控程序 监控kill 和 ctrl+c 退出
// @param config
func exec(config string) {
	syncConfig, _ := configs.ReadConf(config)

	sourceDir := syncConfig.Source.SourceDir
	if sourceDir == "" {
		log.Fatalf("未设置监控目录，请设置source_dir")
	}
	fileInfo, err := os.Stat(sourceDir)
	if err != nil || !fileInfo.IsDir() {
		log.Fatalf("监控路径不存在或不是目录，请检查source_dir")
	}

	watch := cmd.NewWatcherFile()
	watch.WatchDir(sourceDir)

	osc := make(chan os.Signal, 1)
	signal.Notify(osc, syscall.SIGTERM, syscall.SIGINT)
	<-osc

	fmt.Println("程序退出")
}
