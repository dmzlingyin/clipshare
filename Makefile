# (C) 2022 Wenchao lv

server:
	go build -tags=server -ldflags="-s -w" -o clipshare main.go

server_debug:
	go build -tags=server -o clipshare main.go

client:
	go build -tags=client -ldflags="-s -w" -o clipshare main.go

install: client
	rm /usr/bin/clipshare
	ln -s /home/lingyin/go/my_src/clipshare/clipshare /usr/bin/clipshare
