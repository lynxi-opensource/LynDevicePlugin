use tokio::{io::AsyncReadExt, net::TcpStream};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let mut socket = TcpStream::connect("127.0.0.1:8080").await?;
    let dst = &mut String::new();
    while let Ok(_) = socket.read_to_string(dst).await {
        println!("{dst}");
    }

    Ok(())
}
