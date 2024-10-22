const socketUrl = 'ws://localhost:{{ . }}/dev';

let socket = new WebSocket(socketUrl);

socket.addEventListener('close', () => {
  const interAttemptTimeoutMilliseconds = 100;
  const maxDisconnectedTimeMilliseconds = 3000;
  const maxAttempts = Math.round(
    maxDisconnectedTimeMilliseconds / interAttemptTimeoutMilliseconds,
  );

  let attempts = 0;
  const reloadIfCanConnect = () => {
    attempts++;
    if (attempts > maxAttempts) {
      console.error('Could not reconnect to dev server.');
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
});
