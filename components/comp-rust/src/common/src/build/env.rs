use crate::log::impl_slog_value;

pub const BUILD_VERSION: &str =
    option_env!("QUITSH_COMPONENT_VERSION").unwrap_or("0.0.0");

// BuildType encodes the build type of the application.
#[derive(Debug, PartialEq, strum_macros::EnumString, strum_macros::Display)]
#[strum(serialize_all = "camelCase")]
pub enum EnvironmentType {
    Development,
    Testing,
    Staging,
    Production,
}

#[derive(Debug, PartialEq, strum_macros::EnumString, strum_macros::Display)]
#[strum(serialize_all = "camelCase")]
pub enum BuildType {
    Debug,
    Release,
}

// Reports if compiling release mode.
pub const fn is_release() -> bool {
    cfg!(not(debug_assertions))
}

// Reports if compiling in debug mode.
pub const fn is_debug() -> bool {
    !is_release()
}

// Report the build type.
pub const fn build_type() -> BuildType {
    if is_release() {
        BuildType::Release
    } else {
        BuildType::Debug
    }
}

// Reports the environment type. Should only activate one feature.
pub const fn env_type() -> EnvironmentType {
    if cfg!(feature = "production") {
        EnvironmentType::Production
    } else if cfg!(feature = "staging") {
        EnvironmentType::Staging
    } else if cfg!(feature = "testing") {
        EnvironmentType::Testing
    } else {
        EnvironmentType::Development
    }
}

const fn count_enabled(flags: &[bool]) -> usize {
    let mut i = 0;
    let mut count = 0;
    while i < flags.len() {
        if flags[i] {
            count += 1;
        }
        i += 1;
    }
    count
}

const ENABLED_COUNT: usize = count_enabled(&[
    cfg!(feature = "development"),
    cfg!(feature = "testing"),
    cfg!(feature = "staging"),
    cfg!(feature = "production"),
]);

const _: () = assert!(
    ENABLED_COUNT == 1,
    "Exactly one of: development, testing, staging, production must be enabled.",
);

impl_slog_value!(BuildType, EnvironmentType);
