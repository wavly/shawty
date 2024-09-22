tailwind:
	@bunx tailwindcss -i static/init.css -o static/dist.css --watch

tailminify:
	@bunx tailwindcss -i static/init.css -o static/dist.css --minify

server:
	@watchexec -c -r -e go,html go run .

sqlc:
	@sqlc generate
