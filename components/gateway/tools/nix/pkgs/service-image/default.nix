{
  cacert,
  dockerTools,
  # Own arguments.
  modos,
  service,
  ...
}:
dockerTools.buildLayeredImage {
  name = "modos-rs/${service.pname}-service";
  tag = service.version;

  contents = [
    modos.packages.image.etcGroupAndPasswd
    cacert
    service
  ];

  fakeRootCommands = ''
    mkdir -p workspace/data workspace/config tmp
    chown -R 1000:1000 workspace tmp
    chmod -R u+rw workspace tmp
  '';

  config = {
    Entrypoint = [ "${service}/bin/${service.pname}" ];
    WorkingDir = "/workspace";
    Volumes = {
      "/workspace/config" = { };
      "/workspace/data" = { };
    };
    Env = [
      "SSL_CERT_FILE=${cacert}/etc/ssl/certs/ca-bundle.crt"
    ];
    Labels = {
      "org.opencontainers.image.source" = "https://gitlab.com/data-custodian/dac-portal";
      "org.opencontainers.image.description" = service.meta.description;
      "org.opencontainers.image.license" = service.meta.license.shortName;
      "org.opencontainers.image.version" = service.version;
    };
    User = "non-root";
  };
}
