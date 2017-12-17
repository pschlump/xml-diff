
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

test03:


