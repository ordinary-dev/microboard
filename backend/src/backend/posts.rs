use chrono::prelude::*;
use serde::{Deserialize, Serialize};
use sqlx::{Pool, Postgres};

#[derive(sqlx::FromRow, Serialize, Deserialize, Debug, Clone)]
pub struct Post {
    pub id: i64,
    pub body: String,
    pub thread_id: i64,
    pub created_at: DateTime<Utc>,
}

impl Post {
    /// Get all posts from thread.
    pub async fn get_all(pool: &Pool<Postgres>, thread_id: i64) -> anyhow::Result<Vec<Post>> {
        let posts = sqlx::query_as::<_, Post>(
            "
            SELECT *
            FROM posts
            WHERE thread_id = $1
            ORDER BY id
            ",
        )
        .bind(thread_id)
        .fetch_all(pool)
        .await?;

        Ok(posts)
    }

    /// Get post info by id.
    pub async fn get(pool: &Pool<Postgres>, post_id: i64) -> anyhow::Result<Post> {
        let post = sqlx::query_as::<_, Post>("SELECT * FROM posts WHERE id = $1")
            .bind(post_id)
            .fetch_one(pool)
            .await?;

        Ok(post)
    }

    /// Write a new record to the database.
    pub async fn insert(&self, pool: &Pool<Postgres>) -> anyhow::Result<()> {
        sqlx::query("INSERT INTO posts (body, thread_id) VALUES ($1, $2)")
            .bind(&self.body)
            .bind(self.thread_id)
            .execute(pool)
            .await?;

        Ok(())
    }

    /// Delete post from the database.
    pub async fn delete(pool: &Pool<Postgres>, post_id: i64) -> anyhow::Result<()> {
        sqlx::query("DELETE FROM posts WHERE id = $1")
            .bind(post_id)
            .execute(pool)
            .await?;

        Ok(())
    }
}
