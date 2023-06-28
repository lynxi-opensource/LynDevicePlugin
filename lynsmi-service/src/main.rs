use lynsmi_service::prelude::*;
use tokio::net::TcpListener;
use tracing::{info, warn};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    const ADDR: &'static str = "127.0.0.1:5432";
    info!("listen on {}", ADDR);
    let listener = TcpListener::bind(ADDR).await?;

    let smi_watcher = SMIWatcher::new().await?;
    loop {
        let (socket, addr) = listener.accept().await?;
        let mut smi = smi_watcher.subscribe();

        tokio::spawn(async move {
            let mut conn = Connection::new(socket);
            info!("accept client {}", addr);
            while smi.changed().await.is_ok() {
                let all_props: Vec<_> = smi
                    .borrow()
                    .iter()
                    .map(|v| match v.as_ref() {
                        Err(e) => {
                            warn!("device err {:?}", e);
                            None
                        }
                        Ok(v) => Some(v.to_owned()),
                    })
                    .collect();
                if let Err(e) = conn.send(&all_props).await {
                    warn!("failed to send to client; addr = {}; err = {:?}", addr, e);
                    return;
                }
            }
        });
    }
}
