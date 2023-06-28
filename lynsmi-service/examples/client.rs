use anyhow::Ok;
use lynsmi_service::connection::Connection;
use tokio::net::TcpStream;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let socket = TcpStream::connect("127.0.0.1:5432").await?;
    let mut conn = Connection::new(socket);
    while let Some(v) = conn.next().await {
        println!("{v:?}");
    }
    Ok(())
}
