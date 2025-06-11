import subprocess
import os

from libunreal.supported_devices import supported_devices
from libunreal.edid import UnrealXRDisplayMetadata
import pyedid

def upload_new_device_edid(display_spec: UnrealXRDisplayMetadata, edid: bytes | bytearray):
    pass

def fetch_xr_glass_edid(allow_unsupported_devices) -> UnrealXRDisplayMetadata:
    # Scan for all VGA devices and their IDs
    pci_device_comand = subprocess.run(["lspci"], capture_output=True)

    if pci_device_comand.returncode != 0:
        raise OSError("Failed to scan PCI devices")

    pci_devices: list[str] = pci_device_comand.stdout.decode("utf-8").split("\n")
    pci_devices = pci_devices[:-1]

    vga_devices: list[str] = []

    for pci_device in pci_devices:
        if "VGA compatible controller:" in pci_device:
            vga_devices.append(pci_device[:pci_device.index(" ")])

    # Attempt to find any XR glasses
    for vga_device in vga_devices:
        card_devices = list(os.listdir(f"/sys/devices/pci0000:00/0000:{vga_device}/drm/"))

        for card_device in card_devices:
            if "card" not in card_device:
                continue

            monitors = list(os.listdir(f"/sys/devices/pci0000:00/0000:{vga_device}/drm/{card_device}/"))

            for monitor in monitors:
                if card_device not in monitor:
                    continue

                with open(f"/sys/devices/pci0000:00/0000:{vga_device}/drm/{card_device}/{monitor}/edid", "rb") as edid:
                    raw_edid_file = edid.read()

                    if len(raw_edid_file) == 0:
                        continue

                    edid = pyedid.parse_edid(raw_edid_file)

                    for manufacturer, manufacturer_supported_devices in supported_devices.items():
                        if edid.manufacturer_pnp_id == manufacturer and (edid.name in manufacturer_supported_devices or allow_unsupported_devices):
                            max_width = 0
                            max_height = 0
                            max_refresh = 0

                            for resolution in edid.resolutions:
                                if resolution[0] > max_width and resolution[1] > max_height:
                                    max_width = resolution[0]
                                    max_height = resolution[1]

                                max_refresh = max(max_refresh, int(resolution[2]))

                            if max_width == 0 or max_height == 0:
                                if "max_width" not in manufacturer_supported_devices[edid.name] or "max_height" not in manufacturer_supported_devices[edid.name]:
                                   raise ValueError("Couldn't determine maximum width and height, and the maximum width and height isn't defined in the device quirks section")

                                max_width = int(manufacturer_supported_devices[edid.name]["max_width"])
                                max_height = int(manufacturer_supported_devices[edid.name]["max_height"])

                            if max_refresh == 0:
                                if "max_refresh" not in manufacturer_supported_devices[edid.name]:
                                   raise ValueError("Couldn't determine maximum refresh rate, and the maximum refresh rate isn't defined in the device quirks section")

                                max_refresh = int(manufacturer_supported_devices[edid.name]["max_refresh"])

                            return UnrealXRDisplayMetadata(raw_edid_file, edid.manufacturer_pnp_id, manufacturer_supported_devices[edid.name], max_width, max_height, max_refresh, card_device, monitor.replace(f"{card_device}-", ""))

    raise ValueError("Could not find supported device. Check if the device is plugged in. If it is plugged in and working correctly, check the README or open an issue.")

def upload_edid_firmware(display: UnrealXRDisplayMetadata, fw: bytes | bytearray):
    if display.linux_drm_connector == "" or display.linux_drm_card == "":
        raise ValueError("Linux DRM connector and/or Linux DRM card not specified!")

    with open(f"/sys/kernel/debug/dri/{display.linux_drm_card.replace("card", "")}/{display.linux_drm_connector}/edid_override", "wb") as kernel_edid:
        kernel_edid.write(fw)
