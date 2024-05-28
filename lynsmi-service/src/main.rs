use axum::extract::State;
use axum::Json;
use axum::{routing::get, Router};
use lyndriver::drv::ErrMsg;
use lyndriver::{smi::DriverVersion, Result};
use lynsmi_service::drv_exception::{listen, DRV_EXCEPTION_MAP};
use lynsmi_service::models::*;
use serde_json::json;
use std::collections::HashMap;
use std::ops::Deref;
use std::sync::{Arc, Mutex};
use std::thread;
use std::time::Duration;
use tokio::task::spawn_blocking;
use tracing::{info, warn};

struct SMIData {
    device_count: i32,
    driver_version: DriverVersion,
    is_support_topology: bool,
    devices: Mutex<PropsMap>,
    device_topology_list: Mutex<P2PAttrList>,
}

async fn get_devices(State(smi_data): State<Arc<SMIData>>) -> Json<serde_json::Value> {
    Json(json!(smi_data
        .devices
        .lock()
        .expect("lock smi_data.devices failed")
        .deref()))
}

async fn get_device_count(State(smi_data): State<Arc<SMIData>>) -> Json<i32> {
    Json(smi_data.device_count)
}

async fn get_driver_version(State(smi_data): State<Arc<SMIData>>) -> Json<DriverVersion> {
    Json(smi_data.driver_version.clone())
}

async fn get_drv_exception_map() -> Json<HashMap<u32, ErrMsg>> {
    Json(
        DRV_EXCEPTION_MAP
            .lock()
            .expect("lock DRV_EXCEPTION_MAP failed")
            .clone(),
    )
}

async fn get_device_topology_list(State(smi_data): State<Arc<SMIData>>) -> Json<serde_json::Value> {
    if !smi_data.is_support_topology {
        Json(json!(Option::<&[Result<DeviceP2PAttr>]>::None))
    } else {
        Json(json!(Some(
            smi_data
                .device_topology_list
                .lock()
                .expect("lock smi_data.device_topology_list failed")
                .as_slice()
        )))
    }
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let smi_lib = lyndriver::smi::Lib::try_default()?;
    let smi_common = lyndriver::smi::CommonSymbols::new(&smi_lib)?;
    let smi_props = lyndriver::smi::PropsSymbols::new(&smi_lib)?;
    let smi_topology = lyndriver::smi::TopologySymbols::new(&smi_lib);

    let device_count = smi_common.get_device_cnt()?;
    info!("device_count {}", device_count);
    let driver_version = smi_props.get_driver_version().clone();
    info!("driver_version {:?}", driver_version);
    let is_support_topology = smi_topology.is_ok();
    info!("is_support_topology {:?}", is_support_topology);

    let smi_data = Arc::new(SMIData {
        device_count,
        driver_version,
        is_support_topology,
        devices: Mutex::new(HashMap::new()),
        device_topology_list: Mutex::new(Vec::new()),
    });
    let smi_data_clone = smi_data.clone();

    let smi_thread = spawn_blocking(move || {
        thread::scope(|s| {
            s.spawn(move || {
                info!("start listen drv_exception");
                if let Err(e) = listen() {
                    warn!("listen drv_exception failed {:?}", e);
                }
            });

            let smi_props =
                lyndriver::smi::PropsSymbols::new(&smi_lib).expect("init PropsSymbols failed");
            for device_id in 0..device_count {
                let smi = smi_props.clone();
                let smi_data = smi_data.clone();
                s.spawn(move || loop {
                    let no_exception = DRV_EXCEPTION_MAP
                        .lock()
                        .expect("lock DRV_EXCEPTION_MAP failed")
                        .get(&(device_id as u32))
                        .is_none();
                    if no_exception {
                        info!("start update props for device {}", device_id);
                        loop {
                            let props = smi.get_props(device_id);
                            let is_err = props.is_err();
                            smi_data
                                .devices
                                .lock()
                                .expect("lock smi_data.devices failed")
                                .insert(device_id, props);
                            if is_err {
                                info!(
                                    "get props for device {} failed, retry after chip recovery",
                                    device_id
                                );
                                break;
                            }
                        }
                    } else {
                        thread::sleep(Duration::from_secs(1));
                    }
                });
            }
            if is_support_topology {
                let smi_topology = lyndriver::smi::TopologySymbols::new(&smi_lib)
                    .expect("support topology but init TopologySymbols failed");
                s.spawn(move || loop {
                    info!("start get device_topology_list");
                    let mut device_topology_list = Vec::new();
                    let mut any_err = false;
                    for src_device in 0..device_count {
                        for dst_device in (src_device + 1)..device_count {
                            let attr = smi_topology
                                .get_device_p2p_attr(src_device, dst_device)
                                .map(|v| DeviceP2PAttr {
                                    device_pair: (src_device, dst_device),
                                    attr: v,
                                });
                            if attr.is_err() {
                                warn!("get_device_p2p_attr return err {:?}", attr);
                                any_err = true;
                            }
                            device_topology_list.push(attr);
                        }
                    }

                    *smi_data
                        .device_topology_list
                        .lock()
                        .expect("lock smi_data.device_topology_list failed") = device_topology_list;
                    if any_err {
                        thread::sleep(Duration::from_secs(60));
                    } else {
                        info!("get device topology list successed");
                        break;
                    }
                });
            }
        });
    });

    // build our application with a single route
    let app = Router::new()
        .route("/devices", get(get_devices))
        .route("/device_count", get(get_device_count))
        .route("/driver_version", get(get_driver_version))
        .route("/device_topology_list", get(get_device_topology_list))
        .route("/drv_exception_map", get(get_drv_exception_map))
        .with_state(smi_data_clone);

    // run our app with hyper, listening globally on port 5432
    let listener = tokio::net::TcpListener::bind("0.0.0.0:5432").await?;
    axum::serve(listener, app).await?;
    smi_thread.await?;

    Ok(())
}
