use std::{
    collections::{BTreeMap, HashMap},
    env,
};

use k8s_openapi::{api::core::v1::Node, serde::Deserialize, serde_json::json};
use kube::{
    api::{Patch, PatchParams},
    Api, Client, ResourceExt,
};
use lynsmi::{DriverVersion, Props};
use tokio::{
    task::spawn_blocking,
    time::{interval, Duration},
};

use tracing::{error, info, warn};

#[derive(Debug, Deserialize)]
enum DeResult<T> {
    Ok(T),
    Err(String),
}

type PropsMap = HashMap<i32, DeResult<Props>>;

async fn get_devices() -> reqwest::Result<PropsMap> {
    let resp = reqwest::get("http://localhost:5432/devices")
        .await?
        .json::<HashMap<i32, DeResult<Props>>>()
        .await?;
    Ok(resp)
}

async fn get_driver_version_label() -> anyhow::Result<(String, String)> {
    const DRIVER_VERSION_LABEL_KEY: &str = "lynxi.com/driver-version";

    let driver_version = spawn_blocking(|| DriverVersion::local()).await??;
    let driver_version_str = format!(
        "{}.{}.{}",
        driver_version.0, driver_version.1, driver_version.2
    );

    Ok((DRIVER_VERSION_LABEL_KEY.to_string(), driver_version_str))
}

async fn patch_labels(nodes: &Api<Node>, labels: BTreeMap<String, String>) -> anyhow::Result<()> {
    let node_name = env::var("NODE_NAME")?;

    info!("update labels {:?}", labels);
    let patch = json!({"metadata": {
        "labels": labels
    }});
    let pp = PatchParams::apply("lynxi-device-discovery");
    nodes
        .patch(&node_name, &pp, &Patch::Strategic(patch))
        .await?;

    Ok(())
}

fn get_target_labels(devices: &PropsMap) -> BTreeMap<String, String> {
    let mut target_labels = BTreeMap::new();
    for (_, device) in devices {
        match device {
            DeResult::Ok(v) => {
                target_labels.insert(
                    format!("lynxi.com/{}.present", v.device.name),
                    "true".to_string(),
                );
            }
            DeResult::Err(e) => {
                warn!("get device props return err: {}", e);
            }
        };
    }
    target_labels
}

async fn patch_lynxi_labels() -> anyhow::Result<()> {
    let client = Client::try_default().await?;
    let api_node: Api<Node> = Api::all(client);

    {
        let driver_version_label = get_driver_version_label().await?;
        let mut driver_version_label_map = BTreeMap::new();
        driver_version_label_map.insert(
            driver_version_label.0.clone(),
            driver_version_label.1.clone(),
        );
        patch_labels(&api_node, driver_version_label_map).await?;
    }

    let filter = |k: &str| -> bool {
        k.starts_with("lynxi.com/") && k.ends_with(".present") && k != "lynxi.com/apu.present"
    };

    let mut ticker = interval(Duration::from_secs(60));

    loop {
        ticker.tick().await;
        let devices_result = get_devices().await;
        match devices_result {
            Err(e) => {
                error!("request devices from lynsmi-service failed: {}", e);
                ticker.reset_immediately();
            }
            Ok(devices) => {
                let mut target_labels = get_target_labels(&devices);

                let node_name = env::var("NODE_NAME")?;
                let node = api_node.get(&node_name).await?;
                let mut node_labels = node.labels().clone();
                for (_, v) in node_labels.iter_mut().filter(|(k, _)| filter(k)) {
                    *v = "false".to_string();
                }

                node_labels.append(&mut target_labels);

                if node.labels() != &node_labels {
                    let node_labels: BTreeMap<String, String> =
                        node_labels.into_iter().filter(|(k, _)| filter(k)).collect();
                    patch_labels(&api_node, node_labels).await?;
                }
            }
        };
    }
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    patch_lynxi_labels().await
}
