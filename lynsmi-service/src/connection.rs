use futures_util::{SinkExt, StreamExt};
use lynsmi::Props;
use serde::{Deserialize, Serialize};
use tokio::net::TcpStream;
use tokio_util::codec::{AnyDelimiterCodec, AnyDelimiterCodecError, Decoder, Framed};
use tracing::instrument;

#[derive(Debug, thiserror::Error)]
pub enum Error {
    #[error("{0}")]
    AnyDelimiterCodecError(#[from] AnyDelimiterCodecError),
    #[error("{0}")]
    SeardeJsonError(#[from] serde_json::Error),
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PropsWithID {
    pub id: usize,
    pub props: Option<Props>,
    pub err: Option<String>,
}

pub type Result<T> = std::result::Result<T, Error>;

#[derive(Debug)]
pub struct Connection {
    framed: Framed<TcpStream, AnyDelimiterCodec>,
}

impl Connection {
    pub fn new(socket: TcpStream) -> Self {
        let framed = AnyDelimiterCodec::new(vec![0], vec![0]).framed(socket);
        Self { framed }
    }

    #[instrument]
    pub async fn send(&mut self, data: &PropsWithID) -> Result<()> {
        let b = serde_json::to_string(&data)?;
        self.framed.send(b).await?;
        Ok(())
    }

    #[instrument]
    pub async fn next(&mut self) -> Option<Result<PropsWithID>> {
        Some(
            self.framed
                .next()
                .await?
                .map_err(|e| Error::AnyDelimiterCodecError(e))
                .and_then(|v| serde_json::from_slice::<PropsWithID>(&v).map_err(|e| e.into())),
        )
    }
}
