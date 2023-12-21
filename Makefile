dev: templ css run

templ:
	templ generate

css:
	tailwindcss -c ./assets/tailwind.config.js -i ./assets/css/app.css -o ./static/css/app.css

run:
	go run main.go

build:
	templ generate && tailwindcss -c ./assets/tailwind.config.js -i ./assets/css/app.css -o ./static/css/app.css && go build -o ./tmp/main .