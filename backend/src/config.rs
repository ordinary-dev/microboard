//! Configuration
use std::env;

/// Service configuration: database credentials, port...
pub struct Config {
    /// Port (it will be used in a string, no need to convert it to integer)
    pub port: String,
    /// A long secret string used to authenticate users.
    pub secret: String,
    /// Maximum number of connections
    pub pg_max_connections: u32,
    /// Database host (e.g. localhost:5432)
    pg_host: String,
    /// Database name
    pg_db: String,
    /// Database user
    pg_user: String,
    /// Database password
    pg_password: String,
}

impl Config {
    /// Get a new instance of config
    pub fn new() -> Config {
        // Try to load .env file
        _ = dotenvy::dotenv();

        // Check secret string
        let secret = env::var("MB_SECRET").expect("MB_SECRET is undefined");
        if secret.len() < 16 {
            panic!("MB_SECRET is less than 16 characters");
        }

        Config {
            port: env::var("MB_PORT").unwrap_or("8080".to_string()),
            secret,
            pg_max_connections: env::var("MB_PG_MAX_CONNECTIONS")
                .unwrap_or("5".to_string())
                .parse()
                .expect("MB_PG_MAX_CONNECTIONS is not an integer"),
            pg_host: env::var("MB_PG_HOST").expect("MB_PG_HOST is undefined"),
            pg_db: env::var("MB_PG_DB").expect("MB_PG_DB is undefined"),
            pg_user: env::var("MB_PG_USER").expect("MB_PG_USER is undefined"),
            pg_password: env::var("MB_PG_PASSWORD").expect("MB_PG_PASSWORD is undefined"),
        }
    }

    /// Get a url to postgres db for sqlx
    pub fn get_db_url(&self) -> String {
        format!(
            "postgres://{}:{}@{}/{}",
            self.pg_user, self.pg_password, self.pg_host, self.pg_db
        )
    }
}
