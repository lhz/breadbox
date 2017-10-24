koala2png=$(GOPATH)/bin/koala2png
hires2png=$(GOPATH)/bin/hires2png
vsfinject=$(GOPATH)/bin/vsfinject
mempetscii=$(GOPATH)/bin/mempetscii
prgmerge=$(GOPATH)/bin/prgmerge

default: all

all: $(koala2png) $(hires2png) $(vsfinject) $(mempetscii) $(prgmerge)

godeps:
	go get -d ./...

$(hires2png): cmd/hires2png.go gfx/*.go
	go build -o $@ $<

$(koala2png): cmd/koala2png.go gfx/*.go
	go build -o $@ $<

$(vsfinject): cmd/vsfinject.go file/snapshot.go
	go build -o $@ $<

$(mempetscii): cmd/mempetscii.go file/memory.go
	go build -o $@ $<

$(prgmerge): cmd/prgmerge.go file/program.go
	go build -o $@ $<
