
#include $(GOROOT)/src/Make.inc

all:	install test

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

test:
	#cd cmd && go build draw2dgl.go
	cd cmd && go build gettingStarted.go
	cd cmd && go build testandroid.go
	cd cmd && go build testdraw2d.go
	cd cmd && go build testgopher.go
	cd cmd && go build testimage.go
	cd cmd && go build testpostscript.go
	#cd cmd && go build testX11draw.go

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

