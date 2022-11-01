# ClipShare

---

[English](README.md) | 简体中文

![image clipshare](./docs/clipshare.png)

## 灵感来自 [欧神的midgard项目](https://github.com/changkun/midgard)

## 特点

1. 数据备份
2. 多特备共享剪贴板
3. 多用户支持

## 运行
### Linux/Windows/MacOS

> 服务端

```shell
1. git clone https://github.com/dmzlingyin/clipshare.git
2. cd clipshare
3. make server
4. ./clipshare server
```
> 客户端

```shell
1. git clone https://github.com/dmzlingyin/clipshare.git
2. cd clipshare
3. make client
4. ./clipshare client
```

### Android
> 客户端

要构建安卓客户端, 需使用[gomobile](https://golang.org/x/mobile), 详细信息 [GoMobile wiki](https://github.com/golang/go/wiki/Mobile).
```shell
gomobile build -v -target=android -androidapi 19 -o clipshare.apk cmd/gui/main.go
```

### **根据用户名、密码等个人配置，更新conf/server.yaml和conf/client.yaml**

## ToDo

- [x] 基本通信
- [x] 跨平台
- [ ] 数据备份
- [ ] 数据共享
- [ ] 多用户
