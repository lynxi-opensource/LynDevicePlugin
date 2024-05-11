use crate::errors::*;
use std::ffi::{c_char, CStr};

pub(crate) fn string_from_c(data: &[c_char]) -> Result<String> {
    unsafe { Ok(CStr::from_ptr(data.as_ptr()).to_str()?.to_owned()) }
}
