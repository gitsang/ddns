
default: build

build:

	go build -o bin/ddns cmd/ddns.go

pack: build 

	mkdir -p target/bin target/conf target/log
	cp bin/* target/bin/
	cp configs/* target/conf/

install:

	mkdir -p /usr/local/etc/ddns/bin /usr/local/etc/ddns/conf /usr/local/etc/ddns/log
	cp bin/ddns             /usr/local/bin/ddns
	cp configs/private.yml  /usr/local/etc/ddns/ddns.conf
	cp configs/ddns.service /etc/systemd/system/ddns.service
	systemctl status ddns

uninstall:

	rm -fr /usr/local/bin/ddns /usr/local/etc/ddns /etc/systemd/system/ddns.service

