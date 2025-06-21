from dataclasses import dataclass
from io import BufferedWriter
from sys import int_info
from typing import Union
import ctypes
import time
import math

from loguru import logger
from raylib import rl
import PyEvdi
import pyray

from libunreal import UnrealXRDisplayMetadata, MCUCallbackWrapper, start_mcu_event_listener

vertical_size = 0.0
horizontal_sizing_constant = 1

previous_pitch = 0.0
previous_yaw = 0.0
previous_roll = 0.0

current_pitch = 0.0
current_yaw = 0.0
current_roll = 0.0

has_gotten_pitch_callback_before = False
has_gotten_yaw_callback_before = False
has_gotten_roll_callback_before = False

@dataclass
class RectMetadata:
    card: PyEvdi.Card
    buffer_ptr: pyray.ffi.CData | None
    texture: Union[pyray.Texture, None]
    model: Union[pyray.Model, None]
    angle: int
    relative_position: int

def pitch_callback(new_pitch: float):
    global current_pitch
    global previous_pitch
    global has_gotten_pitch_callback_before

    if not has_gotten_pitch_callback_before:
        has_gotten_pitch_callback_before = True
        previous_pitch = new_pitch
        current_pitch = new_pitch
    else:
        previous_pitch = current_pitch
        current_pitch = new_pitch

def yaw_callback(new_yaw: float):
    global current_yaw
    global previous_yaw
    global has_gotten_yaw_callback_before

    if not has_gotten_yaw_callback_before:
        has_gotten_yaw_callback_before = True
        previous_yaw = new_yaw
        current_yaw = new_yaw
    else:
        previous_yaw = current_yaw
        current_yaw = new_yaw

def roll_callback(new_roll: float):
    global current_roll
    global previous_roll
    global has_gotten_roll_callback_before

    if not has_gotten_roll_callback_before:
        has_gotten_roll_callback_before = True
        previous_roll = new_roll
        roll = new_roll
    else:
        previous_roll = current_roll
        current_roll = new_roll

def text_message(message: str):
    logger.debug(f"Got message from AR's MCU: {message}")

def stub_brightness_function(brightness: int):
    pass

def find_max_vertical_size(fovy_deg: float, distance: float) -> float:
    fovy_rad = math.radians(fovy_deg)
    return 2 * distance * math.tan(fovy_rad / 2)

def find_optimal_horizonal_res(vertical_display_res: int, horizontal_display_res: int) -> float:
    aspect_ratio = horizontal_display_res/vertical_display_res
    horizontal_size = vertical_size * aspect_ratio
    horizontal_size = horizontal_size * horizontal_sizing_constant

    return horizontal_size

