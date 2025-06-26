APP_DIR := ./app
OUTPUT := ./uxr
TAGS := xreal noaudio drm drm_leasing drm_disable_input

.PHONY: all build clean

all: build

build:
	cd $(APP_DIR) && go build -v -tags '$(TAGS)' -o ../$(OUTPUT) .

clean:
	rm -f $(OUTPUT)
