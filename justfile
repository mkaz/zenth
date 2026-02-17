build:
    go build -o zenth ./cmd/zenth
    cp ./zenth ~/bin/

install:
    go install ./cmd/zenth

test:
    go test ./...

clean:
    rm -f zenth
