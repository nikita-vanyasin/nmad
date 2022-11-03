OUT_DIR := $(shell pwd)/build

clean:
	rm -rf build/dist
	rm build/nmad-app

build-prod: clean
	mkdir -p "${OUT_DIR}"
	cd app && CGO_ENABLED=0 go build -ldflags "-s -w" -o "${OUT_DIR}/nmad-app"
	cd map && ./run.sh npm run build-only
	cp -r map/dist "${OUT_DIR}/dist"