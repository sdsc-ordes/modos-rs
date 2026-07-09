{ config, ... }:
{
  # Make the `modos` function parameter available in flake-parts modules.
  _module.args.modos = config.modos;

  # Make the `modos'` function parameter available in `perSystem`.
  perSystem =
    {
      config,
      ...
    }:
    {
      _module.args.modos' = config.modos;
    };
}
