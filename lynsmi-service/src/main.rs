use lynsmi::Props;
use lynsmi::{errors::Error as SMIError, Lib, SMI};
use serde::Serialize;
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::TcpListener;
use tokio::sync::broadcast;
use tokio::task::spawn_blocking;

#[derive(Debug, Clone, Serialize)]
enum Event {
    InitError(String),
    AllProps(Vec<Option<Props>>),
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let (tx, _) = broadcast::channel(10);

    let tx_cloned = tx.clone();
    spawn_blocking(move || match Lib::try_default() {
        Err(e) => tx_cloned.send(Event::InitError(e.to_string())).unwrap(),
        Ok(lib) => match SMI::new(&lib) {
            Err(e) => tx_cloned.send(Event::InitError(e.to_string())).unwrap(),
            Ok(smi) => loop {
                let mut results = Vec::new();
                smi.get_devices(&mut results);
                tx_cloned
                    .send(Event::AllProps(
                        results.into_iter().map(|v| v.ok()).collect(),
                    ))
                    .unwrap();
            },
        },
    });

    let listener = TcpListener::bind("127.0.0.1:8080").await?;

    loop {
        let mut rx_cloned = tx.subscribe();
        let (mut socket, _) = listener.accept().await?;

        tokio::spawn(async move {
            while let Ok(event) = rx_cloned.recv().await {
                let b = serde_json::to_vec(&event).unwrap();
                async || -> _ {
                    socket.write_all(&b).await?;
                    socket.write_all(&[0]).await
                }()
                .await;
                if let Err(e) = socket.write_all(&b).await {
                    eprintln!("failed to write to socket; err = {:?}", e);
                    return;
                }
            }
        });
    }
}