def render_loop(display_metadata: UnrealXRDisplayMetadata, config: dict[str, str | int], cards: list[PyEvdi.Card]):
    global vertical_size
    global core_mesh
    logger.info("Starting sensor event listener")

    mcu_callbacks = MCUCallbackWrapper(roll_callback, pitch_callback, yaw_callback, text_message, stub_brightness_function, stub_brightness_function)
    start_mcu_event_listener(display_metadata.device_vendor, mcu_callbacks)

    logger.info("Beginning sensor initialization. Awaiting first sensor update")

    while (not has_gotten_pitch_callback_before) or (not has_gotten_yaw_callback_before) or (not has_gotten_roll_callback_before):
        time.sleep(0.01)

    logger.info("Initialized sensors")

    camera = pyray.Camera3D()
    camera.fovy = 45.0

    vertical_size = find_max_vertical_size(camera.fovy, 5.0)

    camera.position = pyray.Vector3(0.0, vertical_size/2, 5.0)
    camera.target = pyray.Vector3(0.0, vertical_size/2, 0.0)
    camera.up = pyray.Vector3(0.0, 1.0, 0.0)
    camera.projection = pyray.CameraProjection.CAMERA_PERSPECTIVE

    core_mesh = pyray.gen_mesh_plane(find_optimal_horizonal_res(display_metadata.max_height, display_metadata.max_width), vertical_size, 1, 1)

    movement_vector = pyray.Vector3()
    look_vector = pyray.Vector3()

    has_z_vector_disabled_quirk = False
    has_sensor_init_delay_quirk = False
    sensor_init_start_time = time.time()

    if "z_vector_disabled" in display_metadata.device_quirks and bool(display_metadata.device_quirks["z_vector_disabled"]):
        logger.warning("QUIRK: The Z vector has been disabled for your specific device")
        has_z_vector_disabled_quirk = True

    if "sensor_init_delay" in display_metadata.device_quirks:
        logger.warning(f"QUIRK: Waiting {str(display_metadata.device_quirks["sensor_init_delay"])} second(s) before reading sensors")
        logger.warning("|| MOVEMENT WILL NOT BE OPERATIONAL DURING THIS TIME. ||")
        sensor_init_start_time = time.time()
        has_sensor_init_delay_quirk = True

    rects: list[RectMetadata] = []

    if int(config["display_count"]) >= 2:
        display_angle = int(config["display_angle"])
        display_spacing = int(config["display_spacing"])
        total_displays = int(config["display_count"])

        highest_possible_angle_on_both_sides = (total_displays-1)*display_angle
        highest_possible_pixel_spacing_on_both_sides = (total_displays-1)*display_spacing

        for i in range(total_displays):
            current_angle = (-highest_possible_angle_on_both_sides)+(display_angle*i)
            current_display_spacing = (-highest_possible_pixel_spacing_on_both_sides)+(display_spacing*i)

            rect_metadata = RectMetadata(cards[i], None, None, None, current_angle, current_display_spacing)

            has_acquired_fb = False

            def fb_acquire_handler(evdi_buffer: PyEvdi.Buffer):
                nonlocal has_acquired_fb

                if has_acquired_fb:
                    return

                has_acquired_fb = True
                logger.info(f"Acquired buffer for card #{i+1} with ID {evdi_buffer.id}")

                address = ctypes.pythonapi.PyCapsule_GetPointer
                address.restype = ctypes.c_void_p
                address.argtypes = [ctypes.py_object, ctypes.c_char_p]

                buffer_void_ptr = address(evdi_buffer.bytes, None)
                rect_metadata.buffer_ptr = pyray.ffi.cast("void *", buffer_void_ptr)

                pyray_image = pyray.Image()

                pyray_image.data = rect_metadata.buffer_ptr
                pyray_image.width = display_metadata.max_width
                pyray_image.height = display_metadata.max_height
                pyray_image.mipmaps = 1
                pyray_image.format = pyray.PixelFormat.PIXELFORMAT_UNCOMPRESSED_R8G8B8A8

                rect_metadata.texture = pyray.load_texture_from_image(pyray_image)
                rect_metadata.model = pyray.load_model_from_mesh(core_mesh)

                pyray.set_material_texture(rect_metadata.model.materials[0], pyray.MaterialMapIndex.MATERIAL_MAP_ALBEDO, rect_metadata.texture)

            cards[i].acquire_framebuffer_handler = fb_acquire_handler
            cards[i].handle_events(1000)

            rects.append(rect_metadata)

    while not pyray.window_should_close():
        if has_sensor_init_delay_quirk:
            if time.time() - sensor_init_start_time >= int(display_metadata.device_quirks["sensor_init_delay"]):
                # Unset the quirk state
                logger.info("Movement is now enabled.")
                has_sensor_init_delay_quirk = False
        else:
            look_vector.x = (current_yaw-previous_yaw)*6.5
            look_vector.y = -(current_pitch-previous_pitch)*6.5

            if not has_z_vector_disabled_quirk:
                look_vector.z = (current_roll-previous_roll)*6.5

            pyray.update_camera_pro(camera, movement_vector, look_vector, 0.0)

        pyray.begin_drawing()
        pyray.clear_background(pyray.BLACK)
        pyray.begin_mode_3d(camera)

        for rect_count in range(len(rects)):
            rect = rects[rect_count]

            if rect.buffer_ptr is None or rect.texture is None or rect.model is None:
                continue

            cards[rect_count].handle_events(1)
            pyray.update_texture(rect.texture, rect.buffer_ptr)

            pyray.draw_model_ex(
                rect.model,
                pyray.Vector3(0, vertical_size/2, 0),
                pyray.Vector3(1, 0, 0), # rotate around X to make it vertical
                90,
                pyray.Vector3(1, 1, 1),
                pyray.WHITE
            )

            break

        pyray.end_mode_3d()
        pyray.end_drawing()

    logger.info("Goodbye!")
    pyray.close_window()
