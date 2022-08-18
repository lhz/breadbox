koala2png=bin/koala2png
hires2png=bin/hires2png
png2koala=bin/png2koala
png2hires=bin/png2hires
vsfinject=bin/vsfinject
mempetscii=bin/mempetscii
prgmerge=bin/prgmerge

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
