
#include $(GOROOT)/src/Make.inc

all:	install test

install:
	go install ./...

build:
	go build ./...

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

fmt:
	gofmt -w . 

