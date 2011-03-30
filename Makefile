
#include $(GOROOT)/src/Make.inc

all:	install

install:
	cd draw2d && make install
	cd postscript && make install
	cd wingui && make install

clean:
	cd draw2d && make clean
	cd postscript && make clean
	cd wingui && make clean

nuke:
	cd draw2d && make nuke
	cd postscript && make nuke
	cd wingui && make nuke

fmt:
	gofmt -w draw2d postscript wingui cmd
