.PHONY: build install clean

build:
	go build -o ts .

install: build
	mkdir -p ~/bin
	cp ts ~/bin/ts

clean:
	rm -f ts
