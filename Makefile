dev: templ css app

templ:
	templ generate

css:
	tailwindcss -i ./assets/tailwind.css -o ./assets/app.css

app:
	go run main.go