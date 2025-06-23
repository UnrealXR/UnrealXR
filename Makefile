APP_DIR := ./app
OUTPUT := uxr
TAGS := drm drm_leasing drm_disable_input

.PHONY: all build clean

all: build

build:
	go build -tags '$(TAGS)' -o $(OUTPUT) $(APP_DIR)

clean:
	rm -f $(OUTPUT)
