use std::{
    path::{Path, PathBuf},
    sync::Arc,
};

use async_compression::tokio::bufread::GzipDecoder;
use axum::{
    Json,
    body::Body,
    extract::{Request, State},
    http::{Response, StatusCode, header},
};
use common::{
    fs,
    response::{
        ErrorBadRequestCtx, ErrorInternalCtx, HttpError, ResponseError,
    },
    tracing::{ResultExt, err},
};
use futures_util::TryStreamExt;
use pdf_render::render;
use tempdir::TempDir;
use tokio::{fs::File, task::spawn_blocking};
use tokio_util::io::{ReaderStream, StreamReader, SyncIoBridge};
use tracing::{self, info};
use utoipa::{OpenApi, ToSchema};

use crate::{openapi::ApiDoc, state::AppState};

// Serves the OpenAPI specification.
pub async fn serve_openapi() -> Json<utoipa::openapi::OpenApi> {
    Json(ApiDoc::openapi())
}

#[derive(ToSchema)]
#[schema(format = Binary, value_type = String)]
#[allow(dead_code)]
struct BinaryBody(Vec<u8>);

// Health endpoint which tells if the service is up and running.
#[utoipa::path(
    get,
    path = "/health",
    description = "Check if service is healthy.",
    responses(
        (status = StatusCode::OK, description = "If the service is ready.")
    ),
)]
pub async fn health() -> StatusCode {
    StatusCode::OK
}

// Renders a PDF by extracting a payload folder and
// running it with pandoc.
//
// # Examples
//
// ```shell
// just render-example
// ```
#[utoipa::path(
    post,
    path = "/render",
    description = "Render a PDF.",
    request_body(
        content = BinaryBody,
        description = "Tar gzip archive containing the markdown content to render to a PDF.",
        content_type = "application/octet-stream"
    ),
    responses(
        (status = StatusCode::OK,
            description = "Successfully rendered the PDF.",
            body = BinaryBody,
            content_type = "application/pdf",
            headers(
                ("Content-Length" = u64, description = "The length of the rendered file."),
                ("Content-Disposition" = String, description = "The content disposition header."),
            ),
        ),
        (status = StatusCode::BAD_REQUEST,
            description = common::response::NAME_BAD_REQUEST,
            body = HttpError),
        (status = StatusCode::UNAUTHORIZED,
            description = common::response::NAME_UNAUTHORIZED,
            body = HttpError),
        (status = StatusCode::INTERNAL_SERVER_ERROR,
            description = common::response::NAME_SERVER_ERROR,
            body = HttpError),
        (status = "default", description = "Default Error", body = HttpError),
    ),
)]
pub async fn render(
    State(s): State<Arc<AppState>>,
    request: Request,
) -> Result<Response<Body>, ResponseError> {
    info!("Render PDF handler.");
    let stream = request
        .into_body()
        .into_data_stream()
        .map_err(std::io::Error::other);

    let stream_reader = StreamReader::new(stream);
    let gz_decoder = GzipDecoder::new(stream_reader);
    let sync_reader = SyncIoBridge::new(gz_decoder);

    let temp_dir = TempDir::new("pdf-render")
        .trace_err_m("Failed to create build. dir.")
        .map_err(|_| ErrorInternalCtx {}.build())?;

    let src_dir = temp_dir.path().to_owned().join("src");
    let src_dir2 = src_dir.clone();
    let out_dir = src_dir.join("out");
    fs::create_dirs([&src_dir, &out_dir], 0o700).map_err(|e| {
        err!(e, "Could not create directories.");
        ErrorInternalCtx {}.build()
    })?;

    info!("Unpacking payload into '{:?}'.", temp_dir.path());
    spawn_blocking(move || {
        let mut archive = tar::Archive::new(sync_reader);
        archive
            .unpack(&src_dir)
            .trace_err_m("Unpacking archive failed.")
            .map_err(|_| {
                ErrorBadRequestCtx {
                    message: "Payload malformed, cannot unpack.",
                }
                .build()
            })
    })
    .await
    .trace_err_m("Failed to join spawned rendering process.")
    .unwrap_or_else(|_| Err(ErrorInternalCtx.build()))?;

    info!("Rendering ...");
    let output = spawn_blocking(move || call_render(&s, &src_dir2, &out_dir))
        .await
        .unwrap_or_else(|e| {
            err!(e, "Failed to join spawned rendering process.");
            Err(ErrorInternalCtx.build())
        })?;

    info!(output = %output.display(), "Streaming output file.");

    let file = File::open(output)
        .await
        .trace_err_m("failed to open output file")
        .map_err(|_| ErrorInternalCtx {}.build())?;
    let meta = file
        .metadata()
        .await
        .trace_err_m("failed to read meta data of output file")
        .map_err(|_| ErrorInternalCtx {}.build())?;

    let stream = ReaderStream::new(file);
    let body = Body::from_stream(stream);

    Response::builder()
        .header(header::CONTENT_TYPE, "application/pdf")
        .header(header::CONTENT_LENGTH, meta.len())
        .header(
            header::CONTENT_DISPOSITION,
            "attachment; filename=\"output.pdf\"",
        )
        .body(body)
        .trace_err_m("could not build response")
        .map_err(|_| ErrorInternalCtx {}.build())
}

fn call_render(
    s: &Arc<AppState>,
    src_dir: &Path,
    out_dir: &Path,
) -> Result<PathBuf, ResponseError> {
    render::to_pdf(
        src_dir,
        out_dir,
        &s.config.render.template_dir,
        &s.config.render.template_name_default,
    )
    .trace_err_m("Rendering failed.")
    .map_err(|_| {
        ErrorBadRequestCtx {
            message: "PDF rendering failed. See log.",
        }
        .build()
    })
}
