dev: templ css app

templ:
	templ generate

css:
	bun run build

app:
	go run main.go