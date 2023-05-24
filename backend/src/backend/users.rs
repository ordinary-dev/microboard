//! Users module.
use sqlx::{FromRow, Pool, Postgres};

#[derive(Debug, FromRow)]
pub struct User {
    pub username: String,
    bcrypt: String,
}

impl User {
    /// Create a new user.
    /// Does not store anything in the database.
    pub fn new(username: &str, plain_password: &str) -> anyhow::Result<User> {
        Ok(User {
            username: username.to_string(),
            bcrypt: bcrypt::hash(plain_password, 10)?,
        })
    }

    /// Get user by username.
    pub async fn get(pool: &Pool<Postgres>, username: &str) -> anyhow::Result<User> {
        let user =
            sqlx::query_as::<_, User>("SELECT username, bcrypt FROM users WHERE username = $1")
                .bind(username)
                .fetch_one(pool)
                .await?;

        Ok(user)
    }

    /// Check if provided password is correct.
    pub async fn check_password(
        pool: &Pool<Postgres>,
        username: &str,
        password: &str,
    ) -> anyhow::Result<bool> {
        let user = Self::get(pool, username).await?;
        Ok(bcrypt::verify(password, &user.bcrypt)?)
    }

    /// Save a user in the database.
    pub async fn insert(&self, pool: &Pool<Postgres>) -> anyhow::Result<()> {
        sqlx::query("INSERT INTO users (username, bcrypt) VALUES ($1, $2)")
            .bind(self.username.clone())
            .bind(self.bcrypt.clone())
            .execute(pool)
            .await?;

        Ok(())
    }

    /// Update user.
    pub async fn update(&self, pool: &Pool<Postgres>) -> anyhow::Result<()> {
        sqlx::query("UPDATE users SET bcrypt = $1 WHERE username = $2")
            .bind(self.bcrypt.clone())
            .bind(self.username.clone())
            .execute(pool)
            .await?;

        Ok(())
    }

    /// Delete user.
    pub async fn delete(pool: &Pool<Postgres>, username: &str) -> anyhow::Result<()> {
        sqlx::query("DELETE FROM users WHERE code = $1")
            .bind(username)
            .execute(pool)
            .await?;

        Ok(())
    }
}
