{
  cacert,
  dockerTools,
  # Own arguments.
  etcGroupAndPasswd,
  service,
  ...
}:
dockerTools.buildLayeredImage {
  name = "modos-rs/${service.pname}-service";
  tag = service.version;

  contents = [
    etcGroupAndPasswd
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
      "org.opencontainers.image.source" = "https://github.com/sdsc-ordes/modos-rs";
      "org.opencontainers.image.description" = service.meta.description;
      "org.opencontainers.image.license" = service.meta.license.shortName;
      "org.opencontainers.image.version" = service.version;
    };
    User = "non-root";
  };
}
