use tracing::{info, instrument};

use crate::config::Config;
use crate::connection::Connection;
use crate::event::Event;
use crate::event_queue::EventQueue;
use crate::state::State;
use crate::packets::message::Message;

// event driven state machine
#[derive(Debug)]
pub struct Peer {
    state: State,
    event_queue: EventQueue,
    tcp_connection: Option<Connection>,
    config: Config,
}

impl Peer {
    pub fn new(config: Config) -> Self {
        let state = State::Idle;
        let event_queue = EventQueue::new();
        Self {
            state,
            event_queue,
            tcp_connection: None,
            config,
        }
    }

    #[instrument]
    pub fn start(&mut self) {
        info!("peer is started.");
        self.event_queue.enqueue(Event::ManualStart);
    }

    #[instrument]
    pub async fn next(&mut self) {
        if let Some(event) = self.event_queue.dequeue() {
            info!("event is occurred, event={:?}.", event);
            self.handle_event(event).await;
        }

        if let Some(conn) = &mut self.tcp_connection {
            if let Some(message) = conn.get_message().await {
                info!("message is received, message={:?}.", message);
                self.handle_message(message);
            }
        }
    }

    async fn handle_event(&mut self, event: Event) {
        match &self.state {
            State::Idle => match event {
                Event::ManualStart => {
                    self.tcp_connection = Connection::connect(&self.config).await.ok();
                    if self.tcp_connection.is_some() {
                        self.event_queue.enqueue(Event::TcpConnectionConfirmed);
                    } else {
                        panic!("failed to establish TCP Connection: {:?}", self.config)
                    }
                    self.state = State::Connect;
                }
                _ => {},
            },
            State::Connect => match event {
                Event::TcpConnectionConfirmed => {
                    self.tcp_connection.as_mut().expect("tcp connection is not established")
                        .send(Message::new_open(
                                self.config.local_as,
                                self.config.local_ip,
                                ))
                        .await;
                    self.state = State::OpenSent;
                },
                _ => {},
            }
            State::OpenSent => match event {
                Event::BgpOpen(open) => {
                    self.tcp_connection.as_mut().expect("tcp connection is not established")
                        .send(Message::new_keepalive()).await;
                    self.state = State::OpenConfirm;
                },
                _ => {},
            },
            State::OpenConfirm => {
                todo!();
            },
        }
    }

    fn handle_message(&mut self, message: Message) {
        match message {
            Message::Open(open) => {
                self.event_queue.enqueue(Event::BgpOpen(open));
            },
            Message::KeepAlive(keepalive) => {
                self.event_queue.enqueue(Event::KeepAliveMsg(keepalive));
            },
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tokio::time::{sleep, Duration};

    #[tokio::test]
    async fn peer_can_transition_to_connect_state() {
        let config: Config = "64512 127.0.0.1 64513 127.0.0.2 active".parse().unwrap();
        let mut peer = Peer::new(config);
        peer.start();

        tokio::spawn(async move {
            let remote_config = "64513 127.0.0.2 64512 127.0.0.1 passive".parse().unwrap();
            let mut remote_peer = Peer::new(remote_config);
            remote_peer.start();
            remote_peer.next().await;
        });

        // wait to ensure that the `remote_peer` does its job before the `peer`.
        tokio::time::sleep(Duration::from_secs(1)).await;
        peer.next().await;
        assert_eq!(peer.state, State::Connect);
    }

    #[tokio::test]
    async fn peer_can_transition_to_opensent_state() {
        let config: Config = "64512 127.0.0.1 64513 127.0.0.2 active".parse().unwrap();
        let mut peer = Peer::new(config);
        peer.start();

        tokio::spawn(async move {
            let remote_config = "64513 127.0.0.2 64512 127.0.0.1 passive".parse().unwrap();
            let mut remote_peer = Peer::new(remote_config);
            remote_peer.start();
            remote_peer.next().await;
            remote_peer.next().await;
        });

        // wait to ensure that the `remote_peer` does its job before the `peer`.
        tokio::time::sleep(Duration::from_secs(1)).await;
        peer.next().await;
        peer.next().await;
        assert_eq!(peer.state, State::OpenSent);
    }

    #[tokio::test]
    async fn peer_can_transition_to_open_confirm_state() {
        let config: Config = "64512 127.0.0.1 64513 127.0.0.2 active".parse().unwrap();
        let mut peer = Peer::new(config);
        peer.start();

        tokio::spawn(async move {
            let remote_config = "64513 127.0.0.2 64512 127.0.0.1 passive".parse().unwrap();
            let mut remote_peer = Peer::new(remote_config);
            remote_peer.start();
            let max_step = 50;
            for _ in 0..max_step {
                remote_peer.next().await;
                if remote_peer.state == State::OpenConfirm {
                    break;
                }
                tokio::time::sleep(Duration::from_secs_f32(0.1)).await;
            }
        });

        // wait to ensure that the `remote_peer` does its job before the `peer`.
        tokio::time::sleep(Duration::from_secs(1)).await;
        let max_step = 50;
        for _ in 0..max_step {
            peer.next().await;
            if peer.state == State::OpenConfirm {
                break;
            }
            tokio::time::sleep(Duration::from_secs_f32(0.1)).await;
        }
        assert_eq!(peer.state, State::OpenConfirm);
    }
}
