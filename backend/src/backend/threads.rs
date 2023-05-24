use serde::{Deserialize, Serialize};
use sqlx::{Pool, Postgres};

#[derive(sqlx::FromRow, Serialize, Deserialize, Debug, Clone)]
pub struct Thread {
    pub id: i64,
    pub board_code: String,
    // Hidden field: updated_at: DateTime<Utc>,
}

impl Thread {
    /// Get all threads from the specified board.
    pub async fn get_all(pool: &Pool<Postgres>, board_code: &str) -> anyhow::Result<Vec<Thread>> {
        let boards = sqlx::query_as::<_, Thread>(
            "
            SELECT id, board_code
            FROM threads
            WHERE board_code = $1
            ORDER BY updated_at DESC
            ",
        )
        .bind(board_code)
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
    pub async fn insert(&self, pool: &Pool<Postgres>) -> anyhow::Result<()> {
        sqlx::query("INSERT INTO threads (board_code) VALUES ($1)")
            .bind(&self.board_code)
            .execute(pool)
            .await?;

        Ok(())
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
