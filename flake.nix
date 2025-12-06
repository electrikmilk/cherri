{
  description = "Cherri - A language that compiles to Apple Shortcuts";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        version =
          let
            content = builtins.readFile ./version.go;
            # Extract version string using regex
            matches = builtins.match ''.*version = "([^"]+)".*'' content;
          in
          builtins.head matches;
      in
      {
        packages = {
          cherri = pkgs.buildGoModule {
            pname = "cherri";
            version = version;

            src = ./.;

            vendorHash = "sha256-5hACp6B2eO82pv470E/dIMbKAO1Z/PI978pSzOB7Wtk=";

            doCheck = false;

            meta = with pkgs.lib; {
              description = "A language that compiles to Apple Shortcuts";
              homepage = "https://github.com/electrikmilk/cherri";
              license = licenses.gpl2;
              maintainers = [ ];
              mainProgram = "cherri";
            };
          };

          default = self.packages.${system}.cherri;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
          ];
        };
      }
    );
}
