use rayon::ThreadPoolBuildError;
use std::{
    ffi::FromBytesWithNulError, io::Error as IoError, num::ParseIntError, str::Utf8Error,
    string::FromUtf8Error,
};

#[non_exhaustive]
#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("lynxi error code {0}")]
    Lyn(i32),
    #[error("{0}")]
    FromBytesWithNul(#[from] FromBytesWithNulError),
    #[error("{0}")]
    FromUtf8(#[from] FromUtf8Error),
    #[error("{0}")]
    Utf8(#[from] Utf8Error),
    #[error("{0}")]
    LibLoading(#[from] libloading::Error),
    #[error("{0}")]
    Io(#[from] IoError),
    #[error("{0}")]
    ParseInt(#[from] ParseIntError),
    #[error("StripVersionPrefix")]
    StripVersionPrefix,
    #[error("VersionSplit")]
    SplitVersion,
    #[error("ThreadPoolBuildError")]
    ThreadPoolBuildError(#[from] ThreadPoolBuildError),
}

impl Error {
    pub(crate) fn check(code: i32) -> Result<(), Self> {
        if code == 0 {
            Ok(())
        } else {
            Err(Self::Lyn(code))
        }
    }
}
