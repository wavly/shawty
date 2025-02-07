# Include the local env file
include .env.local

BUILD_DIR = ./build

tailwind:
	@bunx tailwindcss -i static/init.css -o static/dist.css --watch

tailminify:
	@bunx tailwindcss -i static/init.css -o static/dist.css --minify

# Templ handles live-reloading out-of-the-box
server:
	@templ generate --watch --proxy="http://localhost:$(PORT)" --cmd="make run"

build | $(BUILD_DIR):
	clear && go build -o ./build/surf

$(BUILD_DIR):
	mkdir -p $@

run:
	clear && go build -o ./build/surf && ./build/surf

sqlc:
	@sqlc generate
