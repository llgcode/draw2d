
#include $(GOROOT)/src/Make.inc

all:	install test

install:
	go install
	go install ./draw2dgl
#	cd wingui && make install

build:
	go build
	go build ./draw2dgl
#	cd wingui && make build

test:
	cd cmd && go build draw2dgl.go
	cd cmd && go build gettingStarted.go
	cd cmd && go build testandroid.go
	cd cmd && go build testdraw2d.go
	cd cmd && go build testgopher.go
	cd cmd && go build testimage.go
	#cd cmd && go build testX11draw.go

clean:
	go clean ./...
#	cd wingui && make clean

command:
	cd cmd && make

fmt:
	gofmt -w . 

