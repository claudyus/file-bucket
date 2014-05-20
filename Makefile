all:
	python rst2html.py ../README.rst > site-src/pages/README.html
	md5sum *.deb > md5sum
	#build than copy
	cd site-src; cactus build; cd ..
	cp -r site-src/.build/* .

serve: all
	cd site-src/; cactus serve
