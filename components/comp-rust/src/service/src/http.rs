use std::sync::Arc;

use axum::{
    Router,
    extract::{MatchedPath, Request},
    routing::{get, post},
};
use common::result::Res;
use snafu::{ResultExt, whatever};
use tower_http::trace::TraceLayer;
use tracing::info;

use crate::{
    handler::{health, render, serve_openapi},
    state::AppState,
};

pub async fn start_server(s: Arc<AppState>) -> Res<()> {
    let app = Router::new()
        .route("/health", get(health))
        .route("/render", post(render))
        .route("/api/openapi.json", get(serve_openapi))
        .layer(
            TraceLayer::new_for_http()
                // Create our own span for the request and include the matched path. The matched
                // path is useful for figuring out which handler the request was routed to.
                .make_span_with(|req: &Request| {
                    let method = req.method();
                    let uri = req.uri();

                    let matched_path = req
                        .extensions()
                        .get::<MatchedPath>()
                        .map(axum::extract::MatchedPath::as_str);

                    tracing::debug_span!("request", %method, %uri, matched_path)
                })
                // By default `TraceLayer` will log 5xx responses but we're doing our specific
                // logging of errors so disable that
                .on_failure(()),
        )
        .with_state(s.clone());

    let host = format!("{}:{}", s.config.server.hostname, s.config.server.port);
    info!("Starting server on {host}");

    let listener = tokio::net::TcpListener::bind(&host)
        .await
        .whatever_context(format!("could not create listener at {}", &host))?;

    match axum::serve(listener, app).await {
        Ok(()) => {
            info!("Successfully closed listener.");
            Ok(())
        }
        Err(e) => whatever!("axum listener failed: {e}"),
    }
}
