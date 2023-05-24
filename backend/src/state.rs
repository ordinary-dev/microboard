//! Service status.
use sqlx::{Pool, Postgres};
use std::sync::Arc;

/// Shared app state.
///
/// Any request handler can access this structure.
pub struct AppState {
    /// Database pool.
    pub db: Pool<Postgres>,
}

impl AppState {
    pub fn new(db_pool: Pool<Postgres>) -> Arc<AppState> {
        let state = AppState { db: db_pool };
        Arc::new(state)
    }
}
