use notify::{Event, RecursiveMode, Watcher};

use std::{
  net::{TcpListener, TcpStream},
  path::Path,
  sync::mpsc,
  thread,
};

fn main() {
  println!("Listening on: 9090");
  println!("Watching the filesystem");
  let listener = TcpListener::bind("0.0.0.0:9090").unwrap();

  // Accept incoming connections on separate thread
  for stream in listener.incoming() {
    // Painc the main thread if error in the TcpStream
    thread::spawn(|| accept_conns(stream.expect("tcp stream")));
  }
}

fn accept_conns(stream: TcpStream) {
  // Upgrade the connection to websocket
  let mut websocket = tungstenite::accept(stream).expect("tcp stream accept");

  // Create channel for notify::Watcher to send/receive events
  let (tx, rx) = mpsc::channel::<Result<Event, notify::Error>>();
  let mut watcher = notify::recommended_watcher(tx).expect("watcher");

  // Watch the files in Recursive mode
  watcher
    .watch(Path::new("../"), RecursiveMode::Recursive)
    .expect("watching files");

  // Blocks forever
  for res in rx {
    match res {
      Ok(event) => {
        let kind = event.kind;

        // Only send a message whenever the files is modify/create/remove
        if kind.is_remove() || kind.is_modify() || kind.is_create() {
          if websocket.send("file changed".into()).is_err() {
            break;
          }

          // Try to read from the connection to prevent idle connections/threads
          let _ = websocket.read();
        }
      }
      Err(err) => eprintln!("Erorr receiving events: {err}"),
    }
  }
}
