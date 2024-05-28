use std::{collections::HashMap, ffi::c_int, slice, sync::Mutex};

use lazy_static::lazy_static;
use lyndriver::{
    drv::{self, ErrMsg, RawErrMsg},
    Result,
};
use tracing::{error, info};

lazy_static! {
    pub static ref DRV_EXCEPTION_MAP: Mutex<HashMap<u32, ErrMsg>> = Mutex::new(HashMap::new());
}

fn cb(raw_err_msgs: *mut RawErrMsg, cnt: *mut c_int) -> c_int {
    let cnt = unsafe { *cnt };
    let raw_err_msg_slice = unsafe { slice::from_raw_parts(raw_err_msgs, cnt as usize) };
    let mut exception_map = DRV_EXCEPTION_MAP
        .lock()
        .expect("lock DRV_EXCEPTION_MAP failed");
    for raw in raw_err_msg_slice {
        match ErrMsg::try_from(raw) {
            Ok(err_msg) => {
                if err_msg.enable_recover == 1 {
                    exception_map.remove(&err_msg.devid);
                    info!("device {} recovery: {:?}", &err_msg.devid, &err_msg);
                } else {
                    info!("device {} exception: {:?}", &err_msg.devid, &err_msg);
                    exception_map.insert(err_msg.devid, err_msg);
                }
            }
            Err(e) => {
                error!("parse err msg failed: {}", e);
            }
        };
    }
    0
}

pub fn listen() -> Result<()> {
    let drv_lib = lyndriver::drv::Lib::try_default()?;
    let drv_exception_symbols = drv::ExceptionSymbols::new(&drv_lib)?;
    let exceptions = drv_exception_symbols.get_device_exception()?;
    let mut exception_map = HashMap::new();
    for e in exceptions {
        exception_map.insert(e.devid, e);
    }
    info!("get device exception: {:?}", &exception_map);
    *DRV_EXCEPTION_MAP.lock().unwrap() = exception_map;
    drv_exception_symbols.get_device_current_exception(cb)?;
    Ok(())
}
