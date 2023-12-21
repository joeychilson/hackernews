dev: templates css app

templates:
	templ generate

css:
	bun run build

app:
	go run main.go