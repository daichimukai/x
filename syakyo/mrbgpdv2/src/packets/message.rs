use std::net::Ipv4Addr;

use bytes::BytesMut;

use crate::bgp_type::AutonomousSystemNumber;
use crate::error::{ConvertBytesToBgpMessageError};
use crate::packets::open::OpenMessage;
use crate::packets::header::{Header, MessageType};
use crate::packets::keepalive::KeepAliveMessage;

#[derive(PartialEq, Eq, Debug, Clone, Hash)]
pub enum Message {
    Open(OpenMessage),
    KeepAlive(KeepAliveMessage),
}

impl Message {
    pub fn new_open(my_as_number: AutonomousSystemNumber, my_ip_addr: Ipv4Addr) -> Self {
        Self::Open(OpenMessage::new(my_as_number, my_ip_addr))
    }

    pub fn new_keepalive() -> Self {
        Self::KeepAlive(KeepAliveMessage::new())
    }
}

impl TryFrom<BytesMut> for Message {
    type Error = ConvertBytesToBgpMessageError;

    fn try_from(bytes: BytesMut) -> Result<Self, Self::Error> {
        let header_bytes_length = 19;
        if bytes.len() < header_bytes_length {
            return Err(Self::Error::from(anyhow::anyhow!(
                "failed to convert bytes to BGP Message: \
                the length of bytes is shorter than the expected header length"
            )));
        }

        let header = Header::try_from(BytesMut::from(&bytes[0..header_bytes_length]))?;
        match header.type_ {
            MessageType::Open => Ok(Message::Open(OpenMessage::try_from(bytes)?)),
            MessageType::KeepAlive => Ok(Message::KeepAlive(KeepAliveMessage::try_from(bytes)?)),
        }
    }
}

impl From<Message> for BytesMut {
    fn from(message: Message) -> BytesMut {
        match message {
            Message::Open(open) => open.into(),
            Message::KeepAlive(keepalive) => keepalive.into(),
        }
    }
}
