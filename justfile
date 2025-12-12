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
    @go build -o ./release/rd
    @echo "Build complete -> ./release/rd"

install: build
    @mkdir -p ~/go/bin
    @mv ./release/rd ~/go/bin
    @echo "Installation complete -> ~/go/bin"
    @just clean
