use crate::errors::*;
use crate::ffi_convert::*;
use libloading::{Library, Symbol};
use serde::{Deserialize, Serialize};
use std::{
    ffi::{c_char, c_int, OsStr},
    fmt::Debug,
    mem::zeroed,
};

pub struct Lib(Library);

impl Lib {
    const DEFAULT_FILENAME: &'static str = "/usr/lib/liblyn_drv.so";

    pub fn new<P>(filename: P) -> Result<Self>
    where
        P: AsRef<OsStr>,
    {
        unsafe { Ok(Lib(Library::new(filename)?)) }
    }

    pub fn try_default() -> Result<Self> {
        Self::new(Self::DEFAULT_FILENAME)
    }
}

#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord, Serialize, Deserialize)]
pub enum ErrType {
    CHIP,
    BOARD,
    NODE,
    OTHER,
}

#[repr(C)]
pub struct RawErrMsg {
    typ: ErrType,
    devid: u32,
    enable_recover: c_int,
    msg: [c_char; 256],
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ErrMsg {
    pub typ: ErrType,
    pub devid: u32,
    pub enable_recover: i32,
    pub msg: String,
}

impl TryFrom<&RawErrMsg> for ErrMsg {
    type Error = Error;

    fn try_from(value: &RawErrMsg) -> std::result::Result<Self, Self::Error> {
        Ok(ErrMsg {
            typ: value.typ,
            devid: value.devid,
            enable_recover: value.enable_recover as i32,
            msg: string_from_c(value.msg.as_ref())?,
        })
    }
}

type RawErrMsgList = [RawErrMsg; 256];

type ErrCb = fn(*mut RawErrMsg, *mut c_int) -> c_int;

#[derive(Clone)]
pub struct ExceptionSymbols<'lib> {
    lynd_get_device_exception: Symbol<'lib, fn(*mut RawErrMsg, *mut c_int) -> c_int>,
    lynd_get_device_current_exception: Symbol<'lib, fn(*mut RawErrMsg, *mut c_int, ErrCb) -> c_int>,
}

impl<'lib> ExceptionSymbols<'lib> {
    pub fn new(lib: &'lib Lib) -> Result<Self> {
        unsafe {
            Ok(Self {
                lynd_get_device_exception: lib.0.get(b"lynd_get_device_exception")?,
                lynd_get_device_current_exception: lib
                    .0
                    .get(b"lynd_get_device_current_exception")?,
            })
        }
    }

    pub fn get_device_exception(&self) -> Result<Vec<ErrMsg>> {
        let mut cnt = 0;
        let mut list: RawErrMsgList = unsafe { zeroed() };
        Error::check((self.lynd_get_device_exception)(
            list.as_mut_ptr(),
            &mut cnt,
        ))?;
        let mut ret = Vec::new();
        for raw in &list[..cnt as usize] {
            let err_msg = ErrMsg::try_from(raw)?;
            ret.push(err_msg);
        }
        Ok(ret)
    }

    pub fn get_device_current_exception(&self, cb: ErrCb) -> Result<Vec<ErrMsg>> {
        let mut cnt = 0;
        let mut list: RawErrMsgList = unsafe { zeroed() };
        Error::check((self.lynd_get_device_current_exception)(
            &mut list[0],
            &mut cnt,
            cb,
        ))?;
        let mut ret = Vec::new();
        for raw in list {
            let err_msg = ErrMsg::try_from(&raw)?;
            ret.push(err_msg);
        }
        Ok(ret)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    // use core::slice;

    // fn cb(raw_err_msgs: *mut RawErrMsg, cnt: *mut c_int) -> c_int {
    //     let cnt = unsafe { *cnt };
    //     let raw_err_msg_slice = unsafe { slice::from_raw_parts(raw_err_msgs, cnt as usize) };
    //     for raw in raw_err_msg_slice {
    //         let err_msg = ErrMsg::try_from(raw).unwrap();
    //         println!("{:?}", err_msg);
    //     }
    //     0
    // }

    #[test]
    fn test_get_device_exception() {
        let lib = Lib::try_default().unwrap();
        let symbols = ExceptionSymbols::new(&lib).unwrap();
        let err_msgs = symbols.get_device_exception().unwrap();
        println!("{err_msgs:#?}");
        // symbols.get_device_current_exception(cb).unwrap();
    }
}
