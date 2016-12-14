koala2png=$(GOPATH)/bin/koala2png
hires2png=$(GOPATH)/bin/hires2png

default: all

all: $(koala2png) $(hires2png)

godeps:
	go get -d ./...

$(hires2png): cmd/hires2png.go gfx/*.go
	go build -o $@ $<

$(koala2png): cmd/koala2png.go gfx/*.go
	go build -o $@ $<

$(vsfinject): cmd/vsfinject.go file/snapshot.go
	go build -o $@ $<

$(prgmerge): cmd/prgmerge file/program.go
	go build -o $@ $<
