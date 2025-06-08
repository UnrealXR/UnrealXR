# UnrealXR

UnrealXR is a display multiplexer for the original Nreal Air (other devices may be supported).

## Development

1. Clone this repository with submodules: `git clone --recurse-submodules https://git.terah.dev/imterah/unrealxr.git`
2. Initialize the development shell: `nix-shell`
3. Create a virtual environment: `python3 -m venv .venv; source .venv/bin/activate`
4. Install Python dependencies: `pip install -r requirements.txt`

### Building `PyEvdi`

1. Build `libevdi`: `cd evdi/library; make -j$(nproc); cd ../..`
2. Build `PyEvdi`: `cd evdi/pyevdi; make -j$(nproc); make install; cd ../..`

### Building patched `PyRay`/`raylib`

1. Build the native version of raylib: `cd raylib-python-cffi/raylib-c; mkdir -p build/out; cd build; cmake  -DCUSTOMIZE_BUILD=ON -DSUPPORT_FILEFORMAT_JPG=ON -DSUPPORT_FILEFORMAT_FLAC=ON -DWITH_PIC=ON -DCMAKE_BUILD_TYPE=Release -DPLATFORM=DRM -DENABLE_WAYLAND_DRM_LEASING=ON -DSUPPORT_CLIPBOARD_IMAGE=ON -DCMAKE_INSTALL_PREFIX:PATH=$PWD/out ..; make install -j$(nproc)`
2. Build the Python version of raylib: `cd ../..; PKG_CONFIG_PATH_FOR_TARGET="$PKG_CONFIG_PATH_FOR_TARGET:$PWD/raylib-c/build/out/lib64/pkgconfig/" PKG_CONFIG_PATH="$PKG_CONFIG_PATH:$PWD/raylib-c/build/out/lib64/pkgconfig/" ENABLE_WAYLAND_DRM_LEASING=YES RAYLIB_PLATFORM=DRM python3 setup.py bdist_wheel; pip install dist/*.whl`
