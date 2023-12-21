# hackernews

A Hacker News clone built with [Go](https://go.dev/), [templ](https://github.com/a-h/templ), and [Tailwind](https://tailwindcss.com/)

<p align="center">
    <img width="1080" src="https://hckrnws.fly.dev/static/img/homepage.png">
    <img width="1080" src="https://hckrnws.fly.dev/static/img/thread.png">
</p>

## Demo

https://hckrnws.fly.dev/news

## Features

- No JS required

## Build Setup

#### Install tailwind cli
```bash
brew install tailwindcss
```

#### Install templ
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

#### Build and run server
```bash
make dev
```

#### Deploy to [Fly](https://fly.io)
```bash
fly launch
```