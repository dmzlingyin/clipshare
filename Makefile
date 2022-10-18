# (C) 2022 Wenchao lv

all: build

build:
	go build -o clipshare main.go

install: build
	rm /usr/bin/clipshare
	ln -s /home/lingyin/go/my_src/clipshare/clipshare /usr/bin/clipshare
