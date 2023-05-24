//! Error structure for http responses.

use axum::{
    http::StatusCode,
    response::{IntoResponse, Json, Response},
};
use serde::ser::{Serialize, SerializeStruct, Serializer};
use std::fmt;

/// Simple error struct that can be returned to the user as a JSON.
#[derive(Debug, Clone)]
pub struct HttpError {
    pub status_code: StatusCode,
    pub message: String,
}

impl HttpError {
    /// Create a new error.
    pub fn new(status_code: StatusCode, message: &str) -> HttpError {
        HttpError {
            status_code,
            message: message.to_string(),
        }
    }
}

impl fmt::Display for HttpError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Code {}: {}", self.status_code.as_u16(), self.message)
    }
}

impl Serialize for HttpError {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut state = serializer.serialize_struct("HttpError", 2)?;
        state.serialize_field("status_code", &self.status_code.as_u16())?;
        state.serialize_field("message", &self.message)?;
        state.end()
    }
}

impl IntoResponse for HttpError {
    fn into_response(self) -> Response {
        (self.status_code, Json(self)).into_response()
    }
}
