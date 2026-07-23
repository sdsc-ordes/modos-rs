use std::{ffi::OsString, path::PathBuf};

// Annoyingly `p.with_extension("sig")` will turn `foo.bar` into `foo.sig`. This
// avoids that issue (and without requiring the path be UTF-8), but is kind of
// tedious. Note that `ext` should be something like `"sig"` and not `".sig"`.
pub fn add_extension(p: impl Into<PathBuf>, ext: &str) -> PathBuf {
    let mut path: PathBuf = p.into();
    let mut name: OsString = path.file_name().unwrap_or_default().to_owned();
    name.push(".");
    name.push(ext);
    path.set_file_name(name);
    path
}

#[test]
fn check() {
    assert_eq!(&add_extension("asdf.bar", "yaml"), "asdf.bar.yaml");
    assert_eq!(&add_extension("asdf", "yaml"), "asdf.yaml");
}
