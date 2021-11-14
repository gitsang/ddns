
default: build

build:

	go build -o bin/ddns cmd/ddns.go

pack: build 

	mkdir -p target/bin target/conf target/log
	cp bin/* target/bin/
	cp configs/* target/conf/

install:

	cp bin/ddns /usr/local/bin/ddns
	mkdir -p /usr/local/etc/ddns
	cp configs/private.yml /usr/local/etc/ddns/ddns.conf
	touch /var/log/ddns.log
	cp configs/ddns.service /etc/systemd/system/ddns.service
	systemctl status ddns

uninstall:

	rm -fr /usr/local/bin/ddns /usr/local/etc/ddns /etc/systemd/system/ddns.service

