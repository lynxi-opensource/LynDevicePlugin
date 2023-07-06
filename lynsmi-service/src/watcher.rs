use std::{
    sync::{Arc, Condvar, Mutex},
    thread,
};

use lynsmi::{Lib, Props, Result, Symbols};
use tokio::sync::broadcast::Sender;
use tracing::{info, warn};

pub type Data = Arc<(usize, lynsmi::Result<Props>)>;

pub fn watch(tx: Sender<Data>, notify: Arc<(Mutex<bool>, Condvar)>) -> Result<()> {
    let lib = Lib::try_default()?;
    let smi = Symbols::new(&lib)?;
    let device_cnt = smi.get_device_cnt()?;
    thread::scope(|s| {
        for id in 0..device_cnt {
            let smi = smi.clone();
            let tx = tx.clone();
            let notify = notify.clone();
            s.spawn(move || {
                let (lock, cvar) = &*notify;
                loop {
                    {
                        let mut quit = lock.lock().unwrap();
                        if *quit {
                            break;
                        }
                        while tx.receiver_count() == 0 {
                            info!("device {} waiting receiver", id);
                            quit = cvar.wait(quit).unwrap();
                            if *quit {
                                break;
                            }
                        }
                    }
                    let result = smi.get_props(id);
                    if let Err(_) = tx.send(Arc::new((id, result))) {
                        warn!("SendError");
                    };
                }
            });
        }
    });
    Ok(())
}
