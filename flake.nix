{
  description = "Go development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go and core tools
            go
            gopls
            go-tools
            
            # Air for live reload
            air
            
            # Build tools
            gnumake
            
            # Database tools
            migrate
            
            # Additional development tools
            golangci-lint
            delve
            
            # Git tools
            git
          ];

          shellHook = ''
            # Set GOPATH to the current directory
            export GOPATH="$PWD/.go"
            export PATH="$GOPATH/bin:$PATH"
            
            # VSCode Go extension environment variables
            export GOROOT="$(go env GOROOT)"
            export GO111MODULE=on
            
            # Create necessary directories
            mkdir -p .go/bin
            
            # Install additional Go tools required for VSCode
            echo "Installing additional Go tools..."
            go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest
            go install github.com/ramya-rao-a/go-outline@latest
            go install github.com/cweill/gotests/gotests@latest
            go install github.com/fatih/gomodifytags@latest
            go install github.com/josharian/impl@latest
            go install honnef.co/go/tools/cmd/staticcheck@latest
            
            # Print available tools
            echo "Go development environment ready!"
            echo "Available tools:"
            echo "- go ($(go version))"
            echo "- gopls ($(gopls version))"
            echo "- air ($(air -v))"
            echo "- migrate ($(migrate -version))"
            echo "- make ($(make -v | head -n1))"
            echo "- golangci-lint ($(golangci-lint --version))"
          '';
        };
      }
    );
}
