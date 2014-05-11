
build: file-bucket.go
	go build file-bucket.go

debian: build
	sudo rm -rf build/
	mkdir -p build/usr/local/bin
	mkdir -p build/etc/init
	mkdir -p build/etc/file-bucket
	cp file-bucket build/usr/local/bin/
	cp hooks/* build/etc/file-bucket/
	cp config.json build/etc/file-bucket/config.example.json
	sudo chown -R root: build/usr
	sudo chown -R root: build/etc
	# copy and modify debian package info
	cp debian/* build
	sed -i "s/SIZE/`du build/ | tail -1 | cut -f 1`/g" build/control
	sed -i "s/VERSION/`git describe --tags`/g" build/control
	sed -i "s/DATE/`date`/g" build/control
	# build the deb file
	tar cvzf build/data.tar.gz -C build usr etc
	tar cvzf build/control.tar.gz -C build control
	cd build && ar rc file-bucket.deb debian-binary control.tar.gz data.tar.gz && cd ..
	mv build/file-bucket.deb gh-pages/file-bucket_`git describe --tags`.deb

site: debian
	make -C gh-pages/