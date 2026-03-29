build:
    go build -o zenth ./cmd/zenth

install:
    go install ./cmd/zenth

test:
    go test ./...

clean:
    rm -f zenth
