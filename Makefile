.PHONY: build
build:
	export GO111MODULE=on
	rm -rf build/*
	go mod download
	go build -o build/journey
	cp -R built-in build/built-in
	mkdir -p build/content/themes/promenade
	cp -R content/themes/promenade build/content/themes/promenade