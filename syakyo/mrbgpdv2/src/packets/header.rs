use bytes::{BufMut, BytesMut};

use crate::error::ConvertBytesToBgpMessageError;

#[derive(PartialEq, Eq, Debug, Clone, Hash)]
pub struct Header {
    length: u16,
    pub type_: MessageType,
}

impl Header {
    pub fn new(length: u16, type_: MessageType) -> Self {
        Self { length, type_ }
    }
}

impl TryFrom<BytesMut> for Header {
    type Error = ConvertBytesToBgpMessageError;

    fn try_from(bytes: BytesMut) -> Result<Self, Self::Error> {
        let _marker = &bytes[0..16];
        let length = u16::from_be_bytes([bytes[16], bytes[17]]);
        let type_ = bytes[18].try_into()?;

        Ok(Header { length, type_ })
    }
}

impl From<Header> for BytesMut {
    fn from(header: Header) -> Self {
        let mut buf = BytesMut::new();
        let marker = [255u8; 16];
        let length = header.length.to_be_bytes();
        let type_: u8 = header.type_.into();

        buf.put(&marker[..]);
        buf.put(&length[..]);
        buf.put_u8(type_);

        buf
    }
}

#[derive(PartialEq, Eq, Debug, Clone, Copy, Hash)]
pub enum MessageType {
    Open,
    KeepAlive,
}

impl TryFrom<u8> for MessageType {
    type Error = ConvertBytesToBgpMessageError;

    fn try_from(num: u8) -> Result<Self, Self::Error> {
        match num {
            1 => Ok(MessageType::Open),
            4 => Ok(MessageType::KeepAlive),
            _ => Err(Self::Error::from(anyhow::anyhow!(
                "failed to convert {} to BGP message type (expected 1-4)",
                num
            ))),
        }
    }
}

impl From<MessageType> for u8 {
    fn from(type_: MessageType) -> Self {
        match type_ {
            MessageType::Open => 1,
            MessageType::KeepAlive => 4,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn convert_bytes_to_header_and_header_to_bytes() {
        let header = Header::new(29, MessageType::Open);
        let header_bytes: BytesMut = header.clone().into();
        let header2: Header = header_bytes.try_into().unwrap();

        assert_eq!(header, header2);
    }
}
