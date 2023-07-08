{ pkgs }: {
    deps = [
        pkgs.just
        pkgs.bitcoind
        pkgs.clightning
        pkgs.go
        pkgs.gopls
    ];
}