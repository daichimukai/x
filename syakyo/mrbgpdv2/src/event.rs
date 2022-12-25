use crate::packets::keepalive::KeepAliveMessage;
use crate::packets::open::OpenMessage;

// https://www.rfc-editor.org/rfc/rfc4271.html#section-8.1.2
#[derive(PartialEq, Eq, Debug, Clone, Hash)]
pub enum Event {
    ManualStart,
    TcpConnectionConfirmed,
    BgpOpen(OpenMessage),
    KeepAliveMsg(KeepAliveMessage),
    // original event; not in RFC
    Established,
}
