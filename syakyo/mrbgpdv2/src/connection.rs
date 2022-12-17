use anyhow::{Context, Result};
use tokio::net::{TcpListener, TcpStream};

use crate::config::{Config, Mode};
use crate::error::CreateConnectionError;

#[derive(Debug)]
pub struct Connection {
    conn: TcpStream,
}

impl Connection {
    pub async fn connect(config: &Config) -> Result<Self, CreateConnectionError> {
        let conn = match config.mode {
            Mode::Active => Self::connect_to_remote_peer(config).await,
            Mode::Passive => Self::wait_connection_from_remote_peer(config).await,
        }?;

        Ok(Self { conn })
    }

    async fn connect_to_remote_peer(config: &Config) -> Result<TcpStream> {
        let bgp_port = 179;

        TcpStream::connect((config.remote_ip, bgp_port))
            .await
            .context(format!(
                "failed to connect to remote peer {0}:{1}",
                config.remote_ip, bgp_port,
            ))
    }

    async fn wait_connection_from_remote_peer(config: &Config) -> Result<TcpStream> {
        let bgp_port = 179;
        let listener = TcpListener::bind((config.local_ip, bgp_port))
            .await
            .context(format!(
                "failed to bind to {0}:{1}",
                config.local_ip, bgp_port,
            ))?;

        Ok(listener
           .accept()
           .await
           .context(format!(
               "failed to accept TCP connection from remote peer: {0}:{1}",
               config.local_ip, bgp_port,
            ))?
           .0)
    }
}
