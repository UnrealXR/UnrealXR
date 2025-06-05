{
  pkgs ? import <nixpkgs> { },
}: pkgs.mkShell {
  buildInputs = with pkgs; [
    # Runtime dependencies
    python3
    pciutils

    # evdi build dependencies
    libdrm
    linuxHeaders

    # raylib build dependencies
    libGL
    xorg.libXi
    xorg.libXcursor
    xorg.libXrandr
    xorg.libXinerama
    waylandpp
    libxkbcommon
  ];

  shellHook = ''
    export LD_LIBRARY_PATH="$PWD/evdi/library:$LD_LIBRARY_PATH"
    mkdir -p "$PWD/data/config" "$PWD/data/data"
    export UNREALXR_CONFIG_PATH="$PWD/data/config"
    export UNREALXR_DATA_PATH="$PWD/data/data"

    if [ ! -d ".venv" ]; then
      python3 -m venv .venv
    fi

    source .venv/bin/activate
  '';
}
