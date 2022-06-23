package cmd

import (
	"fmt"
	"log"
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

func (this *NotifyToSync) Notify(name string) {
	destination := configs.Config.Destination
	identifyFile := destination.SshIdentity
	sshIp := destination.SshIp
	sshPort := destination.SshPort
	sshUser := destination.SshUser
	sedRule := configs.Config.ReplaceRule.SedRule

	filename := filepath.Base(name)
	extname := path.Ext(filename)
	if identifyFile == "" || sshPort == "" || sshIp == "" || sshUser == "" {
		log.Fatalf("远程同步服务器配置缺失")
	}
	if extname == ".swp" {
		log.Printf("临时文件不同步 %s", name)
		return
	}

	dir := filepath.Dir(name)

	c := fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \" mkdir -p %s \"", sshPort, identifyFile, sshUser, sshIp, dir)
	log.Printf("创建目录：%s\n", dir)
	cmd := exec.Command("bash", "-c", c)
	cmd.Run()

	c = fmt.Sprintf("scp -P %s  -i %s  %s  %s@%s:%s", sshPort, identifyFile, name, sshUser, sshIp, name)
	log.Printf("执行复制：%s\n", c)
	cmd = exec.Command("bash", "-c", c)
	cmd.Run()

	c = fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \"sed -i %s %s \"", sshPort, identifyFile, sshUser, sshIp, sedRule, name)
	log.Printf("执行替换：%s\n", c)
	cmd = exec.Command("bash", "-c", c)
	cmd.Run()

	c = fmt.Sprintf("ssh -p %s  -i %s -l %s  %s  \"/usr/local/nginx/sbin/nginx -s reload \"", sshPort, identifyFile, sshUser, sshIp)
	log.Printf("启动服务：%s\n", c)
	cmd = exec.Command("bash", "-c", c)
	cmd.Run()
}
