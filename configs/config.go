package configs

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

type TomlConfig struct {
	Source      Source      `toml:"source"`
	Destination Destination `toml:"destination"`
	ReplaceRule ReplaceRule `toml:"replace_rule"`
	Command     Command     `toml:"command"`
}

type Source struct {
	SourceDir string `toml:"source_dir"`
	TargetDir string `toml:"target_dir"`
	DeleteDir string `toml:"delete_dir"`
}

type Destination struct {
	SshIp          string `toml:"ssh_ip"`
	SshPort        string `toml:"ssh_port"`
	SshUser        string `toml:"ssh_user"`
	SshIdentity    string `toml:"ssh_identify"`
	DestinationDir string `toml:"destination_dir"`
}

type Command struct {
	CompletedCmd []string `toml:"completed_cmd"`
}

type ReplaceRule struct {
	SedRule string `toml:"sed_rule"`
}

var Config = &TomlConfig{}

//
// ReadConf
// @Description: 读取配置文件
// @param path
// @return p
// @return err
//
func ReadConf(path string) (p *TomlConfig, err error) {
	fcontent := loadToml(path)
	if fcontent == nil {
		return
	}
	p = new(TomlConfig)

	if err = toml.Unmarshal(fcontent, p); err != nil {
		fmt.Println("toml.Unmarshal error ", err)
		return
	}
	setConfig(p)
	return
}

//
// setConfig
// @Description: 设置配置文件
// @param config
//
func setConfig(config *TomlConfig) {
	Config = config
}

//
// loadToml
// @Description: 加载toml配置文件
// @param path
// @return fcontent
//
func loadToml(path string) (fcontent []byte) {
	var (
		fp  *os.File
		err error
	)
	if _, err := os.Stat(path); err != nil {
		fmt.Println("file not exist", err)
		return nil
	}
	if fp, err = os.Open(path); err != nil {
		fmt.Println("open error ", err)
		return nil
	}

	if fcontent, err = ioutil.ReadAll(fp); err != nil {
		fmt.Println("ReadAll error ", err)
		return nil
	}
	return fcontent
}
