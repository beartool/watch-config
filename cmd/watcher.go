package cmd

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
)

type WatcherFile struct {
	watch  *fsnotify.Watcher
	notify *NotifyToSync
}

func NewWatcherFile() *WatcherFile {
	w := new(WatcherFile)
	w.watch, _ = fsnotify.NewWatcher()
	w.notify = NewNotifyToSync()
	return w
}

// WatchDir
// @Description: 添加监控目录
// @receiver this
// @param dir
func (this *WatcherFile) WatchDir(dir string) {
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//判断是否为目录，监控目录,目录下文件也在监控范围内，不需要加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = this.watch.Add(path)
			if err != nil {
				return err
			}
			log.Println("监控 : ", path)
		}
		return nil
	})

	go this.WatchEvent() //协程
}

// WatchEvent
// @Description: 处理监控事件
// @receiver this
func (this *WatcherFile) WatchEvent() {
	for {
		select {
		case ev := <-this.watch.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					//获取新创建文件的信息，如果是目录，则加入监控中
					file, err := os.Stat(ev.Name)
					if err == nil && file.IsDir() {
						this.watch.Add(ev.Name)
						log.Println("添加监控 : ", ev.Name)
					}
					if !file.IsDir() {
						log.Println("创建文件 : ", ev.Name)
						this.notify.CreateNotify(ev.Name)
					}
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					log.Println("写入文件 : ", ev.Name)
					this.notify.CreateNotify(ev.Name)
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					//如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						this.watch.Remove(ev.Name)
						log.Println("删除监控 : ", ev.Name)
					}
					this.notify.RemoveNotify(ev.Name)
				}

				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					//如果重命名文件是目录，则移除监控 ,注意这里无法使用os.Stat来判断是否是目录了
					//因为重命名后，go已经无法找到原文件来获取信息了,所以简单粗爆直接remove
					//同时还会有一个创建的通知，所以移除原来的文件就行
					log.Println("重命名文件 : ", ev.Name)
					this.watch.Remove(ev.Name)
					this.notify.RemoveNotify(ev.Name)
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					// 权限先忽略
					log.Println("修改权限 : ", ev.Name)
				}
			}
		case err := <-this.watch.Errors:
			{
				log.Println("error : ", err)
				return
			}
		}
	}
}
