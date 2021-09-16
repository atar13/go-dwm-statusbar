export HOME := $(HOME)
export USER := $(shell logname)

.PHONY: all
all: build

.PHONY: build
build: 
	# go get -d src/*.go
	@echo "Building go-dwm-statusbar"
	go build -o go-dwm-statusbar ./src/*.go       

.PHONY: clean 	
clean:
	@echo "Removing go-dwm-statusbar binary"
	rm -f ./go-dwm-statusbar

.PHONY: install
install: 
	@echo "Installing go-dwm-statusbar to system at /usr/bin"
	cp ./go-dwm-statusbar /usr/bin
	@if [ -f "/home/$(USER)/.config/go-dwm-statusbar/config.yaml" ]; then \
		echo "Found exisisting config. Not overwriting changes"; \
	else \
		echo "Creating config directory at ~/.config/"; \
		mkdir -p /home/$(USER)/.config/go-dwm-statusbar; \
		echo "Copying config file to ~/.config/go-dwm-statusbar"; \
		cp config-sample.yaml /home/$(USER)/.config/go-dwm-statusbar/config.yaml; \
	fi

.PHONY: test
test:
	go run src/main.go

.PHONY: uninstall
uninstall:
	@echo "Uninstalling go-dwm-statusbar"
	rm -f /usr/bin/go-dwm-statusbar
	@echo "Removing config files"
	rm -rf /home/$(USER)/.config/go-dwm-statusbar/
