use chrono::{Duration, Utc};
use jsonwebtoken::{decode, encode, Algorithm, DecodingKey, EncodingKey, Header, Validation};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    /// Subject (whom token refers to).
    pub sub: String,
    /// Expiration time (as UTC timestamp).
    pub exp: i64,
}

/// Generate jwt that may be stored in a cookie.
///
/// The token will expire one day later.
///
/// You must validate password before calling this function.
pub fn generate_jwt(username: String, secret: &str) -> anyhow::Result<String> {
    let exp = (Utc::now() + Duration::days(1)).timestamp();
    let claims = Claims { sub: username, exp };
    let token = encode(
        &Header::default(),
        &claims,
        &EncodingKey::from_secret(secret.as_ref()),
    )?;
    Ok(token)
}

/// Returns true if jwt token is valid.
pub fn validate_jwt(token: &str, secret: &str) -> bool {
    decode::<Claims>(
        token,
        &DecodingKey::from_secret(secret.as_ref()),
        &Validation::new(Algorithm::HS256),
    ).is_ok()
}
