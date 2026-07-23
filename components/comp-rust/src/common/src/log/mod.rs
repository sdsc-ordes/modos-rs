use std::{io, sync::Arc};

use slog::{Drain, o};
use slog_async;

use crate::build::{BuildType, EnvironmentType, env_type};

/// Wrapping our internal type to the outside.
/// TODO: Wrap it better, is a struct with private member possible?
pub type Logger = slog::Logger;

#[allow(dead_code)]
fn no_out(_io: &mut dyn io::Write) -> io::Result<()> {
    Ok(())
}

#[derive(Default, Copy, Clone, PartialEq, PartialOrd)]
pub struct LoggerBuilder {
    pub timestamp: bool,
    pub force_dev_log: bool,
}

// Create a logger.
pub fn new() -> LoggerBuilder {
    LoggerBuilder::default()
}

impl LoggerBuilder {
    pub fn with_timestamp(&mut self, enable: bool) -> &mut LoggerBuilder {
        self.timestamp = enable;
        self
    }

    pub fn with_force_devlog(&mut self, enable: bool) -> &mut LoggerBuilder {
        self.force_dev_log = enable;
        self
    }

    pub fn build(&self) -> Arc<Logger> {
        let drain;

        if [EnvironmentType::Development, EnvironmentType::Testing]
            .contains(&env_type())
            || self.force_dev_log
        {
            let decorator = slog_term::TermDecorator::new().build();
            let mut d1 = slog_term::FullFormat::new(decorator);

            if !self.timestamp {
                d1 = d1.use_custom_timestamp(no_out);
            }

            let d = d1.build().fuse();
            drain = slog_async::Async::new(d)
                .chan_size(5_000_000)
                .build()
                .fuse();
        } else {
            let d = slog_json::Json::new(std::io::stderr())
                .set_pretty(build_type() == BuildType::Debug)
                .add_default_keys()
                .set_flush(true)
                .build()
                .fuse();

            drain = slog_async::Async::new(d)
                .chan_size(5_000_000)
                .build()
                .fuse();
        }

        Arc::new(slog::Logger::root(drain, o!()))
    }
}

/// Log trace level record
#[macro_export]
macro_rules! log_trace(
    ($log:expr, #$tag:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Trace, $tag, $($args)+)
    };
    ($log:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Trace, "", $($args)+)
    };
);

pub use log_trace as trace;

/// Log debug level record
#[macro_export]
macro_rules! log_debug(
    ($log:expr, #$tag:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Debug, $tag, $($args)+)
    };
    ($log:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Debug, "", $($args)+)
    };
);

pub use log_debug as debug;

/// Log info level record
#[macro_export]
macro_rules! log_info(
    ($log:expr, #$tag:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Info, $tag, $($args)+)
    };
    ($log:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Info, "", $($args)+)
    };
);

pub use log_info as info;

/// Log warn level record
#[macro_export]
macro_rules! log_warn(
    ($log:expr, #$tag:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Warning, $tag, $($args)+)
    };
    ($log:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Warning, "", $($args)+)
    };
);

pub use log_warn as warn;

/// Log warn level record
#[macro_export]
macro_rules! log_error(
    ($log:expr, #$tag:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Error, $tag, $($args)+)
    };
    ($log:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Error, "", $($args)+)
    };
);

pub use log_error as error;

#[macro_export]
macro_rules! log_critical(
    ($log:expr, #$tag:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Critical, $tag, $($args)+)
    };
    ($log:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Critical, "", $($args)+)
    };
);

pub use log_critical as critical;

/// Log panic level record
#[macro_export]
macro_rules! log_panic(
    ($log:expr, #$tag:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Error, $tag, $($args)+);
        panic!();
    };
    ($log:expr, $($args:tt)+) => {
        slog::log!($log, slog::Level::Error, "", $($args)+);
        panic!();
    };
);

pub use log_panic;

// Implements the slog::Value trait for enums with `to_string()`.
#[macro_export]
macro_rules! impl_slog_value {
    ($($t:ty),*) => {
        $(
            impl slog::Value for $t {
                fn serialize(
                    &self,
                    _record: &slog::Record,
                    key: slog::Key,
                    serializer: &mut dyn slog::Serializer,
                ) -> slog::Result {
                    serializer.emit_str(key, &self.to_string())
                }
            }
        )*
    };
}

pub use impl_slog_value;

use crate::build::build_type;
