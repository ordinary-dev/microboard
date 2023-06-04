use axum::{
    extract::{Query, Path, State},
    http::StatusCode,
    Json,
};
use serde::Deserialize;
use std::{
    sync::Arc,
    collections::HashMap,
};
use crate::{
    error::HttpError,
    backend::{Thread, Post},
    state::AppState,
};

pub async fn get_thread_by_id() {}

/// Get a list of threads from the board.
///
/// Parameters:
/// - limit (default: 10)
/// - offset (default: 0)
pub async fn get_threads_by_board_code(
    State(state): State<Arc<AppState>>,
    Path(board_code): Path<String>,
    Query(params): Query<HashMap<String, String>>
) -> Result<Json<Vec<Thread>>, HttpError> {
    let limit: i16 = match params.get("limit").unwrap_or(&"10".to_string()).parse() {
        Ok(res) => res,
        Err(err) => return Err(HttpError::new(StatusCode::BAD_REQUEST, &err.to_string())),
    };
    let offset: i16 = match params.get("offset").unwrap_or(&"0".to_string()).parse() {
        Ok(res) => res,
        Err(err) => return Err(HttpError::new(StatusCode::BAD_REQUEST, &err.to_string())),
    };

    match Thread::get_all(&state.db, &board_code, limit, offset).await {
        Ok(threads) => Ok(Json(threads)),
        Err(err) => Err(HttpError::new(StatusCode::BAD_REQUEST, &err.to_string())),
    }
}

#[derive(Deserialize)]
pub struct NewThread {
    pub board_code: String,
    pub body: String,
}

pub async fn create_thread(
    State(state): State<Arc<AppState>>,
    Json(new_thread): Json<NewThread>,
) -> Result<Json<Thread>, HttpError> {
    let mut transaction = match state.db.begin().await {
        Ok(transaction) => transaction,
        Err(err) => return Err(HttpError::new(StatusCode::INTERNAL_SERVER_ERROR, &err.to_string())),
    };

    // Create a new thread.
    let thread_id = match Thread::insert(&mut transaction, &new_thread.board_code).await {
        Ok(id) => id,
        Err(err) => return Err(HttpError::new(StatusCode::INTERNAL_SERVER_ERROR, &err.to_string())),
    };

    // Add the first post.
    match Post::insert(&mut transaction, &new_thread.body, thread_id).await {
        Ok(_) => (),
        Err(err) => return Err(HttpError::new(StatusCode::INTERNAL_SERVER_ERROR, &err.to_string())),
    };

    match transaction.commit().await {
        Ok(_) => (),
        Err(err) => return Err(HttpError::new(StatusCode::INTERNAL_SERVER_ERROR, &err.to_string())),
    };

    match Thread::get(&state.db, thread_id).await {
        Ok(thread) => Ok(Json(thread)),
        Err(err) => Err(HttpError::new(StatusCode::INTERNAL_SERVER_ERROR, &err.to_string())),
    }
}
