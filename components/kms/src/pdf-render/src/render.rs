use std::{
    fs::File,
    io::Read,
    path::{Path, PathBuf},
};

use common::{
    fs,
    result::{IOErrorCtx, Res},
};
use fs_extra::dir::CopyOptions;
use path_clean::PathClean;
use snafu::{ResultExt, prelude::*};
use tracing::{error, info};

#[derive(Default, Clone, Debug)]
struct Pandoc {
    // The build directory where all files to render reside.
    src_dir: PathBuf,
    main_file: PathBuf,

    // The template's source directory.
    template_src_dir: PathBuf,
    // The template's pandoc dir.
    template_pandoc_dir: PathBuf,

    lua_path: PathBuf,
    // Pandoc's data dir.
    data_dir: PathBuf,

    // The output file inside the build folder.
    output_file: PathBuf,
}

impl Pandoc {
    fn new(
        src_dir: &Path,
        output_dir: &Path,
        templates_dir: &Path,
        template_name: &str,
    ) -> Res<Box<Pandoc>> {
        if template_name == "pandoc" {
            whatever!("template name must not be named 'pandoc' [internal]")
        }

        // Make all absolute.
        let cwd = std::env::current_dir()
            .whatever_context("could not get current working dir")?;
        let src_dir = fs::make_absolute_to(&cwd, src_dir);
        let output_dir = fs::make_absolute_to(&cwd, output_dir);

        let templates_dir = fs::make_absolute_to(&cwd, templates_dir);

        let templ_dir: PathBuf = templates_dir.join(template_name).clean();
        ensure_whatever!(
            templ_dir.starts_with(&templates_dir),
            "resolved template dir is outside of '{}",
            templates_dir.display()
        );
        ensure_whatever!(
            templ_dir.exists(),
            "template directory '{}' does not exist",
            templ_dir.display(),
        );

        let main_file = src_dir.join("main.md");
        let template_src_dir = templ_dir.join("src");
        let template_pandoc_dir = templ_dir.join("pandoc");
        let output_file = output_dir.join("output.pdf");
        let data_dir = templates_dir.join("pandoc");
        let lua_path = data_dir.join("lua/?.lua;;");
        for f in [src_dir.as_ref(), output_dir.as_ref(), &data_dir, &main_file]
        {
            ensure_whatever!(f.exists(), "path '{}' not existing", f.display());
        }

        Ok(Pandoc {
            src_dir: src_dir.into(),
            main_file,

            template_src_dir,
            template_pandoc_dir,

            lua_path,
            data_dir,

            output_file,
        }
        .into())
    }

    fn sync_template_to_build(&self) -> Res<()> {
        sync_src_to_dest(&self.template_src_dir, &self.src_dir)
    }

    fn render(&self) -> Res<&Path> {
        let templ_defaults: Vec<String> =
            [self.template_pandoc_dir.join("defaults/main.yaml")]
                .into_iter()
                .filter(|p| p.exists())
                .flat_map(|p| {
                    ["--defaults".to_string(), p.to_string_lossy().to_string()]
                })
                .collect();

        let log_file = self.src_dir.join("output.log");
        let stdout = File::create(&log_file)
            .whatever_context("could not create log file")?;
        let stderr = stdout
            .try_clone()
            .whatever_context("could not clone log file")?;

        let mut p = std::process::Command::new("pandoc");
        p.current_dir(&self.src_dir)
            .arg("--standalone")
            .arg("--fail-if-warnings")
            .arg("--data-dir")
            .arg(&self.data_dir)
            .arg("--defaults=pandoc-dirs.yaml")
            .arg("--defaults=pandoc-general.yaml")
            .arg("--defaults=pandoc-filters.yaml")
            .arg("--defaults=pandoc-typst.yaml")
            .args(templ_defaults)
            .arg(&self.main_file)
            .arg("-o")
            .arg(&self.output_file)
            .env("LUA_PATH", &self.lua_path)
            .env("SOURCE_ROOT", &self.src_dir)
            .stdout(stdout)
            .stderr(stderr);

        info!(args = ?p.get_args(), "Invoking pandoc with:");

        let exit_code =
            p.status().whatever_context("pandoc invocation failed.")?;

        let mut log = String::new();
        File::open(&log_file)
            .whatever_context("could not open log file")?
            .read_to_string(&mut log)
            .whatever_context(format!(
                "could not read log file '{}'",
                log_file.display()
            ))?;

        if !exit_code.success() {
            error!(log = log, "Pandoc log.");
            whatever!(
                "pandoc failed with exit code '{}'",
                exit_code.code().unwrap_or(-1)
            )
        }

        Ok(&self.output_file)
    }
}

// Syncs over all files in `src` to `dest`.
fn sync_src_to_dest(src: &Path, dest: &Path) -> Res<()> {
    if !src.exists() {
        return Ok(());
    }

    ensure_whatever!(
        dest.exists(),
        "destination '{}' does not exist",
        dest.display()
    );

    let mut paths = vec![];
    for dir_entry in src
        .read_dir()
        .whatever_context("could not read source directory '{}'")?
    {
        let p = dir_entry.context(IOErrorCtx {
            message: "could not read directory",
        })?;
        paths.push(p.path());
    }

    info!("Syncing '{:?}' with '{}'", paths, dest.display());
    if let Err(e) = fs_extra::copy_items(
        &paths,
        dest,
        &CopyOptions::new().overwrite(true).skip_exist(false),
    ) {
        whatever!(
            "failed to sync '{}' to {}': {}",
            src.display(),
            dest.display(),
            e,
        )
    }

    Ok(())
}

// Renders a PDF with a given template
// WARNING: User controlled input is `template_name`.
pub fn to_pdf(
    src_dir: &Path,
    out_dir: &Path,
    templates_dir: &Path,
    template_name: &str,
) -> Res<PathBuf> {
    info!(
        "build-dir" = %src_dir.display(),
        "template-dir" = %templates_dir.display(),
        "template-name" = %template_name,
        "Rendering PDF with 'pandoc' -> 'typst'."
    );

    let pandoc = Pandoc::new(src_dir, out_dir, templates_dir, template_name)?;
    info!(paths = ?pandoc, "Render paths:");
    pandoc.sync_template_to_build()?;
    let output = pandoc.render()?;
    info!(output = %output.display(), "Rendering successful.");

    Ok(pandoc.output_file)
}
