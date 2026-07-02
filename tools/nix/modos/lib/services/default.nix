{ ... }:
{
  # Service modules for `service-flake`.
  modules = {
    services.garage = import ./services/garage.nix;
    services.minio = import ./services/minio.nix;
  };
}
