use std::fmt::Display;

use snafu::Report;

use crate::log::{Logger, error};

/// Extension trait to log result if [`Result`] is [`Err`].
pub trait ResultExt<T, E>: Sized {
    fn log(self, log: &Logger) -> Result<T, E>;
}

impl<T, E> ResultExt<T, E> for Result<T, E>
where
    E: Display + std::error::Error,
{
    fn log(self, log: &Logger) -> Result<T, E> {
        if let Err(err) = self {
            error!(log, "{}", Report::from_error(&err));
            return Err(err);
        }

        self
    }
}
