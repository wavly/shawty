# Include the local env file
include .env

tailwind:
	@bunx tailwindcss -i static/init.css -o static/dist.css --watch

tailminify:
	@bunx tailwindcss -i static/init.css -o static/dist.css --minify

# Templ handles live-reloading out-of-the-box
server:
	@templ generate --watch --proxy="http://localhost:$(PORT)" --cmd="make run"

build:
	@go build

run:
	clear && go build && ./surf

sqlc:
	@sqlc generate
