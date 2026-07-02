use std::{path::PathBuf, sync::Arc};

use clap::Parser;
use common::{
    build::{build_type, env_type},
    result::Res,
};
use service::{
    config::{Config, WithDataDir},
    http, openapi,
    state::AppState,
};
use tracing::info;

/// A simple CLI that says hello.
#[derive(Parser)]
#[command(name = "pdf-renderer")]
#[command(version = "0.1.0")]
#[command(about ="The pdf-renderer to render PDFs by pandoc templates.", long_about = None)]
struct ServiceArgs {
    /// The data directory for the service.
    #[clap(short, long, default_value = "data")]
    data_dir: PathBuf,

    /// The config directory for the service.
    #[clap(short, long, default_value = "config")]
    config_dir: PathBuf,
}

#[tokio::main]
#[snafu::report]
async fn main() -> Res<()> {
    let args = ServiceArgs::parse();
    let mut config = Config::load(&args.config_dir)?;
    config.with_data_dir(&args.data_dir);
    openapi::save()?;

    common::tracing::new()
        .with_timestamp(true)
        .with_force_devlog(config.log.force_dev_log)
        .setup();

    let state = Arc::new(AppState {
        config: config.into(),
    });

    info!(
        "build" = build_type().to_string(),
        "env" = env_type().to_string(),
        "Starting pdf-renderer."
    );

    http::start_server(state).await
}
