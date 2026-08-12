# Development Guide

### Nix Cache Setup

Make sure that your daemon nix config at: `/etc/nix/nix.conf` (or
`/etc/nix/nix.custom.conf`) contains:

```ini
extra-trusted-user = <user-name>
```

to enable substituters in `flake.nix`.
