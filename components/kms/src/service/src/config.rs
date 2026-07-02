use std::{
    fs, io,
    path::{Path, PathBuf},
};

use common::result::Res;
use serde::{Deserialize, Serialize};
use snafu::prelude::*;

#[derive(
    Debug, Default, Clone, Serialize, Deserialize, PartialEq, PartialOrd,
)]
#[serde(default, rename_all = "camelCase", deny_unknown_fields)]
pub struct Config {
    pub server: Server,
    pub render: Render,
    pub log: Log,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, PartialOrd)]
#[serde(default, rename_all = "camelCase", deny_unknown_fields)]
pub struct Server {
    pub hostname: String,
    pub port: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, PartialOrd)]
#[serde(default, rename_all = "camelCase", deny_unknown_fields)]
pub struct Render {
    pub template_dir: PathBuf,
    pub template_name_default: String,
}

#[derive(
    Debug, Default, Clone, Serialize, Deserialize, PartialEq, PartialOrd,
)]
#[serde(default, rename_all = "camelCase")]
pub struct Log {
    pub force_dev_log: bool,
}

pub trait WithDataDir {
    fn with_data_dir(&mut self, dir: &Path);
}

impl Config {
    // Load the config from directory `dir`.
    pub fn load(dir: &Path) -> Res<Box<Config>> {
        let file = dir.join("config.yaml");

        tracing::info!("Loading config from {}.", file.display());

        let f = fs::File::open(&file)
            .whatever_context(format!("failed to load {}", &file.display()))?;

        match serde_yaml::from_reader(f) {
            Ok(c) => Ok(Box::new(c)),
            Err(e) => Err(io::Error::other(e.to_string())).whatever_context(
                format!("failed to read {}", &file.display()),
            ),
        }
    }
}

impl WithDataDir for Config {
    fn with_data_dir(&mut self, dir: &Path) {
        self.render.with_data_dir(dir);
    }
}

impl WithDataDir for Render {
    fn with_data_dir(&mut self, dir: &Path) {
        if self.template_dir.is_relative() {
            self.template_dir = dir.join(&self.template_dir);
        }
    }
}

impl Default for Server {
    fn default() -> Self {
        Server {
            hostname: "127.0.0.1".into(),
            port: 3040,
        }
    }
}

impl Default for Render {
    fn default() -> Self {
        Render {
            template_dir: "templates".into(),
            template_name_default: "default".into(),
        }
    }
}
