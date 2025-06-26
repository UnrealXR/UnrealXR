APP_DIR := ./app
OUTPUT := uxr
TAGS := xreal noaudio drm drm_leasing drm_disable_input

.PHONY: all build clean

all: build

build:
	go build -v -tags '$(TAGS)' -o $(OUTPUT) $(APP_DIR)

clean:
	rm -f $(OUTPUT)
