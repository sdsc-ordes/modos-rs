use std::fs;

use common::{
    build::{EnvironmentType, env_type, is_debug},
    result::Res,
};
use snafu::ResultExt;
use utoipa::OpenApi;

/// Due to lint errors, I have to wrap into a module here.
#[allow(clippy::needless_for_each)]
pub mod api {
    use super::OpenApi;

    #[derive(OpenApi)]
    #[openapi(paths(crate::handler::health, crate::handler::render))]
    pub struct ApiDoc;
}
pub use api::ApiDoc;

/// Saves the openapi spec on start of the service
/// (only if debug or development environment)
pub fn save() -> Res<()> {
    if is_debug() || env_type() == EnvironmentType::Development {
        tracing::info!("Writing OpenAPI spec to 'api/openapi.json'.");

        let json = ApiDoc::openapi().to_pretty_json().with_whatever_context(
            |_| "Could not serialize OpenAPI JSON spec.",
        )?;

        let yaml = ApiDoc::openapi()
            .to_yaml()
            .with_whatever_context(|_| "Could not save OpenAPI spec.")?;

        fs::create_dir_all("api")
            .with_whatever_context(|_| "Could not create directory.")?;

        fs::write("api/openapi.json", &json)
            .with_whatever_context(|_| "Could not write OpenAPI spec")?;

        fs::write("api/openapi.yaml", &yaml)
            .with_whatever_context(|_| "Could not write OpenAPI spec")?;
    }

    Ok(())
}
