#!/usr/bin/env python3
from sys import platform
import atexit
import json
import os

from platformdirs import user_data_dir, user_config_dir
from loguru import logger
import PyEvdi
import pyedid
import pyray
import time

from render import render_loop
import libunreal

default_configuration: dict[str, str | int] = {
    "display_angle": 45,
    "display_count": 3,
    "allow_unsupported_devices": False,
    "allow_unsupported_vendors": False,
    "override_default_edid": False,
    "edid_override_path": "/file/here",
    "override_width": 0,
    "override_height": 0,
    "override_refresh_rate": 0,
}

def initialize_configuration():
    pass

used_cards: list[int] = []

def find_suitable_evdi_card() -> int:
    for i in range(20):
        if PyEvdi.check_device(i) == PyEvdi.AVAILABLE and i not in used_cards:
            used_cards.append(i)
            return i

    PyEvdi.add_device()

    for i in range(20):
        if PyEvdi.check_device(i) == PyEvdi.AVAILABLE and i not in used_cards:
            used_cards.append(i)
            return i

    raise ValueError("Failed to allocate virtual display device")

@logger.catch
def main():
    configuration = {}

    logger.info("Loading configuration")
    config_dir = os.environ["UNREALXR_CONFIG_PATH"] if "UNREALXR_CONFIG_PATH" in os.environ else ""
    data_dir = os.environ["UNREALXR_DATA_PATH"] if "UNREALXR_DATA_PATH" in os.environ else ""

    # Use OS defaults if we weren't overriden in env
    if config_dir == "":
        config_dir = user_config_dir("UnrealXR", "Tera")

    if data_dir == "":
        data_dir = user_data_dir("UnrealXR", "Tera")

    try:
        os.stat(data_dir)
    except OSError:
        os.makedirs(data_dir)

    # Read config and create it if it doesn't exist
    config_path = os.path.join(config_dir, "config.json")

    try:
        os.stat(config_path)

        with open(config_path, "r") as config_file:
            configuration = json.load(config_file)
    except OSError:
        try:
            os.makedirs(config_dir)
        except OSError:
            pass

        with open(config_path, "w") as config_file:
            json.dump(default_configuration, config_file, indent=4)

        configuration = default_configuration

    # Set unbound values (ie. if user is using an older version)
    for key, default_value in default_configuration.items():
        if key not in configuration:
            # Add quotes if we're a string
            value = default_value

            if isinstance(value, str):
                value = f'"{value}"'

            logger.warning(f"Setting unbound key '{key}' with default value {value}. You might want to define this! If not, this warning can be safely ignored.")
            configuration[key] = default_value

    # Initialize logging to files
    logger.add(os.path.join(data_dir, "unrealxr.log"), format="{time:YYYY-MM-DD at HH:mm:ss} | {level} | {message}")
    logger.info("Loaded configuration")

    if os.geteuid() != 0:
        raise OSError("You are not running as root! Running as root is necessary to talk to the EVDI service")

    # Get the display EDID
    logger.info("Attempting to read display EDID file")
    edid: libunreal.EvdiDisplaySpec | None = None

    if configuration["override_default_edid"] or configuration["allow_unsupported_vendors"]:
        # We need to parse it to get the maximum width, height, and refresh rate for EVDI's calculations
        with open(configuration["edid_override_path"], "rb") as edid_file:
            edid_file = edid_file.read()
            parsed_edid_file = pyedid.parse_edid(edid_file)

            max_width = int(configuration["override_width"])
            max_height = int(configuration["override_height"])
            max_refresh = int(configuration["override_refresh_rate"])

            if configuration["override_width"] == 0 or configuration["override_height"] == 0 or configuration["override_refresh_rate"] == 0:
                for resolution in parsed_edid_file.resolutions:
                    if configuration["override_width"] == 0 or configuration["override_height"] == 0:
                        if resolution[0] > max_width and resolution[1] > max_height:
                            max_width = resolution[0]
                            max_height = resolution[1]

                    if configuration["override_refresh_rate"] == 0:
                        max_refresh = max(max_refresh, int(resolution[2]))

            if max_width == 0 or max_height == 0:
                raise ValueError("Could not determine maximum width and/or height from EDID file, and the width and/or height overrides aren't set!")

            if max_refresh == 0:
                raise ValueError("Could not determine maximum refresh rate from EDID file, and the refresh rate overrides aren't set!")

            edid = libunreal.EvdiDisplaySpec(edid_file, max_width, max_height, max_refresh, "", "")
    else:
        edid = libunreal.fetch_xr_glass_edid(configuration["allow_unsupported_devices"])

    assert(edid is not None)
    logger.info("Got EDID file")

    if platform == "linux" or platform == "linux2":
        # TODO: implement EDID patching for overridden displays
        logger.info("Patching EDID firmware")
        patched_edid = libunreal.patch_edid_to_be_specialized(edid.edid)
        libunreal.upload_edid_firmware(edid, patched_edid)

        def unload_custom_fw():
            with open(f"/sys/kernel/debug/dri/{edid.linux_drm_card.replace("card", "")}/{edid.linux_drm_connector}/edid_override", "w") as kernel_edid:
                kernel_edid.write("reset")

            logger.info("Please unplug and plug in your XR device to restore it back to normal settings.")

        atexit.register(unload_custom_fw)
        input("Press the Enter key to continue loading after you unplug and plug in your XR device.")

    # Raylib gets confused if there's multiple dri devices so we initialize the window before anything
    logger.info("Initializing XR headset")
    pyray.set_target_fps(edid.max_refresh_rate)
    pyray.init_window(edid.max_width, edid.max_height, "UnrealXR")

    logger.info("Initializing virtual displays")
    cards = []

    for i in range(int(configuration["display_count"])):
        suitable_card_id = find_suitable_evdi_card()
        card = PyEvdi.Card(suitable_card_id)
        card.connect(edid.edid, len(edid.edid), edid.max_width*edid.max_height, edid.max_width*edid.max_height*edid.max_refresh_rate)
        cards.append(card)

        logger.debug(f"Initialized card #{str(i+1)}")

        atexit.register(lambda: card.close())

    logger.info("Initialized displays. Entering rendering loop")
    render_loop(edid, cards)
if __name__ == "__main__":
    print("Welcome to UnrealXR!\n")
    main()
