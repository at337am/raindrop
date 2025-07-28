default:
    @just --list

cp:
    fd -e go -x sh -c 'echo "===== {} ====="; cat {}; echo' | wl-copy

tidy:
    @go mod tidy
    @echo "Dependencies have been tidied up"

clean:
    @echo "Cleaning up..."
    @rm -rfv release/
    @echo "Cleanup complete"

build:
    @mkdir -p release/
    @go build -o ./release/rdrop
    @echo "Build complete -> ./release/rdrop"

install: build
    @mv ./release/rdrop ~/go/bin
    @echo "Installation complete -> ~/go/bin"
    @just clean
