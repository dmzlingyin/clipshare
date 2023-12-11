# (C) 2022 Wenchao lv
client:
	go build -ldflags="-s -w" -o clipshare main.go

install: client
	rm /usr/bin/clipshare
	ln -s /home/lingyin/go/my_src/clipshare/clipshare /usr/bin/clipshare
