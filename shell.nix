{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  name = "dev-shell";
  buildInputs = with pkgs; [
    go
    ngrok
  ];
}
