use std::sync::{Arc, Condvar, Mutex};

use lynsmi::{AllProps, Lib, SMI};
use tokio::{
    sync::{oneshot, watch},
    task::{spawn_blocking, JoinError, JoinHandle},
};

pub struct SMIWatcher {
    rx: watch::Receiver<Arc<AllProps>>,
    notify: Arc<(Mutex<bool>, Condvar)>,
    h: JoinHandle<()>,
}

impl SMIWatcher {
    pub async fn new() -> lynsmi::Result<Self> {
        let (tx, rx) = watch::channel(Arc::new(AllProps::new()));
        let notify = Arc::new((Mutex::new(false), Condvar::new()));
        let notify_clone = notify.clone();
        let (init_tx, init_rx) = oneshot::channel::<lynsmi::Result<()>>();
        let h = spawn_blocking(move || match Lib::try_default() {
            Err(e) => {
                init_tx.send(Err(e)).unwrap();
                return;
            }
            Ok(lib) => match SMI::new(&lib) {
                Err(e) => {
                    init_tx.send(Err(e)).unwrap();
                    return;
                }
                Ok(smi) => {
                    init_tx.send(Ok(())).unwrap();
                    let (lock, cvar) = &*notify_clone;
                    loop {
                        {
                            let mut quit = lock.lock().unwrap();
                            if *quit {
                                break;
                            }
                            while tx.receiver_count() < 2 {
                                quit = cvar.wait(quit).unwrap();
                                if *quit {
                                    break;
                                }
                            }
                        }
                        let mut results = Vec::new();
                        smi.get_devices(&mut results);
                        tx.send(Arc::new(results)).unwrap();
                    }
                }
            },
        });
        init_rx.await.unwrap()?;
        Ok(Self { rx, notify, h })
    }

    pub fn subscribe(&self) -> watch::Receiver<Arc<AllProps>> {
        let rx = self.rx.clone();
        let __ = self.notify.0.lock().unwrap();
        self.notify.1.notify_one();
        rx
    }

    pub async fn close(self) -> std::result::Result<(), JoinError> {
        let mut quit = self.notify.0.lock().unwrap();
        *quit = true;
        self.notify.1.notify_one();
        self.h.await
    }
}
