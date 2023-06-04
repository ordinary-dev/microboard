//! Database migrations.

use futures::TryStreamExt;
use sqlx::{Pool, Postgres, Row};

/// List of queries that will be executed once.
///
/// Do not change existing queries, create new ones.
const MIGRATIONS: [&str; 7] = [
    // Create 'boards' table.
    "CREATE TABLE boards (
        code VARCHAR(5) PRIMARY KEY,
        name VARCHAR(20) UNIQUE NOT NULL
    )",
    // Create 'threads' table.
    // - updated_at: timestamp of the last post, used to speed up sorting.
    "CREATE TABLE threads (
        id BIGSERIAL PRIMARY KEY,
        board_code VARCHAR(5) REFERENCES boards(code) ON DELETE CASCADE,
        updated_at TIMESTAMPTZ NOT NULL
    )",
    // Create 'posts' table.
    "CREATE TABLE posts (
        id BIGSERIAL PRIMARY KEY,
        body VARCHAR(1024) NOT NULL,
        thread_id BIGINT REFERENCES threads(id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL
    )",
    // Create 'file_contents' table.
    // - hash: sha256 digest.
    // - size: file size in bytes.
    // - width: the width of the photo or video in pixels.
    // - height: the height of the photo or video in pixels.
    // - duration: the duration of the video or audio in seconds.
    "CREATE TABLE file_contents (
        hash VARCHAR(64) PRIMARY KEY,
        mime_type VARCHAR(10) NOT NULL,
        size BIGINT NOT NULL,
        width SMALLINT,
        height SMALLINT,
        duration SMALLINT
    )",
    // Create 'previews' table.
    "CREATE TABLE previews (
        hash VARCHAR(64) PRIMARY KEY,
        width SMALLINT NOT NULL,
        height SMALLINT NOT NULL
    )",
    // Create 'files' table.
    // - content_hash: sha256 digest of file content.
    // - preview_hash: sha256 digest of generated preview.
    "CREATE TABLE files (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(256) NOT NULL,
        post_id BIGINT REFERENCES posts(id) ON DELETE CASCADE,
        content_hash VARCHAR(64) REFERENCES file_contents(hash) ON DELETE CASCADE,
        preview_hash VARCHAR(64) REFERENCES previews(hash) ON DELETE SET NULL,
        nsfw BOOLEAN NOT NULL
    )",
    // Create 'users' table.
    // Users are administrators or moderators. Those who need additional rights.
    // We use JWT for authorization.
    "CREATE TABLE users (
        username VARCHAR(32) PRIMARY KEY,
        bcrypt VARCHAR(60) NOT NULL
    )",
];

/// Apply all pending migrations.
pub async fn apply_migrations(pool: &Pool<Postgres>) -> anyhow::Result<()> {
    create_migrations_table(pool).await?;

    let last_id = get_last_migration_id(pool).await?;

    // The transaction will either apply all migrations or apply nothing.
    // I think being in the middle is not very good.
    let mut transaction = pool.begin().await?;

    for (index, query) in MIGRATIONS.into_iter().enumerate() {
        let migration_id: i16 = index.try_into().unwrap();
        if migration_id > last_id {
            println!("Applying migration #{}", migration_id);
            // Run the query.
            sqlx::query(query).execute(&mut transaction).await?;
            // Save information about the last migration.
            sqlx::query("INSERT INTO migrations (id) VALUES ($1)")
                .bind(migration_id)
                .execute(&mut transaction)
                .await?;
        }
    }

    transaction.commit().await?;

    Ok(())
}

pub async fn create_migrations_table(pool: &Pool<Postgres>) -> anyhow::Result<()> {
    sqlx::query(
        "
        CREATE TABLE IF NOT EXISTS migrations (
            id SMALLINT PRIMARY KEY
        )",
    )
    .execute(pool)
    .await?;

    Ok(())
}

pub async fn get_last_migration_id(pool: &Pool<Postgres>) -> anyhow::Result<i16> {
    let mut rows = sqlx::query(
        "
        SELECT id FROM migrations
        ORDER BY id DESC
        LIMIT 1
    ",
    )
    .fetch(pool);

    if let Some(row) = rows.try_next().await? {
        let id: i16 = row.try_get("id")?;
        return Ok(id);
    }

    // No migrations were found
    Ok(-1)
}
