# hackernews

This project is aimed at recreating the popular Hacker News using Go and the templ library. 

It's designed to closely mirror the functionality of hacker news as much with the limitations of the public API.

The application is a close drop-in replacement for the actual hacker news application with the limitations of the API. 

You can just replace the hacker news url https://news.ycombinator.com/ with the url https://hckrnws.fly.dev/. 

Reply links will redirect to the actual hacker news website.

### requirements 

#### install tailwind cli
```bash
brew install tailwindcss
```

#### install templ
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

#### run dev
```bash
make dev
```

#### deploy
```bash
fly launch
```