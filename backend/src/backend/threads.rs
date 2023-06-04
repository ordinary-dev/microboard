use serde::{Deserialize, Serialize};
use sqlx::{Pool, Postgres, Transaction};
use chrono::Utc;

#[derive(sqlx::FromRow, Serialize, Deserialize, Debug, Clone)]
pub struct Thread {
    pub id: i64,
    pub board_code: String,
    // Hidden field: updated_at: DateTime<Utc>,
}

impl Thread {
    /// Get all threads from the specified board.
    pub async fn get_all(pool: &Pool<Postgres>, board_code: &str, limit: i16, offset: i16) -> anyhow::Result<Vec<Thread>> {
        let boards = sqlx::query_as::<_, Thread>(
            "
            SELECT id, board_code
            FROM threads
            WHERE board_code = $1
            ORDER BY updated_at DESC
            LIMIT $2
            OFFSET $3
            ",
        )
        .bind(board_code)
        .bind(limit)
        .bind(offset)
        .fetch_all(pool)
        .await?;

        Ok(boards)
    }

    /// Get thread info by id.
    pub async fn get(pool: &Pool<Postgres>, id: i64) -> anyhow::Result<Thread> {
        let board = sqlx::query_as::<_, Thread>("SELECT id, board_code FROM threads WHERE id = $1")
            .bind(id)
            .fetch_one(pool)
            .await?;

        Ok(board)
    }

    /// Write a new record to the database.
    ///
    /// A transaction is required to create the first post after the thread is created.
    pub async fn insert(pool: &mut Transaction<'_, Postgres>, board_code: &str) -> anyhow::Result<i64> {
        let id: (i64,) = sqlx::query_as("INSERT INTO threads (board_code, updated_at) VALUES ($1, $2) RETURNING id")
            .bind(board_code)
            .bind(Utc::now())
            .fetch_one(pool)
            .await?;

        Ok(id.0)
    }

    /// Delete thread from the database.
    pub async fn delete(pool: &Pool<Postgres>, thread_id: i64) -> anyhow::Result<()> {
        sqlx::query("DELETE FROM threads WHERE id = $1")
            .bind(thread_id)
            .execute(pool)
            .await?;

        Ok(())
    }
}
