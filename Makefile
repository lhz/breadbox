koala2png=$(GOPATH)/bin/koala2png
hires2png=$(GOPATH)/bin/hires2png
vsfinject=$(GOPATH)/bin/vfsinject
prgmerge=$(GOPATH)/bin/prgmerge

default: all

all: $(koala2png) $(hires2png) $(vsfinject) $(prgmerge)

godeps:
	go get -d ./...

$(hires2png): cmd/hires2png.go gfx/*.go
	go build -o $@ $<

$(koala2png): cmd/koala2png.go gfx/*.go
	go build -o $@ $<

$(vsfinject): cmd/vsfinject.go file/snapshot.go
	go build -o $@ $<

$(prgmerge): cmd/prgmerge.go file/program.go
	go build -o $@ $<
