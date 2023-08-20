{
  description = "A minimalistic image board engine written in Go.";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-23.05";

  outputs = { self, nixpkgs }:
    let
      version = "0.0.5";

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          microboard = pkgs.buildGoModule {
            nativeBuildInputs = [
              pkgs.pkg-config
            ];
            buildInputs = [
              pkgs.vips
            ];
            pname = "microboard";
            inherit version;
            src = ./.;

            # This hash locks the dependencies of this package. It is
            # necessary because of how Go requires network access to resolve
            # VCS.  See https://www.tweag.io/blog/2021-03-04-gomod2nix/ for
            # details. Normally one can build with a fake sha256 and rely on native Go
            # mechanisms to tell you what the hash should be or determine what
            # it should be "out-of-band" with other tooling (eg. gomod2nix).
            # To begin with it is recommended to set this, but one must
            # remeber to bump this hash when your dependencies change.
            #vendorSha256 = pkgs.lib.fakeSha256;
            vendorSha256 = "sha256-N7zJlTBr7FjaW1twfmL4Nprl5Prx35Hz2GkBQUEHvd0=";
          };
        });

      # The default package for 'nix build'. This makes sense if the
      # flake provides only one package or there is a clear "main"
      # package.
      defaultPackage = forAllSystems (system: self.packages.${system}.microboard);

      # NixOS service
      nixosModule = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        { config, lib, pkgs, ... }:
          with lib; let
            cfg = config.services.microboard;
          in
          {
            options.services.microboard = {
              enable = mkEnableOption "Enable microboard engine";

              port = mkOption {
                type = types.port;
                default = 55006;
                description = "Port to serve microboard on.";
              };

              dataDir = mkOption {
                type = types.path;
                default = "/var/lib/microboard";
                description = "Directory where uploaded images and attachments are stored.";
              };

              database.name = {
                type = types.str;
                default = "microboard";
                description = "Name of the postgres database";
              };

              database.user = {
                type = types.str;
                default = "microboard";
                description = "Postgres user";
              };

              database.passwordFile = {
                type = types.nullOr types.path;
                default = null;
                description = "Path to the file with the password for postgres.";
              };

              database.host = {
                type = types.str;
                default = "/run/postgresql";
                description = "Postgresql host";
              };

              database.port = {
                type = types.nullOr types.port;
                default = null;
                description = "Postgres port";
              };
            };
            config = mkIf cfg.enable {
              users.users.microboard = {
                group = "microboard";
                home = cfg.dataDir;
                createHome = true;
              };

              users.groups.microboard = { };

              systemd.services.microboard = {
                description = "Microboard engine";
                wantedBy = ["multi-user.target"];

                environment = {
                  MB_LOGLEVEL = "warning";
                  MB_UPLOADDIR = "${cfg.dataDir}/uploads";
                  MB_PREVIEWDIR = "${cfg.dataDir}/previews";
                  MB_DBHOST = cfg.database.host;
                  MB_DBUSER = cfg.database.user;
                  MB_DBNAME = cfg.database.name;
                  MB_PORT = cfg.port;
                } // lib.optionalAttrs (cfg.database.passwordFile != null) {
                  MB_DBPASSWORD = "$(cat ${cfg.database.passwordFile})";
                } // lib.optionalAttrs (cfg.databse.port != null) {
                  MB_DBPORT = cfg.database.port;
                };

                serviceConfig = {
                  User = cfg.user;
                  Group = cfg.user.group;
                  ExecStart = "${cfg.package}/bin/microboard";
                  Restart = "on-failure";
                  Type = "exec";
                  WorkingDirectory = cfg.dataDir;

                  # Security Hardening
                  LockPersonality = true;
                  NoNewPrivileges = true;
                  ProtectSystem = "strict";
                  ReadWritePaths = [ cfg.dataDir ];
                  RestrictAddressFamilies = [ "AF_INET" "AF_INET6" "AF_UNIX" ];
                  RestrictSUIDSGID = true;
                };
              };
            };
          });

    };
}
