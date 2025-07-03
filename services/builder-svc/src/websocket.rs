use axum::extract::ws::{WebSocket, Message};
use futures_util::{SinkExt, StreamExt};
use tracing::{error, info, instrument};

use crate::AppState;

#[instrument(skip(socket, state))]
pub async fn websocket_handler(socket: WebSocket, state: AppState) {
    info!("WebSocket connection established");
    
    let (mut sender, mut receiver) = socket.split();
    
    // Send welcome message
    if let Err(e) = sender.send(Message::Text("Welcome to Page Magic!".to_string())).await {
        error!("Failed to send welcome message: {}", e);
        return;
    }
    
    // Handle incoming messages
    while let Some(msg) = receiver.next().await {
        match msg {
            Ok(Message::Text(text)) => {
                info!("Received text message: {}", text);
                
                // Echo the message back
                if let Err(e) = sender.send(Message::Text(format!("Echo: {}", text))).await {
                    error!("Failed to send message: {}", e);
                    break;
                }
            }
            Ok(Message::Binary(_)) => {
                info!("Received binary message");
            }
            Ok(Message::Ping(ping)) => {
                info!("Received ping");
                if let Err(e) = sender.send(Message::Pong(ping)).await {
                    error!("Failed to send pong: {}", e);
                    break;
                }
            }
            Ok(Message::Pong(_)) => {
                info!("Received pong");
            }
            Ok(Message::Close(_)) => {
                info!("WebSocket connection closed");
                break;
            }
            Err(e) => {
                error!("WebSocket error: {}", e);
                break;
            }
        }
    }
    
    info!("WebSocket connection ended");
}
