use crate::{backend::Board, state::AppState};
use axum::{
    extract::{Path, State},
    response::IntoResponse,
    http::StatusCode,
    Json,
};
use std::sync::Arc;

/// Get all boards.
pub async fn get_all(State(state): State<Arc<AppState>>) -> Result<Json<Vec<Board>>, impl IntoResponse> {
    match Board::get_all(&state.db).await {
        Ok(board_list) => Ok(Json(board_list)),

        // For now, let's just return 500.
        // TODO: add different error codes for different errors.
        Err(err) => Err((StatusCode::INTERNAL_SERVER_ERROR, err.to_string())),
    }
}

/// Create a new board.
pub async fn create(
    State(state): State<Arc<AppState>>,
    Json(board): Json<Board>,
) -> Result<Json<Board>, impl IntoResponse> {
    match board.insert(&state.db).await {
        Ok(_) => Ok(Json(board)),
        Err(err) => Err((StatusCode::BAD_REQUEST, err.to_string())),
    }
}

/// Get board.
pub async fn get(
    State(state): State<Arc<AppState>>,
    Path(code): Path<String>,
) -> Result<Json<Board>, impl IntoResponse> {
    match Board::get(&state.db, &code).await {
        Ok(board) => Ok(Json(board)),
        Err(err) => Err((StatusCode::BAD_REQUEST, err.to_string())),
    }
}

/// Update board info.
pub async fn update(
    State(state): State<Arc<AppState>>,
    Json(board): Json<Board>,
) -> Result<Json<Board>, impl IntoResponse> {
    match board.update(&state.db).await {
        Ok(_) => Ok(Json(board)),
        Err(err) => Err((StatusCode::BAD_REQUEST, err.to_string())),
    }
}

/// Delete board.
pub async fn delete(
    State(state): State<Arc<AppState>>,
    Path(code): Path<String>,
) -> Result<(), impl IntoResponse> {
    match Board::delete(&state.db, &code).await {
        Ok(_) => Ok(()),
        Err(err) => Err((StatusCode::BAD_REQUEST, err.to_string())),
    }
}
