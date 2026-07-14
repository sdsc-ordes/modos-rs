mkIf cfg.enable {
  assertions = [
    {
      assertion = cfg.database.type == "dev-mem" -> realmsExport == [ ];
      message = ''
        You cannot export realms with `realms.«name».export == true` when
        using `database.type == 'dev-mem'`, import however works.
        You can disable realms export with `exportRealms = true` globally.
      '';
    }
  ];

  services.keycloak.settings = mkMerge [
    {
      # We always enable http since we also use it to check the health.
      http-enabled = true;
      db = cfg.database.type;

      health-enabled = true;
      http-management-relative-path = "/";

      log-console-level = "info";
      log-level = "info";

      https-certificate-file =
        if providedSSLCerts then cfg.sslCertificate else "${dummyCertificates}/ssl-cert.crt";
      https-certificate-key-file =
        if providedSSLCerts then cfg.sslCertificateKey else "${dummyCertificates}/ssl-cert.key";
    }
  ];

  packages = [ keycloakBuild ];

  env = {
    KC_HOME_DIR = config.env.DEVENV_STATE + "/keycloak";
    KC_CONF_DIR = config.env.DEVENV_STATE + "/keycloak/conf";
    KC_TMP_DIR = config.env.DEVENV_STATE + "/keycloak/tmp";

    KC_BOOTSTRAP_ADMIN_USERNAME = "admin";
    KC_BOOTSTRAP_ADMIN_PASSWORD = "${lib.escapeShellArg cfg.initialAdminPassword}";
  };

  processes.keycloak =
    let
      keycloak-start = pkgs.writeShellScriptBin "keycloak-start" ''
        set -euo pipefail
        mkdir -p "$KC_HOME_DIR"
        mkdir -p "$KC_HOME_DIR/conf"
        mkdir -p "$KC_HOME_DIR/tmp"

        # Always remove the symlinks for the realm imports.
        rm -rf "$KC_HOME_DIR/data/import" || true
        mkdir -p "$KC_HOME_DIR/data/import"

        ln -fs ${keycloakBuild}/providers "$KC_HOME_DIR/"
        ln -fs ${keycloakBuild}/lib "$KC_HOME_DIR/"
        install -D -m 0600 ${confFile} "$KC_HOME_DIR/conf/keycloak.conf"

        echo "Keycloak config:"
        ${keycloakBuild}/bin/kc.sh show-config || true

        echo "Import realms (if any)..."
        ${builtins.concatStringsSep "\n" realmImport}
        echo "========================"

        echo "Start keycloak:"
        exec ${keycloakBuild}/bin/kc.sh start --optimized --import-realm
      '';
    in
    {
      ports.http.allocate = baseHttpPort;
      ports.https.allocate = baseHttpsPort;
      ports.management.allocate = baseManagementPort;
      exec = "${keycloak-start}/bin/keycloak-start";

      ready = {
        exec = "${keycloak-health}/bin/keycloak-health";
        initial_delay = 20;
        probe_timeout = 4;
        failure_threshold = 20;
      };
    };

  # Export a single realm.
  scripts.keycloak-realm-export = mkIf cfg.scripts.exportRealm {
    exec = "${keycloak-realm-export}/bin/keycloak-realm-export";
    description = ''
      Export a realm '$1' (first argument) from keycloak to location '$2' (second argument).
    '';
  };

  # Export all configured realms.
  scripts.keycloak-realm-export-all = mkIf (realmsExport != [ ]) {
    exec = "${keycloak-realm-export-all}/bin/keycloak-realm-export-all";
    description = ''
      Save the configured realms from keycloak, to back them up. You can run it manually.
    '';
  };

  # Process to start for exporting the above.
  processes.keycloak-realm-export-all = mkIf (realmsExport != [ ]) {
    start.enable = false;
    exec = "${keycloak-realm-export-all}/bin/keycloak-realm-export-all";
    after = [ "devenv:processes:keycloak@completed" ];
  };
}
