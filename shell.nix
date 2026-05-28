{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  name = "dev-shell";
  buildInputs = with pkgs; [
    go
  ];
  go-migrate-pg = pkgs.go-migrate.overrideAttrs(oldAttrs: {
    tags = ["postgres"];
  });
}

# nix-shell -p go-migrate -I nixpkgs=https://github.com/NixOS/nixpkgs/archive/f62d6734af4581af614cab0f2aa16bcecfc33c11.tar.gz
