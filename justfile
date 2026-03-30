build:
    go build -o zenth ./cmd/zenth

install:
    go install ./cmd/zenth

test:
    go test ./...

clean:
    rm -f zenth

docs:
    git submodule update --init --recursive site/themes/hugo-book
    hugo -s site -D

docs-serve:
    git submodule update --init --recursive site/themes/hugo-book
    hugo server -s site -D --baseURL http://localhost:1313/ --appendPort=false
