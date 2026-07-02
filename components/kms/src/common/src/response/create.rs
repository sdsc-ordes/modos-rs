use std::sync::Arc;

use axum::{
    extract::FromRequest,
    http::StatusCode,
    response::{IntoResponse, Response},
};
use serde::Serialize;
use utoipa::ToSchema;

use crate::response::ResponseError;

pub const NAME_SERVER_ERROR: &str = "Internal Error";
pub const NAME_BAD_REQUEST: &str = "Bad Request";
pub const NAME_UNAUTHORIZED: &str = "Unauthorized";

pub const ID_SERVER_ERROR: &str = "internal";
pub const ID_BAD_REQUEST: &str = "bad-request";
pub const ID_UNAUTHORIZED: &str = "unauthorized";

// HttpError represents how we serialize the errors.
#[derive(Serialize, ToSchema)]
pub struct HttpError {
    pub id: String,
    pub title: String,
    pub message: String,
}

impl HttpError {
    fn unauthorized(message: &str) -> (StatusCode, Self) {
        (
            StatusCode::UNAUTHORIZED,
            HttpError {
                id: ID_UNAUTHORIZED.to_owned(),
                title: NAME_UNAUTHORIZED.to_owned(),
                message: message.to_owned(),
            },
        )
    }

    fn bad_request(message: &str) -> (StatusCode, Self) {
        (
            StatusCode::BAD_REQUEST,
            HttpError {
                id: ID_BAD_REQUEST.to_owned(),
                title: NAME_BAD_REQUEST.to_owned(),
                message: message.to_owned(),
            },
        )
    }

    fn internal() -> (StatusCode, Self) {
        (
            StatusCode::INTERNAL_SERVER_ERROR,
            HttpError {
                id: ID_SERVER_ERROR.to_owned(),
                title: NAME_SERVER_ERROR.to_owned(),
                message: NAME_SERVER_ERROR.to_owned(),
            },
        )
    }
}

// Create our own JSON extractor by wrapping `axum::Json`. This makes it easy to override the
// rejection and provide our own which formats errors to match our application.
//
// `axum::Json` responds with plain text if the input is invalid.
#[derive(FromRequest)]
#[from_request(via(axum::Json), rejection(ResponseError))]
pub struct ResponseJson<T>(pub T);

impl<T> IntoResponse for ResponseJson<T>
where
    axum::Json<T>: IntoResponse,
{
    fn into_response(self) -> Response {
        axum::Json(self.0).into_response()
    }
}

// Tell axum how `AppError` should be converted into a response.
impl IntoResponse for ResponseError {
    fn into_response(self) -> Response {
        let resp = match &self {
            ResponseError::Unauthorized { message, .. } => {
                HttpError::unauthorized(message)
            }
            ResponseError::BadRequest { message, .. } => {
                HttpError::bad_request(message)
            }
            _ => HttpError::internal(),
        };

        let mut response: Response =
            (resp.0, ResponseJson(resp.1)).into_response();

        // Set the error on the extension.
        response.extensions_mut().insert(Arc::new(self));

        response
    }
}
