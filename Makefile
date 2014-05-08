all:
	cd site-src; cactus build; cp -r .build/* ..; cd ..