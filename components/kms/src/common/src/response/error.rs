use snafu::{Location, prelude::*};

#[derive(Debug, Snafu)]
#[snafu(visibility(pub))]
pub enum ResponseError {
    #[snafu(
        display("@{location}: HTTP Error: Unauthorized: {message}"),
        context(name(ErrorUnauthorizedCtx))
    )]
    Unauthorized {
        message: String,
        #[snafu(implicit)]
        location: Location,
    },

    #[snafu(
        display("@{location}: HTTP Error: Bad Request: {message}"),
        context(name(ErrorBadRequestCtx))
    )]
    BadRequest {
        message: String,
        #[snafu(implicit)]
        location: Location,
    },

    #[snafu(
        display("@{location}: HttpError: Internal"),
        context(name(ErrorInternalCtx))
    )]
    Internal {
        #[snafu(implicit)]
        location: Location,
    },
}
