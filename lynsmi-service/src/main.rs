use std::{
    net::SocketAddr,
    sync::{Arc, Condvar, Mutex},
};

use lynsmi::Props;
use lynsmi_service::prelude::*;
use tokio::{
    net::{TcpListener, TcpStream},
    select, spawn,
    sync::broadcast::{channel, error::RecvError, Receiver, Sender},
    task::spawn_blocking,
};
use tracing::{info, warn};

async fn handle(socket: TcpStream, addr: SocketAddr, mut rx: Receiver<Data>) {
    let mut conn = Connection::new(socket);
    info!("accept client {}", addr);
    loop {
        match rx.recv().await {
            Ok(v) => {
                let id = v.0;
                let (props, err) = match v.1.as_ref() {
                    Err(e) => {
                        warn!("device {} err {:?}", id, e);
                        (None, Some(e.to_string()))
                    }
                    Ok(v) => (Some(v.to_owned()), None),
                };
                if let Err(e) = conn.send(&PropsWithID { id, props, err }).await {
                    warn!("failed to send to client; addr = {}; err = {:?}", addr, e);
                    return;
                }
            }
            Err(RecvError::Lagged(v)) => warn!("lagged {}", v),
            Err(RecvError::Closed) => break,
        }
    }
}

async fn listen(listener: TcpListener, tx: Sender<Data>, notify: Arc<(Mutex<bool>, Condvar)>) {
    let (lock, cvar) = &*notify;
    loop {
        match listener.accept().await {
            Err(e) => warn!("failed to accept client {}", e),
            Ok((socket, addr)) => {
                let rx = tx.subscribe();
                {
                    let _l = lock.lock().unwrap();
                    cvar.notify_all();
                }
                tokio::spawn(handle(socket, addr, rx));
            }
        };
    }
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    const ADDR: &'static str = "0.0.0.0:5432";
    info!("listen on {}", ADDR);
    let listener = TcpListener::bind(ADDR).await?;

    let (tx, _) = channel::<Arc<(usize, lynsmi::Result<Props>)>>(100);
    let notify = Arc::new((Mutex::new(false), Condvar::new()));
    select!(
        _ = spawn(listen(listener, tx.clone(), notify.clone())) => {},
        result = spawn_blocking(move || watch(tx, notify)) => result??
    );

    Ok(())
}
