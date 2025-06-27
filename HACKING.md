# UnrealXR Development Documentation

## Dummies and You

The dummy/development build generally isn't recommended for use if you're developing headset or OS drivers. It's intended for development on the 3D environment itself only (ie. virtual display positioning, widgets, etc.)

To use it, all you need to do is run the Makefile with the `dev` target:

```bash
make dev
```

In the code, this mode is typically referred to as "dummy mode" (`dummy_ar.go`) or "fake mode" (ie. `patching_tools_fake_patching.go`).

To control the X11 window, you can use the arrow keys to move the camera around. You **cannot** use the mouse due to the nature of this application. If you use the mouse for this, it would make development really difficult as you're using virtual displays and you need your mouse to interact with the desktop.

If you're adding new headset drivers, but you still want to test in a window, you still can. It's not the dummy mode, but you can compile UnrealXR without the DRM flags (`drm drm_leasing drm_disable_input`) manually. This may differ in the future, but for now, the command is:

```bash
cd app; go build -v -tags "headset_driver_goes_here noaudio" -o ../uxr .; cd ..
```

Replace `headset_driver_goes_here` with the name of your headset driver. For example, `xreal` is the Xreal driver.
