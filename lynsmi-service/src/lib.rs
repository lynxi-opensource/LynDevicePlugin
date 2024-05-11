pub mod models {
    use lyndriver::{
        smi::{P2PAttr, Props},
        Result,
    };
    use serde::Serialize;
    use std::collections::HashMap;

    #[derive(Debug, Clone, Serialize)]
    pub struct DeviceP2PAttr {
        pub device_pair: (i32, i32),
        pub attr: P2PAttr,
    }

    pub type PropsMap = HashMap<i32, Result<Props>>;
    pub type P2PAttrList = Vec<Result<DeviceP2PAttr>>;
}

pub mod drv_exception;
