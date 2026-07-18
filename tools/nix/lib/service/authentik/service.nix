{
  config,
  lib,
  pkgs,
  ...
}:
let
  name = "authentik";
  cfg = config.services.authentik;
  dataDir = cfg.dataDir + "/authentik";

  settingsFormat = pkgs.formats.yaml { };

  # The authentik components.
  gopkgs = requireComponent "gopkgs";
  rust = requireComponent "rust";
  rawMigrate = requireComponent "migrate";
  pythonEnv = requireComponent "pythonEnv";
  rawStaticWorkdirDeps = requireComponent "staticWorkdirDeps";

  requireComponent =
    attr:
    if cfg.components ? ${attr} then
      cfg.components.${attr}
    else
      throw ''
        services.authentik.${name}.authentikComponents is missing the `${attr}` attribute.
        (see the module docs).
      '';

  # NOTE: Workaround -> authentik-nix's `migrate` ships `bin/migrate.py`, a shell wrapper that execs an inner
  # `bin/.migrate.py-wrapped` whose shebang is `#!/usr/bin/env python`.
  # `/usr/bin/env` does not exist in the Nix build sandbox that
  # `nix build .#checks…` (`just test`) runs the process-compose stack in,
  # so the migration crashes with "bad interpreter".
  migrate = pkgs.runCommand "authentik-migrate-sandbox-safe" { } ''
    cp -R --no-preserve=mode,ownership ${rawMigrate} $out
    substituteInPlace $out/bin/.migrate.py-wrapped \
      --replace-fail '#!/usr/bin/env python' '#!${pythonEnv}/bin/python'
    substituteInPlace $out/bin/migrate.py \
      --replace-fail '${rawMigrate}/bin/.migrate.py-wrapped' "$out/bin/.migrate.py-wrapped"
    chmod +x "$out/bin/.migrate.py-wrapped" "$out/bin/migrate.py"
  '';

  # Fix upstream `except` syntax in upstream worker_process.py so it parses under Python 3.
  staticWorkdirDeps = pkgs.runCommand "authentik-static-workdir-deps-patched" { } ''
    cp -LR --no-preserve=mode,ownership ${rawStaticWorkdirDeps} $out



    # chmod +w -R $out/lifecycle
    substituteInPlace $out/lifecycle/worker_process.py \
      --replace-fail 'except OSError, FileNotFoundError:' 'except (OSError, FileNotFoundError):'
  '';

  connEnv = {
    AUTHENTIK_SECRET_KEY = cfg.secretKey;
    AUTHENTIK_BOOTSTRAP_PASSWORD = cfg.initialAdminPassword;
    AUTHENTIK_BOOTSTRAP_EMAIL = cfg.initialAdminEmail;
  };

  listenEnv =
    type:
    let
      c = cfg.${type};
    in
    {
      AUTHENTIK_LISTEN__HTTP = "${c.http.host}:${toString c.http.port}";
      AUTHENTIK_LISTEN__METRICS = "${c.metrics.host}:${toString c.metrics.port}";
    }
    // lib.optionalAttrs (c ? https) {
      AUTHENTIK_LISTEN__HTTPS = "${c.https.host}:${toString c.https.port}";
    };

  # Copy the enabled blueprints into the Authentik blueprints dir.
  # Copied (not symlinked) for the same reason as the built-in
  # blueprints: authentik rejects blueprints whose resolved path escapes `blueprints_dir`.
  blueprintImport = lib.mapAttrsToList (
    bp: e:
    # Bash
    ''
      f=$(realpath "${e.path}")
      echo "Copying blueprint '${bp}' from '$f' into './blueprints'."
      if [ ! -f "$f" ]; then
        echo "Blueprint file '$f' does not exist!" >&2
        exit 1
      fi
      cp -L "$f" "./blueprints/"
      chmod u+w "./blueprints/$(basename "$f")"
      unset f
    '') (lib.filterAttrs (_: v: v.import && v.path != null) cfg.blueprints);

  loadEnvFile =
    lib.optionalString (cfg.environmentFile != null)
      # Bash
      ''
        echo "Load extra authentic environment file '${cfg.environmentFile}'."
        set -a
        # shellcheck disable=SC1091
        . "${cfg.environmentFile}"
        set +a
      '';

  runtimeEnv =
    # Bash
    ''
      dataDir="$(realpath ${dataDir})"
      echo "Authentik data dir: '$dataDir'."

      mkdir -p "$dataDir/data" "$dataDir/media" "$dataDir/certs" "$dataDir/prometheus"

      export PROMETHEUS_MULTIPROC_DIR="$dataDir/prometheus"

      # The server's Python imports `authentik` from a read-only store path and
      # loads its config relative to that (`__file__`-based base_dir in
      # authentik/lib/config.py), so it never reads the merged default.yml in the
      # data dir. Force the blueprints dir via env var, which config loading
      # applies last (overrides the file value) regardless of which store the
      # Python was imported from. Points at the read-only, patched blueprints
      # (contains system/bootstrap.yaml); custom blueprint import is not supported.
      export AUTHENTIK_BLUEPRINTS_DIR="${staticWorkdirDeps}/blueprints"
      echo "Bluepints directory: $AUTHENTIK_BLUEPRINTS_DIR"

      export TMPDIR="$dataDir/.temp"
      export TEMPDIR="$TMPDIR"

      export PATH="${pythonEnv}/bin:$PATH"


      cd "$dataDir"
    '';

  settingsFile = settingsFormat.generate "authentik.yml" cfg.settings;

  setup =
    # Bash
    ''
      set -euo pipefail

      dataDir="$(realpath ${dataDir})"
      export CUSTOM_TMPDIR=$(mktemp -d)
      mkdir -p "$CUSTOM_TMPDIR" "$dataDir"
      ln -s "$CUSTOM_TMPDIR" "$dataDir/.temp"

      ${runtimeEnv}
      ${loadEnvFile}

      echo "Working dir: $(pwd)"

      # Bring in Authentik's working-directory dependencies
      # (authentik/, templates/, static assets, ...).
      cp -R --no-preserve=mode,ownership "${staticWorkdirDeps}/." ./
      ${builtins.concatStringsSep "\n" blueprintImport}
      # chmod -R -w ./blueprints

      src="$dataDir/authentik/lib/default.yml"
      echo "Merging settings file into '$src'."
      ${lib.getExe pkgs.yq-go} eval-all '. as $item ireduce ({}; . *+ $item)' \
        "$src" "${settingsFile}" > $dataDir/.merged.yml
      mv "$src" "$src.original"
      mv "$dataDir/.merged.yml" "$src"
      echo "Settings file '$src':"
      echo "====================="
      cat "$src"
      echo "====================="
    '';

  authentik-migrate =
    pkgs.writeShellScriptBin "authentik-migrate"
      # Bash
      ''
        ${setup}
        echo "Starting authentik migrate ..."
        exec ${migrate}/bin/migrate.py
      '';

  authentik-worker =
    pkgs.writeShellScriptBin "authentik-worker"
      # Bash
      ''
        ${runtimeEnv}
        ${loadEnvFile}
        echo "Starting authentik worker ..."
        exec ${rust}/bin/authentik worker
      '';

  authentik-server =
    pkgs.writeShellScriptBin "authentik-server"
      # Bash
      ''
        ${runtimeEnv}
        ${loadEnvFile}
        echo "Starting authentik server ..."
        exec ${gopkgs}/bin/server
      '';

  authentik-health =
    pkgs.writeShellScriptBin "authentik-health"
      # Bash
      ''
        ${pkgs.curl}/bin/curl -fsS "http://${cfg.server.http.host}:${toString cfg.server.http.port}/-/health/ready/"
      '';
