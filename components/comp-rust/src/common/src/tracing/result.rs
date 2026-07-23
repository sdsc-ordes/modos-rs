use std::fmt::Display;

use snafu::Report;
use tracing::error;

/// Extension trait to trace result if [`Result`] is [`Err`].
pub trait ResultExt<T, E>: Sized {
    fn trace_err(self) -> Result<T, E>;
    fn trace_err_m(self, msg: &str) -> Result<T, E>;
}

impl<T, E> ResultExt<T, E> for Result<T, E>
where
    E: Display + std::error::Error,
{
    fn trace_err(self) -> Result<T, E> {
        if let Err(err) = &self {
            error!(error = %Report::from_error(err));
        }

        self
    }
    fn trace_err_m(self, msg: &str) -> Result<T, E> {
        if let Err(err) = &self {
            error!(error = %Report::from_error(err), "{}", msg);
        }

        self
    }
}
