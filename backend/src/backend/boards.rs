use serde::{Deserialize, Serialize};
use sqlx::{Pool, Postgres};

#[derive(sqlx::FromRow, Serialize, Deserialize, Debug, Clone)]
pub struct Board {
    pub code: String,
    pub name: String,
}

impl Board {
    /// Get all boards.
    pub async fn get_all(pool: &Pool<Postgres>) -> anyhow::Result<Vec<Board>> {
        let boards = sqlx::query_as::<_, Board>("SELECT code, name FROM boards")
            .fetch_all(pool)
            .await?;

        Ok(boards)
    }

    /// Get board by code.
    pub async fn get(pool: &Pool<Postgres>, code: &str) -> anyhow::Result<Board> {
        let board = sqlx::query_as::<_, Board>("SELECT code, name FROM boards WHERE code = $1")
            .bind(code)
            .fetch_one(pool)
            .await?;

        Ok(board)
    }

    /// Write a new record to the database.
    pub async fn insert(&self, pool: &Pool<Postgres>) -> anyhow::Result<()> {
        sqlx::query("INSERT INTO boards (code, name) VALUES ($1, $2)")
            .bind(&self.code)
            .bind(&self.name)
            .execute(pool)
            .await?;

        Ok(())
    }

    /// Update existing record.
    pub async fn update(&self, pool: &Pool<Postgres>) -> anyhow::Result<()> {
        sqlx::query("UPDATE boards SET name = $1 WHERE code = $2")
            .bind(&self.name)
            .bind(&self.code)
            .execute(pool)
            .await?;

        Ok(())
    }

    /// Delete record from the database.
    pub async fn delete(pool: &Pool<Postgres>, code: &str) -> anyhow::Result<()> {
        sqlx::query("DELETE FROM boards WHERE code = $1")
            .bind(code)
            .execute(pool)
            .await?;

        Ok(())
    }
}
