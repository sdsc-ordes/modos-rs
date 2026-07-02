pub mod error;
pub mod log;

pub use error::*;
pub use log::*;

/// The common error shared across the components.
pub type Res<T> = Result<T, Error>;
