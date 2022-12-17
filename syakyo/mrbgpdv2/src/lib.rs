#![feature(exclusive_range_pattern, arc_unwrap_or_clone)]
#![allow(dead_code)]

pub mod config;
pub mod peer;

mod bgp_type;
mod connection;
mod error;
mod event;
mod event_queue;
mod state;
