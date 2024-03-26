{
  description = "A Nix flake to configure dependencies and development environment";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go # The Go / Golang programming language
            gopls # Official language server for the Go language
            openssl # Cryptographic library
            python311 # Python 3.11
            python311Packages.pip # Python package manager
            python311Packages.virtualenv # Virtual environment module
          ];
          shellHook = ''
            # deactivate any currently activated virtual env and
            # silence output in case there is no virtual env
            deactivate > /dev/null 2> /dev/null
            echo "Checking for existing virutal environment..."
            if [ ! -d ".venv" ]; then
              echo "No existing virtual environment found."
              echo "Creating new virtual environment..."
              python -m venv .venv
            fi
            # activate virtual environment
            echo "Activating virtual environment..."
            source .venv/bin/activate
            echo "Checking for exisiting dependencies..."
            if [ ! -d ".venv/lib/python3.11/site-packages/bitcoin" ]; then
              echo "No existing dependencies found."
              echo "Installing dependencies..."
              pip install -r requirements.txt
            fi
            echo "Done."
          '';
        };
      });
}
