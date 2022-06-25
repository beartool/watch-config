package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
	res, err := n.FileSync(name)
	if err != nil || res == false {
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
	isTmp := n.CheckIsTmpFile(name)
	if isTmp {
		return
	}
	remoteDir := name
	sourceDir := n.Configs.Source.SourceDir
	relativePath, err := filepath.Rel(sourceDir, remoteDir)
	if targetDir := n.Configs.Source.TargetDir; targetDir != "" {
		if err != nil {
			return
		}
		remoteDir = filepath.Join(targetDir, relativePath)
	}

	// 获取删除时的文件目录
	deleteDir := n.Configs.Source.DeleteDir
	if deleteDir == "" {
		deleteDir = "/tmp/watch-file"
	}

	mvPath := filepath.Join(deleteDir, relativePath)
	mvDir, _ := filepath.Split(mvPath)
	c := fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \" mkdir -p %s \"", n.SshPort, n.SshIdentify, n.SshUser, n.SshIp, mvDir)
	cmd := exec.Command("bash", "-c", c)
	err = cmd.Run()
	if err != nil {
		return
	}

	c = fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \" mv -f %s %s\"", n.SshPort, n.SshIdentify, n.SshUser, n.SshIp, remoteDir, mvPath)
	log.Printf("删除远程文件：%s, 移动至 %s\n", remoteDir, mvPath)
	cmd = exec.Command("bash", "-c", c)
	err = cmd.Run()
	if err != nil {
		return
	}
}

// FileSync
// @Description: 同步文件到远端
// @receiver this
// @param name
func (n *NotifyToSync) FileSync(name string) (res bool, err error) {
	isTmp := n.CheckIsTmpFile(name)
	if isTmp {
		return false, nil
	}
	remoteDir := filepath.Dir(name)
	if targetDir := n.Configs.Source.TargetDir; targetDir != "" {
		sourceDir := n.Configs.Source.SourceDir
		relativePath, err := filepath.Rel(sourceDir, remoteDir)
		if err == nil {
			remoteDir = filepath.Join(targetDir, relativePath)
		}
	}

	c := fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \" mkdir -p %s \"", n.SshPort, n.SshIdentify, n.SshUser, n.SshIp, remoteDir)
	log.Printf("创建远程目录：%s\n", remoteDir)
	cmd := exec.Command("bash", "-c", c)
	err = cmd.Run()
	if err != nil {
		return false, nil
	}
	// 目录则不执行同步
	fileInfo, err := os.Stat(name)
	if err != nil || fileInfo.IsDir() {
		return
	}
	filename := filepath.Base(name)
	remotePath := filepath.Join(remoteDir, filename)

	c = fmt.Sprintf("scp -P %s  -i %s  %s  %s@%s:%s", n.SshPort, n.SshIdentify, name, n.SshUser, n.SshIp, remotePath)
	log.Printf("执行复制：%s to %s\n", name, remotePath)
	cmd = exec.Command("bash", "-c", c)
	err = cmd.Run()
	if err != nil {
		return
	}
	n.ContentReplace(remotePath)
	return true, nil
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
	log.Printf("执行替换：%s\n", sedRule)
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
	for _, cmd := range command {
		c := fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \" %s \"", n.SshPort, n.SshIdentify, n.SshUser, n.SshIp, cmd)
		log.Printf("后置命令执行：%s\n", command)
		cmd := exec.Command("bash", "-c", c)
		err := cmd.Run()
		if err != nil {
			log.Println("后置命令执行失败")
		}
	}
}

//
// CheckIsTmpFile
// @Description: 检测是否是临时文件
// @receiver n
// @param name
// @return tmp
//
func (n *NotifyToSync) CheckIsTmpFile(name string) (tmp bool) {
	filename := filepath.Base(name)
	matched, err := regexp.MatchString(`\w+\.(swp|swx|swpx|\w+~)$`, filename)
	if err == nil && matched == true {
		return true
	}
	return false
}
