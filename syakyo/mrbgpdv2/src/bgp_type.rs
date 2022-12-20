use crate::error::ConvertBytesToBgpMessageError;

#[derive(PartialEq, Eq, Debug, Clone, Copy, Hash, PartialOrd, Ord)]
pub struct AutonomousSystemNumber(u16);

impl From<AutonomousSystemNumber> for u16 {
    fn from(as_number: AutonomousSystemNumber) -> u16 {
        as_number.0
    }
}

impl From<u16> for AutonomousSystemNumber {
    fn from(as_number: u16) -> AutonomousSystemNumber {
        Self(as_number)
    }
}

#[derive(PartialEq, Eq, Debug, Clone, Copy, Hash, PartialOrd, Ord)]
pub struct HoldTime(u16);

impl HoldTime {
    pub fn new() -> Self {
        Default::default()
    }
}

impl From<HoldTime> for u16 {
    fn from(t: HoldTime) -> u16 {
        t.0
    }
}

impl From<u16> for HoldTime {
    fn from(t: u16) -> HoldTime {
        HoldTime(t)
    }
}

impl Default for HoldTime {
    fn default() -> Self {
        HoldTime(0)
    }
}

#[derive(PartialEq, Eq, Debug, Clone, Hash, PartialOrd, Ord)]
pub struct Version(u8);

impl Version {
    pub fn new() -> Self {
        Default::default()
    }
}

impl From<Version> for u8 {
    fn from(v: Version) -> u8 {
        v.0
    }
}

impl TryFrom<u8> for Version {
    type Error = ConvertBytesToBgpMessageError;

    fn try_from(v: u8) -> Result<Self, Self::Error> {
        if v <= 4 {
            Ok(Version(v))
        } else {
            Err(Self::Error::from(anyhow::anyhow!(
                "failed to convert to BGP version: expected <= 4, got: {}",
                v
            )))
        }
    }
}

impl Default for Version {
    fn default() -> Self {
        Version(4)
    }
}
