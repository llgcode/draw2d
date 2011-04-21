
#include $(GOROOT)/src/Make.inc

all:	install

install:
	cd draw2d && make install
	cd draw2dgl && make install
	cd postscript && make install
#	cd wingui && make install

clean:
	cd draw2d && make clean
	cd draw2dgl && make clean
	cd postscript && make clean
	cd cmd && make clean
#	cd wingui && make clean

nuke:
	cd draw2d && make nuke
	cd draw2dgl && make nuke
	cd postscript && make nuke
#	cd wingui && make nuke

command:
	cd cmd && make

fmt:
	gofmt -w . 

