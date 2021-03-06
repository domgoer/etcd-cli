# ETCD-CLI

## 使用etcd-cli连接etcd服务之后，可以使用常用的linux命令来操作etcd中的数据

```bash
etcd-cli -s 127.0.0.1 -p 2380
```

### 支持命令

- cd
- ls
- mkdir
- touch
- rm
- mv
- cp
- pwd
- cat

`（暂时不支持使用额外参数，如 -f、-r、-p等）`

**你甚至可以使用`vim`来修改etcd中可以被翻译成文本的文件**

需要注意的是`rm`、`mv`、`cp`这些命令在操作时需要在后加上 **"/"**,用来区分是文件夹还是文件

```bash
# 删除整个dir文件夹
rm dir/
# 删除file文件
rm file
```

mv、cp 同理

![option](./images/option.gif)

### 额外支持的命令

- upload [etcd-path] [local-path] 上传本地**文件**到etc指定的路径
- download [etcd-path] [local-path] 下载etcd中指定的**文件**到本地

如果在连接状态下使用upload或者download，local-path需要写绝对路径

也可以直接使用etcd-cli，这时local-path可以使用相对路径

```bash
etcd-cli -s 127.0.0.1 -p 2380 download /etcd-path/testfile ./
```

### 如何安装

安装go语言环境

```bash
go get github.com/domgoer/etcd-cli
cd $GOPATH/bin
./etcd-cli
```
