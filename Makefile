dev: templ run

templ:
	templ generate

run:
	go run main.go

build:
	templ generate && go build -o ./tmp/main .