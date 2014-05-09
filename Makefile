all:
	cd site-src; cactus build; cd ..
	md5sum *.deb > site-src/static/md5sum
	python md52json.py md5sum > site-src/static/md5sum.json
	cp -r site-src/.build/* .