use service::openapi;

fn main() {
    common::tracing::new()
        .with_timestamp(true)
        .with_force_devlog(true)
        .setup();

    openapi::save().expect("Could not save OpenAPI spec.");
}
