use snafu::{Backtrace, Location, prelude::*};

#[derive(Debug, Snafu)]
#[snafu(visibility(pub), context(suffix(ErrorCtx)))] // Sets the default visibility for these context selectors
pub enum Error {
    #[snafu(display("@{location}:: IO Error: {message}"))]
    IOError {
        message: String,
        source: std::io::Error,

        #[snafu(implicit)]
        backtrace: Backtrace,

        #[snafu(implicit)]
        location: Location,
    },
    #[snafu(whatever, display("Generic Error: {message}"))]
    GenericError {
        message: String,
        // Having a `source` is optional, but if it is present, it must
        // have this specific attribute and type:
        #[snafu(source(from(Box<dyn std::error::Error + Send + Sync>, Some)))]
        source: Option<Box<dyn std::error::Error>>,

        #[snafu(implicit)]
        location: Location,

        #[snafu(implicit)]
        backtrace: Backtrace,
    },
}
