use notify::{Event, RecursiveMode, Watcher};

use std::{net::TcpListener, path::Path, sync::mpsc, thread};

fn main() {
  println!("Listening on: 9090");
  println!("Watching the filesystem");
  let server = TcpListener::bind("0.0.0.0:9090").unwrap();

  for stream in server.incoming() {
    thread::spawn(|| {
      let stream = stream.expect("tcp stream");
      let mut websocket = tungstenite::accept(stream).expect("tcp stream accept");
      let (tx, rx) = mpsc::channel::<Result<Event, notify::Error>>();

      let mut watcher = notify::recommended_watcher(tx).expect("watcher");
      watcher
        .watch(Path::new("../../"), RecursiveMode::Recursive)
        .expect("watching files");

      // Blocks forever
      for res in rx {
        match res {
          Ok(event) => {
            let kind = event.kind;
            if kind.is_remove() || kind.is_modify() || kind.is_create() {
              if websocket.send("file changed".into()).is_err() {
                break;
              }
            }
          }
          Err(err) => eprintln!("Erorr receiving events: {err}"),
        }
      }
    });
  }
}
