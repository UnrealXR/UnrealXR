{
  pkgs ? import <nixpkgs> { },
}: pkgs.mkShell {
  buildInputs = with pkgs; [
    # Runtime dependencies
    pciutils

    # UnrealXR build dependencies
    go
    gopls

    # evdi build dependencies
    libdrm
    linuxHeaders

    # raylib build dependencies
    cmake
    clang-tools
    pkg-config
    wayland
    libGL
    libgbm
    libdrm
    xorg.libXi
    xorg.libXcursor
    xorg.libXrandr
    xorg.libXinerama
    xorg.libX11

    # raylib X11 build dependencies (for dev-only non-XR builds)
    libxkbcommon

    # nreal driver build dependencies
    hidapi
    json_c
    libusb1
  ];

  shellHook = ''
    export LD_LIBRARY_PATH="$PWD/modules/evdi/library:${pkgs.lib.makeLibraryPath [ pkgs.libGL ]}:$LD_LIBRARY_PATH"
    mkdir -p "$PWD/data/config" "$PWD/data/data"
    export UNREALXR_CONFIG_PATH="$PWD/data/config"
    export UNREALXR_DATA_PATH="$PWD/data/data"
    export UNREALXR_LOG_LEVEL="debug"
  '';
}
