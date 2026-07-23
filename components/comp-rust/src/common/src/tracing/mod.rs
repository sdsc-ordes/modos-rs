pub mod result;
pub use result::*;
use tracing_subscriber::{
    EnvFilter,
    fmt::{
        self,
        time::{FormatTime, UtcTime},
    },
    layer::SubscriberExt,
    util::SubscriberInitExt,
};

use crate::build::{EnvironmentType, env_type};

#[derive(Default, Copy, Clone, PartialEq, PartialOrd)]
pub struct TraceBuilder {
    pub timestamp: bool,
    pub force_dev_log: bool,
}

// Create a tracing seutp.
pub fn new() -> TraceBuilder {
    TraceBuilder::default()
}

enum MaybeTimer {
    On(UtcTime<time::format_description::well_known::Rfc3339>),
    Off,
}

impl FormatTime for MaybeTimer {
    fn format_time(
        &self,
        w: &mut tracing_subscriber::fmt::format::Writer<'_>,
    ) -> std::fmt::Result {
        match self {
            MaybeTimer::On(t) => t.format_time(w),
            MaybeTimer::Off => Ok(()), // write nothing
        }
    }
}

impl TraceBuilder {
    fn make_timer(&self) -> MaybeTimer {
        if self.timestamp {
            MaybeTimer::On(UtcTime::rfc_3339())
        } else {
            MaybeTimer::Off
        }
    }

    pub fn with_timestamp(&mut self, enable: bool) -> &mut Self {
        self.timestamp = enable;
        self
    }

    pub fn with_force_devlog(&mut self, enable: bool) -> &mut Self {
        self.force_dev_log = enable;
        self
    }

    pub fn setup(&self) {
        if [EnvironmentType::Development, EnvironmentType::Testing]
            .contains(&env_type())
            || self.force_dev_log
        {
            tracing_subscriber::registry()
                .with(EnvFilter::new("trace"))
                .with(
                    fmt::layer().with_ansi(true).with_timer(self.make_timer()),
                )
                .init();
        } else {
            tracing_subscriber::registry()
                .with(EnvFilter::new("trace"))
                .with(fmt::layer().json().with_timer(self.make_timer()))
                .init();
        }
    }
}

// Error macro to trace an error.
#[macro_export]
macro_rules! err {
    ($err:expr) => {
        tracing::error!(error = %snafu::Report::from_error($err))
    };
    ($err:expr, $($rest:tt)*) => {
        tracing::error!(error = %snafu::Report::from_error($err), $($rest)*)
    };
}

pub use err;
