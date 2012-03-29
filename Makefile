
#include $(GOROOT)/src/Make.inc

all:	install

install:
	cd draw2d && go install
#	cd draw2dgl && make install
	cd postscript && go install
#	cd wingui && make install

build:
	cd draw2d && go build
#	cd draw2dgl && make build
	cd postscript && go build
#	cd wingui && make build

clean:
	cd draw2d && go clean
#	cd draw2dgl && make clean
	cd postscript && go clean
	cd cmd && go clean
#	cd wingui && make clean

command:
	cd cmd && make

fmt:
	gofmt -w . 

