{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
  outputs = { self, nixpkgs }:
    let
      pkgs = import nixpkgs {
        system = "aarch64-linux";
      };
    in
    {
      packages.aarch64-linux.rpi4-motd-panel = pkgs.buildGoModule {
        pname = "rpi4-motd-panel";
        version = "0.1";
        src = ./.;
        vendorHash = "sha256-KH67bTaRKjaM3a4j38G4TdJejfv8FlAHfFtlhMW1oUQ=";
      };

      packages.aarch64-linux.default = self.packages.aarch64-linux.rpi4-motd-panel;
    };
}
