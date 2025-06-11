# UnrealXR

UnrealXR is a spatial multi-display renderer for the Xreal line of devices, enabling immersive, simultaneous viewing of multiple desktops and applications in 3D space.

## Why not Breezy Desktop?

1. **Limited desktop support**: Breezy only supports GNOME Wayland, whereas UnrealXR supports all desktops as long as they are Wayland and implement the [drm-lease-v1](https://wayland.app/protocols/drm-lease-v1) protocol. All major desktops aside from [Sway](https://swaywm.org/) and [Gamescope](https://wiki.archlinux.org/title/Gamescope) support this. [UnrealXR may add X11 support in the future](https://git.terah.dev/imterah/unrealxr/issues/1).
2. **Pricing**: UnrealXR will always be free and open source. The catch is that since Breezy costs money, it will likely have more of an "out-of-the-box" working experience, and additionally have more polish.

*However*, Breezy Desktop supports more AR headsets. We only support Xreal devices as that's all I have.

**TL;DR**: Breezy Desktop costs money and it has limited desktop support. However, due to its cost, it'll likely be more stable in the long term as it has funding. Additionally, Breezy Desktop supports more AR devices, while we only support Xreal devices.

## Why not the normal Nebula app?

1. **Incompatibility**: The normal Nebula app doesn't support Linux. Of course, we won't support Windows for the forseeable future, and macOS isn't on our roadmap. So this isn't really a fair point.
1. **Bloat**: The drivers for the Xreal line of devices are bundled into a streaming app, which is a little bit absurd. Instead, UnrealXR is just UnrealXR with no extra dependencies required.
2. **Closed Source**: The Nebula app and the normal SDK is closed source, so we don't really know what it's doing internally.

**TL;DR**: The normal Nebula app is only compatible with Windows and macOS. However, we only support Linux, so that isn't a fair comparison. The Nebula app is more bloated, and it is also closed source.

## Dependencies

Before anything, this depends on the `evdi` Linux kernel module. This is packaged in Debian-based distributions as `evdi-dkms`. If you already have DisplayLink drivers installed for their devices, you likely do not need to do this step. After installing this, please reboot your computer.

First, install the runtime dependencies. For Debian-based distros, the dependency list should be: `python3 python3-pip`

You'll need to install build dependencies after this. For Debian-based distros, the dependency list should be: `git build-essential libdrm libdrm-dev linux-headers-$(uname -r) cmake clang-tools pkg-config libwayland-client++1 libgl1-mesa-dev libglu1-mesa-dev libwayland-dev libxkbcommon-dev libhidapi-dev libjson-c-dev	libudev-dev libusb-1.0-0	libusb-1.0-0-dev libopencv-dev`

If you're using Nix/NixOS, use the `nix-shell` to enter the development environment.

After that, create a virtual environment for Python (done automatically in Nix): `python3 -m venv .venv`. Then, activate it: `source .venv/bin/activate`

Finally, install the Python dependencies: `pip install -r requirements.txt`

From there, you need to follow all the below steps if applicable to your current platform.
.
### Building patched `raylib` and `PyRay` (all platforms)

1. First, you need to build the native version of raylib. To do that, go inside the `modules/raylib-python-cffi/raylib-c` directory.
2. Then, make the build directories and go into them: `mkdir -p build/out; cd build`
3. Configure raylib: `cmake  -DCUSTOMIZE_BUILD=ON -DSUPPORT_FILEFORMAT_JPG=ON -DSUPPORT_FILEFORMAT_FLAC=ON -DWITH_PIC=ON -DCMAKE_BUILD_TYPE=Release -DPLATFORM=DRM -DENABLE_WAYLAND_DRM_LEASING=ON -DSUPPORT_CLIPBOARD_IMAGE=ON -DBUILD_EXAMPLES=OFF -DSUPPORT_SSH_KEYBOARD_RPI=OFF -DDISABLE_EVDEV_INPUT=ON -DCMAKE_INSTALL_PREFIX:PATH=$PWD/out ..`
4. Finally, build and install raylib: `make install -j$(nproc)`
5. After that, you need to build the Python bindings. To do that, go to the `modules/raylib-python-cffi` directory. Assuming you did everything correctly, you should be able to go 2 directories back (`../..`) to get there.
6. If you're on normal Linux and are not using Nix, do this command to build the package: `PKG_CONFIG_PATH="$PKG_CONFIG_PATH:$PWD/raylib-c/build/out/lib64/pkgconfig/" ENABLE_WAYLAND_DRM_LEASING=YES RAYLIB_PLATFORM=DRM python3 setup.py bdist_wheel`
7. If you are using Nix/NixOS, do this command to build the package: `PKG_CONFIG_PATH_FOR_TARGET="$PKG_CONFIG_PATH_FOR_TARGET:$PWD/raylib-c/build/out/lib64/pkgconfig/" ENABLE_WAYLAND_DRM_LEASING=YES RAYLIB_PLATFORM=DRM python3 setup.py bdist_wheel`
8. Finally, install the package: `pip install dist/*.whl`

### Building `PyEvdi` (Linux)

1. First, build the original libevdi. To start that, go inside the `modules/evdi/library` directory.
2. Then, build `libevdi`: `make -j$(nproc)`
3. After that, you need to build the Python bindings. To do that, go to the `modules/evdi/pyevdi` directory. Assuming you did everything correctly, you should be able to go a directory back (`../pyevdi`) to get there.
4. Then, build `PyEvdi`: `make -j$(nproc); make install`

### Building `nreal-driver` (all platforms)

1. First, create the `drivers` directory in the project root.
2. Then, go inside the `modules/nreal-driver` directory.
3. After that, make the build directories and go into them: `mkdir -p build/out; cd build`
4. Configure nreal-driver: `cmake -DCMAKE_BUILD_TYPE=Release ..`
5. Build the driver: `make -j$(nproc)`
6. Move the driver to the correct directory: `mv xrealAirLinuxDriver ../../../drivers/xreal_ar_driver`
