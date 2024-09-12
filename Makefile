tailwind:
	@bunx tailwindcss -i static/init.css -o static/dist.css --watch

server:
	@watchexec -c -r -e go,html go run .
