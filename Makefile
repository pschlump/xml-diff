
all:
	go build

test: build test01 test02
	@echo PASS

build: 
	go build

test01:
	./xml-diff -l ./testdata/left.xml -r ./testdata/right.xml -lo ./out/lc.xml -ro ./out/rc.xml  >out/test01.out
	diff out/lc.xml ref
	diff out/rc.xml ref
	diff out/test01.out ref

test02:
	./xml-diff -byLine -l ./testdata/left.xml -r ./testdata/right.xml -lo ./out/lc02.xml -ro ./out/rc02.xml  >out/test02.out
	diff out/lc02.xml ref
	diff out/rc02.xml ref
	diff out/test02.out ref

test03: build
	./xml-diff -byLine -l ./testdata/l_t03.xml -r ./testdata/r_t03.xml -lo ./out/lc03.xml -ro ./out/rc03.xml  -lcfg ./testdata/lcfg.json

a03: build
	./xml-diff -byLine -l ./testdata/l_t03a.xml -r ./testdata/r_t03a.xml -lo ./out/lc03.xml -ro ./out/rc03.xml  -lcfg ./testdata/lcfg.json

a04: build
	./xml-diff -byLine -l ./testdata/l_t04a.xml -r ./testdata/r_t04a.xml -lo ./out/lc04.xml -ro ./out/rc04.xml  -lcfg ./testdata/lcfg_04.json -rcfg ./testdata/lcfg_04.json


install:
	( cd ~/bin ; rm -f xml-diff )
	( cd ~/bin ; ln -s ../go/src/github.com/pschlump/xml-diff/xml-diff . )

