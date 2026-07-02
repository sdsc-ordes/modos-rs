use std::sync::Arc;

use crate::config::Config;

pub struct AppState {
    pub config: Arc<Config>,
}
