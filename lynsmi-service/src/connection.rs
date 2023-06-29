use futures_util::{SinkExt, StreamExt};
use lynsmi::Props;
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

pub type AllProps = Vec<Option<Props>>;

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
    pub async fn send(&mut self, data: &AllProps) -> Result<()> {
        let b = serde_json::to_string(&data)?;
        self.framed.send(b).await?;
        Ok(())
    }

    #[instrument]
    pub async fn next(&mut self) -> Option<Result<AllProps>> {
        Some(
            self.framed
                .next()
                .await?
                .map_err(|e| Error::AnyDelimiterCodecError(e))
                .and_then(|v| serde_json::from_slice::<AllProps>(&v).map_err(|e| e.into())),
        )
    }
}
