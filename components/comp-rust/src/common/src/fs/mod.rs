use std::{
    borrow::Cow, fs::DirBuilder, io, os::unix::fs::DirBuilderExt, path::Path,
};

// Makes a path absolute to `to` if its relative.
pub fn make_absolute_to<'a>(to: &Path, p: &'a Path) -> Cow<'a, Path> {
    if p.is_absolute() {
        Cow::Borrowed(p)
    } else {
        Cow::Owned(to.join(p))
    }
}

// Make directories recursive.
pub fn create_dirs<'a, P>(
    dirs: impl IntoIterator<Item = &'a P>,
    mode: u32,
) -> io::Result<()>
where
    P: AsRef<Path> + 'a,
{
    let mut b = DirBuilder::new();
    b.recursive(true).mode(mode);

    for d in dirs.into_iter() {
        b.create(d)?
    }

    Ok(())
}
