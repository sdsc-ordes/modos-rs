{ lib, ... }:
let
  # The same as `read` but only reads some toplevel single-line attributes.
  # This is only to avoid (IFD).
  readSimple =
    path: keywords:
    let
      lines = lib.splitString "\n" (lib.readFile path);

      findKeyword =
        keyword:
        let
          res = builtins.match "${keyword}: (.*)" (lib.findFirst (l: lib.hasPrefix "${keyword}:" l) "" lines);
        in
        assert lib.assertMsg (res != null) "Keyword '${keyword}' not found!";
        builtins.elemAt res 0;
    in
    lib.foldl' lib.mergeAttrs { } (
      lib.map (k: {
        "${k}" = findKeyword k;
      }) keywords
    );
in
{
  inherit readSimple;
}
