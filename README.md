# 永不摆烂队-简易抖音项目

## 项目依赖要求
```
go >= 1.18
mysql >= 8.0
ffmpeg >= 4.2.4
```

## 如何运行
1. 将项目下载至本地
```bash
git clone https://gitee.com/Whitroom/imitate-tiktok.git
cd imitate-tiktok
```
2. 将`confs/database copy.json`复制一份至原文件下，并重命名为`database.json`，输入数据库账号密码及数据库名称。

3. 通过`ipconfig`或`ifconfig`修改`model_change.go`中第61行`play_url`的`your_machine_ip`为机器在路由器中的ip地址

4. (1) 如果使用`vscode`，通过`vscode`进入该目录，点击`main.go`后，直接按下`F5`即可运行。

    (2) 如果使用命令行，输入
    ```bash
    go build && ./imitate-tiktok
    ```
    即可运行。

其他仍在补充...