in
{
  services.authentik.services = {
    postgres = {
      "${name}-pg-db" = {
        enable = true;
        port = cfg.postgres.port;
        initialScript.before = ''
          CREATE USER \"${cfg.postgres.user}\" WITH PASSWORD '${cfg.postgres.password}' CREATEDB;
          CREATE DATABASE \"${cfg.postgres.name}\" OWNER \"${cfg.postgres.user}\"
        '';
      };
    };
  };

  services.authentik.settings = {
    log_level = lib.mkDefault cfg.logLevel;

    blueprints_dir = lib.mkDefault "${dataDir}/blueprints";

    templates_dir = lib.mkDefault "${staticWorkdirDeps}/templates";

    cert_discovery_dir = lib.mkDefault "${dataDir}/certs";

    postgresql = {
      user = lib.mkDefault cfg.postgres.user;
      name = lib.mkDefault cfg.postgres.name;
      host = lib.mkDefault cfg.postgres.host;
      port = lib.mkDefault cfg.postgres.port;
    };

    storage = {
      file = lib.mkDefault {
        path = "${dataDir}/data";
      };

      media = {
        backend = lib.mkDefault "file";
        file = lib.mkDefault {
          path = "${dataDir}/media";
        };
      };
    };
  };

  settings.processes = lib.mkIf cfg.enable {
    "${name}-migrate" = {
      description = "Authentik database migrations: a prerequisite that creates the schema.";
      environment = connEnv;
      command = "${lib.getExe authentik-migrate}";

      depends_on = {
        "${name}-pg-db".condition = "process_healthy";
      };
    };

    "${name}-worker" = {
      description = "Authentik background worker: processes tasks and applies blueprints.";
      environment = connEnv // listenEnv "worker";
      command = "${lib.getExe authentik-worker}";
      depends_on."${name}-migrate".condition = "process_completed_successfully";
    };

    "${name}" = {
      description = "Authentik HTTP/API server.";
      environment = connEnv // listenEnv "server";
      command = "${lib.getExe authentik-server}";

      depends_on = {
        "${name}-migrate".condition = "process_completed_successfully";
        "${name}-worker".condition = "process_started";
      };

      readiness_probe = {
        exec.command = "${lib.getExe authentik-health}";
        initial_delay_seconds = 20;
        period_seconds = 5;
        timeout_seconds = 4;
        failure_threshold = 30;
      };
    };
  };
}
