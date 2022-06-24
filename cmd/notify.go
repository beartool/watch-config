package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"watch-config/configs"
)

type NotifyToSync struct {
	SshIp       string
	SshPort     string
	SshUser     string
	SshIdentify string
	Configs     *configs.TomlConfig
}

func NewNotifyToSync() *NotifyToSync {
	destination := configs.Config.Destination
	m := &NotifyToSync{
		SshIp:       destination.SshIp,
		SshPort:     destination.SshPort,
		SshUser:     destination.SshUser,
		SshIdentify: destination.SshIdentity,
		Configs:     configs.Config,
	}
	return m
}

// CreateNotify
// @Description: 同步远端
// @receiver n
// @param name
func (n *NotifyToSync) CreateNotify(name string) {
	err := n.FileSync(name)
	if err != nil {
		return
	}
	n.AfterOperation()
}

// RemoveNotify
// @Description: 文件删除，远程文件移动到临时目录 rm存在风险
// @receiver n
// @param name
// @return err
func (n *NotifyToSync) RemoveNotify(name string) {
	remoteDir := name
	sourceDir := n.Configs.Source.SourceDir
	relativePath, err := filepath.Rel(remoteDir, sourceDir)
	if targetDir := n.Configs.Source.TargetDir; targetDir != "" {
		if err != nil {
			return
		}
		remoteDir = filepath.Join(targetDir, relativePath)
	}

	mvPath := filepath.Join("/tmp/sync/", relativePath)

	c := fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \" mv -f %s %s\"", n.SshPort, n.SshIdentify, n.SshUser, n.SshIp, remoteDir, mvPath)
	log.Printf("删除远程文件：%s\n", remoteDir)
	cmd := exec.Command("bash", "-c", c)
	err = cmd.Run()
	if err != nil {
		return
	}
}

// FileSync
// @Description: 同步文件到远端
// @receiver this
// @param name
func (n *NotifyToSync) FileSync(name string) (err error) {
	remoteDir := filepath.Dir(name)
	if targetDir := n.Configs.Source.TargetDir; targetDir != "" {
		sourceDir := n.Configs.Source.SourceDir
		relativePath, err := filepath.Rel(remoteDir, sourceDir)
		if err == nil {
			remoteDir = filepath.Join(targetDir, relativePath)
		}
	}

	c := fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \" mkdir -p %s \"", n.SshPort, n.SshIdentify, n.SshUser, n.SshIp, remoteDir)
	log.Printf("创建远程目录：%s\n", remoteDir)
	cmd := exec.Command("bash", "-c", c)
	err = cmd.Run()
	if err != nil {
		return
	}
	// 目录则不执行同步
	fileInfo, err := os.Stat(name)
	if err != nil || fileInfo.IsDir() {
		return
	}
	filename := filepath.Base(name)
	extname := path.Ext(filename)
	if extname == ".swp" {
		log.Printf("临时文件不同步 %s", name)
		return
	}
	remotePath := filepath.Join(remoteDir, filename)

	c = fmt.Sprintf("scp -P %s  -i %s  %s  %s@%s:%s", n.SshPort, n.SshIdentify, name, n.SshUser, n.SshIp, remotePath)
	log.Printf("执行复制：%s\n", c)
	cmd = exec.Command("bash", "-c", c)
	err = cmd.Run()
	if err != nil {
		return
	}
	n.ContentReplace(remotePath)
	return nil
}

// ContentReplace
// @Description: 批量替换文件中的文本内容
// @receiver this
// @param name
func (n *NotifyToSync) ContentReplace(name string) {
	sedRule := n.Configs.ReplaceRule.SedRule
	if sedRule == "" {
		return
	}

	c := fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \"sed -i %s %s \"", n.SshPort, n.SshIdentify, n.SshUser, n.SshIp, sedRule, name)
	log.Printf("执行替换：%s\n", c)
	cmd := exec.Command("bash", "-c", c)
	err := cmd.Run()
	if err != nil {
		log.Println("执行文件内容替换失败 文件路径" + name)
		return
	}
}

// AfterOperation
// @Description: 文件同步后置远程操作
// @receiver this
// @param name//
func (n *NotifyToSync) AfterOperation() {
	command := n.Configs.Command.CompletedCmd

	if command == "" {
		return
	}
	log.Println(n.SshIdentify)
	c := fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \" %s \"", n.SshPort, n.SshIdentify, n.SshUser, n.SshIp, command)
	log.Printf("启动服务：%s\n", c)
	cmd := exec.Command("bash", "-c", c)
	err := cmd.Run()
	if err != nil {
		return
	}
}
