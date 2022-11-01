# ClipShare

---

English | [简体中文](README_zh.md)

![image clipshare](./docs/clipshare.png)

## Inspired by [changkun's midgard project](https://github.com/changkun/midgard)

## Features

1. data backup
2. data share between more one devices
3. multiusers

## How To Run

### Linux/Windows/MacOS

> Server

```shell
1. git clone https://github.com/dmzlingyin/clipshare.git
2. cd clipshare
3. make server
4. ./clipshare server
```
> Client

```shell
1. git clone https://github.com/dmzlingyin/clipshare.git
2. cd clipshare
3. make client
4. ./clipshare client
```

### Android
> Client

To build it, one
must use [gomobile](https://golang.org/x/mobile). You may follow the instructions
provided in the [GoMobile wiki](https://github.com/golang/go/wiki/Mobile) page.
```shell
gomobile build -v -target=android -androidapi 19 -o clipshare.apk cmd/gui/main.go
```

### **Modify the conf/server.yaml or conf/client.yaml based on your username、password etc.**

## ToDo

- [x] fundamental communicaion
- [x] multi platforms
- [ ] data backup
- [ ] data share
- [ ] more than one users
