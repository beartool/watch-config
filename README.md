# watch-config：是一个基于fsnotify开发的文件监控工具，检测文件变化，并同步到其他远程节点

## 目录

* [支持环境](#支持环境)
* [编译和安装](#编译和安装)
* [配置文件说明](#配置文件说明)



## 支持环境
目前是在linux环境下运行的

## 编译和安装
可下载源码包进行编译，或者直接下载watcher-file即可使用

## 配置文件说明

本程序使用

```toml
[source]
source_dir = "/home/www/"
target_dir = "/root/"

[destination]
ssh_ip = "47.24.10.249"
ssh_port = "22333"
ssh_user = "www"
ssh_identify = "～/.ssh/dientify.pem"

[command]
completed_cmd = "touch /root/b.text"

[replace_rule]
sed_rule = "'s/1e15849033-laa84/1149b49065-rku24/g'"
```


基础属性说明

| 属性           | 类型     | 作用                                        |
|--------------|--------|-------------------------------------------|
| source_dir | string | 设置当前程序监控的文件目录                             |
| target_dir | string | 设置同步目标路径，根据source_dir计算相对路径，再拼接target_dir |
| ssh_ip | string | 指定远程；连接地址                                 |
| ssh_port | string   | 远程连接端口                                    |
| ssh_user | string    | 远程同步用户                                    |
| ssh_identify | string   | 远程连接授权秘钥                                  |
| completed_cmd | string   | 完成同步后执行远端命令，比如重启某个服务                      |
| sed_rule | string   | 设置同步远端后 文件内容搜索替换规则                        |



