dev: templ css app

templ:
	templ generate

css:
	tailwindcss -i ./assets/app.css -o ./dist/app.css

app:
	go run main.go