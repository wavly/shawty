const fs = require("node:fs");
const WebSocket = require('ws');

const wss = new WebSocket.Server({ port: 9090 });
console.info("Watching for file changes")

wss.on("connection", (ws) => {
  const watcher = fs.watch(".", { recursive: true })
  watcher.on("change", () => {
    ws.send("file changed")
  })

  ws.on("error", (err) => {
    console.error(err);
  });
});

