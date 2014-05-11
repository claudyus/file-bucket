all:
	wget https://github.com/claudyus/file-bucket/blob/master/README.md -O - | sed '/.*csrf.*/d' - > site-src/static/README.html
	cd site-src; cactus build; cd ..
	md5sum *.deb > site-src/static/md5sum
	python md52json.py md5sum > site-src/static/md5sum.json
	cp -r site-src/.build/* `pwd`

serve: all
	cd site-src; cactus serve
