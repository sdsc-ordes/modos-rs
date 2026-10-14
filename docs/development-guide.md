# Development Guide

### Nix Cache Setup

Make sure that your daemon nix config at: `/etc/nix/nix.conf` (or
`/etc/nix/nix.custom.conf`) contains:

```ini
extra-trusted-user = <user-name>
```

to enable substituters in `flake.nix`.

### Updating Nixpkgs

Updating works with `cd tools/nix && nix flake update multiverse` and then
adjusting the pinned `nixpkgs` sets in
[`nixpkgs.parts.nix`](../tools/nix/lib/nixpkgs.parts.nix). Do a `direnv reload`
and fix the asserts.
