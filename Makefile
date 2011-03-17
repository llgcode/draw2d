
include $(GOROOT)/src/Make.inc

all:	install

install:
	cd draw2d && make install
	cd postscript && make install

clean:
	cd draw2d && make clean
	cd postscript && make clean

nuke:
	cd draw2d && make nuke
	cd postscript && make nuke
