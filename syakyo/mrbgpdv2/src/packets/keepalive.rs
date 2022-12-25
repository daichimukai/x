use bytes::BytesMut;

use super::header::{Header, MessageType};
use crate::error::ConvertBytesToBgpMessageError;

#[derive(PartialEq, Eq, Debug, Clone, Hash)]
pub struct KeepAliveMessage {
    header: Header,
}

impl KeepAliveMessage {
    pub fn new() -> Self {
        let header = Header::new(19, MessageType::KeepAlive);
        Self { header }
    }
}

impl Default for KeepAliveMessage {
    fn default() -> Self {
        Self::new()
    }
}

impl TryFrom<BytesMut> for KeepAliveMessage {
    type Error = ConvertBytesToBgpMessageError;

    fn try_from(bytes: BytesMut) -> Result<Self, Self::Error> {
        let header = Header::try_from(bytes)?;
        if header.type_ != MessageType::KeepAlive {
            return Err(anyhow::anyhow!("type is not keep-alive").into());
        }
        Ok(Self { header })
    }
}

impl From<KeepAliveMessage> for BytesMut {
    fn from(keep_alive: KeepAliveMessage) -> Self {
        keep_alive.header.into()
    }
}
