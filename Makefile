dev: templ css run

templ:
	templ generate

css:
	tailwindcss -i ./assets/tailwind.css -o ./assets/app.css

run:
	go run main.go

build:
	templ generate && tailwindcss -i ./assets/tailwind.css -o ./assets/app.css && go build -o ./tmp/main .