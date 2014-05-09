all:
	cd site-src; cactus build; cp -r .build/* ..; cd ..
	md5sum *.deb > md5sum
	python md52json.py md5sum > md5sum.json