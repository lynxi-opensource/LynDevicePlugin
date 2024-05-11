use std::{collections::HashMap, ffi::c_int, slice, sync::Mutex};

use lazy_static::lazy_static;
use lyndriver::{
    drv::{self, ErrMsg, RawErrMsg},
    Result,
};
use tracing::error;

lazy_static! {
    pub static ref DRV_EXCEPTION_MAP: Mutex<HashMap<u32, ErrMsg>> = Mutex::new(HashMap::new());
}

fn cb(raw_err_msgs: *mut RawErrMsg, cnt: *mut c_int) -> c_int {
    let cnt = unsafe { *cnt };
    let raw_err_msg_slice = unsafe { slice::from_raw_parts(raw_err_msgs, cnt as usize) };
    let mut exception_map = DRV_EXCEPTION_MAP.lock().unwrap();
    for raw in raw_err_msg_slice {
        match ErrMsg::try_from(raw) {
            Ok(err_msg) => {
                exception_map.insert(err_msg.devid, err_msg);
            }
            Err(e) => {
                error!("convert err {}", e);
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
    *DRV_EXCEPTION_MAP.lock().unwrap() = exception_map;
    drv_exception_symbols.get_device_current_exception(cb)?;
    Ok(())
}
