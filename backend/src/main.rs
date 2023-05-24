//! Microboard - Image board server.

use axum::{
    routing::{get, post},
    Router,
};
use sqlx::postgres::PgPoolOptions;

mod api;
mod backend;
mod config;
mod error;
mod state;

#[tokio::main]
async fn main() {
    let cfg = config::Config::new();

    let database_pool = PgPoolOptions::new()
        .max_connections(cfg.pg_max_connections)
        .connect(&cfg.get_db_url())
        .await
        .unwrap();

    backend::apply_migrations(&database_pool).await.unwrap();

    let app_state = state::AppState::new(database_pool);

    let app = Router::new()
        .route(
            "/api/v0/boards",
            get(api::boards::get_all).post(api::boards::create),
        )
        .route(
            "/api/v0/boards/:code",
            get(api::boards::get)
                .put(api::boards::update)
                .delete(api::boards::delete),
        )
        .route("/api/v0/threads", post(api::threads::create_thread))
        .route("/api/v0/threads/:id", get(api::threads::get_thread))
        .route("/api/v0/posts", post(api::posts::create_post))
        .route("/api/v0/posts/:id", get(api::posts::get_post))
        .route("/api/v0/files", post(api::files::create_file))
        .with_state(app_state);

    let addr = format!("0.0.0.0:{}", cfg.port).parse().unwrap();
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
