koala2png=$(GOPATH)/bin/koala2png
hires2png=$(GOPATH)/bin/hires2png
png2koala=$(GOPATH)/bin/png2koala
png2hires=$(GOPATH)/bin/png2hires
vsfinject=$(GOPATH)/bin/vsfinject
mempetscii=$(GOPATH)/bin/mempetscii
prgmerge=$(GOPATH)/bin/prgmerge

default: all

all: $(koala2png) $(hires2png) $(png2koala) $(png2hires) $(vsfinject) $(mempetscii) $(prgmerge)

godeps:
	go get -d ./...

$(hires2png): cmd/hires2png.go pkg/gfx/*.go
	go build -o $@ $<

$(koala2png): cmd/koala2png.go pkg/gfx/*.go
	go build -o $@ $<

$(png2koala): cmd/png2koala.go pkg/gfx/*.go
	go build -o $@ $<

$(png2hires): cmd/png2hires.go pkg/gfx/*.go
	go build -o $@ $<

$(vsfinject): cmd/vsfinject.go pkg/file/snapshot.go
	go build -o $@ $<

$(mempetscii): cmd/mempetscii.go pkg/file/memory.go
	go build -o $@ $<

$(prgmerge): cmd/prgmerge.go pkg/file/program.go
	go build -o $@ $<
