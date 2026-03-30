build:
    go build -o zenth ./cmd/zenth

install:
    go install ./cmd/zenth

test:
    go test ./...

clean:
    rm -f zenth

docs:
    hugo -s site -D

docs-serve:
    hugo server -s site -D
