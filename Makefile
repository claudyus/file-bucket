
build: file-bucket.go
	go build file-bucket.go

package: build
	sudo rm -rf build
	mkdir -p build/
	echo 2.0 > build/debian-binary
	echo "Package: file-bucket" > build/control
	echo "Version: `git describe --tags`" >> build/control
	echo "Vcs-Git: https://github.com/claudyus/file-bucket.git" >> build/control
	echo "Architecture: amd64" >> build/control
	echo "Section: net" >> build/control
	echo "Maintainer: Claudio Mignanti <c.mignanti@gmail.com>" >> build/control
	echo "Priority: optional" >> build/control
	echo "Description: No-auth backup system over HTTP" >> build/control
	echo " Built" `date`
	mkdir -p build/usr/local/bin
	mkdir -p build/etc/init
	mkdir -p build/etc/file-bucket
	cp file-bucket build/usr/local/bin/
	cp hooks/* build/etc/file-bucket/
	cp config.json build/etc/file-bucket/config.example.json
	sudo chown -R root: build/usr
	sudo chown -R root: build/etc
	tar cvzf build/data.tar.gz -C build usr etc
	tar cvzf build/control.tar.gz -C build control
	cd build && ar rc file-bucket.deb debian-binary control.tar.gz data.tar.gz && cd ..
