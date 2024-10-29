const socketUrl = "ws://127.0.0.1:9090";
const socket = new WebSocket(socketUrl);

socket.onopen = () => {
  console.info("Connected to WebSocket server");
};

socket.onmessage = (event) => {
  // Reload the page whenever receving a specific message
  if (event.data == "file changed") {
    console.log("reloading");
    location.reload();
  }
};

socket.onerror = (error) => {
  console.error("WebSocket error:", error);
  console.log("Trying to reconnect");

  const interAttemptTimeoutMilliseconds = 100;
  const maxDisconnectedTimeMilliseconds = 3000;
  const maxAttempts = Math.round(
    maxDisconnectedTimeMilliseconds / interAttemptTimeoutMilliseconds,
  );

  let attempts = 0;
  function reloadIfCanConnect() {
    attempts++;
    if (attempts > maxAttempts) {
      console.error('Could not reconnect to dev server');
      return;
    }

    socket = new WebSocket(socketUrl);
    socket.addEventListener('error', () => {
      setTimeout(reloadIfCanConnect, interAttemptTimeoutMilliseconds);
    });

    socket.addEventListener('open', () => {
      location.reload();
    });
  };

  reloadIfCanConnect();
};

socket.addEventListener('close', () => {
  console.info("WebSocket connection closed!");
});
