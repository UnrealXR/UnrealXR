# UnrealXR

UnrealXR is a spatial multi-display renderer for the Xreal line of devices, enabling immersive, simultaneous viewing of multiple desktops and applications in 3D space.

## Why not Breezy Desktop?

1. **Limited desktop support**: Breezy only supports GNOME Wayland, whereas UnrealXR supports all desktops as long as they are Wayland and implement the [drm-lease-v1](https://wayland.app/protocols/drm-lease-v1) protocol. All major desktops aside from [Sway](https://swaywm.org/) and [Gamescope](https://wiki.archlinux.org/title/Gamescope) support this. [UnrealXR may add X11 support in the future](https://git.terah.dev/imterah/unrealxr/issues/1).
2. **Pricing**: UnrealXR will always be free and open source. The catch is that since Breezy costs money, it will likely have more of an "out-of-the-box" working experience, and additionally have more polish.

*However*, Breezy Desktop supports more AR headsets. We only support Xreal devices as that's all I have.

**TL;DR**: Breezy Desktop costs money and it has limited desktop support. However, due to its cost, it'll likely be more stable in the long term as it has funding. Additionally, Breezy Desktop supports more AR devices, while we only support Xreal devices.

## Why not the normal Nebula app?

1. **Incompatibility**: The normal Nebula app doesn't support Linux. Of course, we won't support Windows for the forseeable future, and macOS isn't on our roadmap. So this isn't really a fair point.
2. **Closed Source**: The Nebula app and the normal SDK is closed source, so we don't really know what it's doing internally.

**TL;DR**: The normal Nebula app is only compatible with Windows and macOS. However, we only support Linux, so that isn't a fair comparison. The Nebula app is also closed source.

## Dependencies

Before anything, this depends on the `evdi` Linux kernel module. This is packaged in Debian-based distributions as `evdi-dkms`. If you already have DisplayLink drivers installed for their devices, you likely do not need to do this step. After installing this, please reboot your computer.

You'll need to install build dependencies after this. For Debian-based distros, the dependency list should be: `git golang build-essential libdrm libdrm-dev linux-headers-$(uname -r) cmake clang-tools pkg-config libwayland-client++1 libgl1-mesa-dev libglu1-mesa-dev libwayland-dev libxkbcommon-dev libhidapi-dev libjson-c-dev	libudev-dev libusb-1.0-0 libusb-1.0-0-dev libopencv-dev`

If you're using Nix/NixOS, all you need to do is use `nix-shell` to enter the development environment.

## Building

Just run `make` in the root directory.
