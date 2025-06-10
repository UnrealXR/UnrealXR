from time import sleep
import math

from loguru import logger
import PyEvdi
import pyray

from libunreal import EvdiDisplaySpec, MCUCallbackWrapper, start_mcu_event_listener

previous_pitch = 0.0
previous_yaw = 0.0
previous_roll = 0.0

current_pitch = 0.0
current_yaw = 0.0
current_roll = 0.0

has_gotten_pitch_callback_before = False
has_gotten_yaw_callback_before = False
has_gotten_roll_callback_before = False

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

def render_loop(display_metadata: EvdiDisplaySpec, cards: list[PyEvdi.Card]):
    logger.info("Starting sensor event listener")

    mcu_callbacks = MCUCallbackWrapper(roll_callback, pitch_callback, yaw_callback, text_message, stub_brightness_function, stub_brightness_function)
    start_mcu_event_listener(display_metadata.device_vendor, mcu_callbacks)

    logger.info("Beginning sensor initialization. Awaiting first sensor update")

    while (not has_gotten_pitch_callback_before) or (not has_gotten_yaw_callback_before) or (not has_gotten_roll_callback_before):
        sleep(0.01)

    logger.info("Initialized sensors")

    camera = pyray.Camera3D()
    camera.position = pyray.Vector3(10.0, 10.0, 10.0)
    camera.target = pyray.Vector3(0.0, 0.0, 0.0)
    camera.up = pyray.Vector3(0.0, 1.0, 0.0)
    camera.fovy = 45.0
    camera.projection = pyray.CameraProjection.CAMERA_PERSPECTIVE

    cube_position = pyray.Vector3(0.0, 0.0, 0.0)

    movement_vector = pyray.Vector3()
    look_vector = pyray.Vector3()

    logger.error("QUIRK: Waiting 10 seconds before reading sensors due to sensor drift bugs")
    sleep(10)
    logger.error("Continuing...")

    while not pyray.window_should_close():
        look_vector.x = (current_yaw-previous_yaw)*6.5
        look_vector.y = (current_pitch-previous_pitch)*6.5
        # the Z vector is more trouble than its worth so it just doesn't get accounted for...
        #look_vector.z = (current_roll-previous_roll)*6.5

        pyray.update_camera_pro(camera, movement_vector, look_vector, 0.0)

        pyray.begin_drawing()
        pyray.clear_background(pyray.BLACK)
        pyray.begin_mode_3d(camera)
        pyray.draw_cube(cube_position, 2.0, 2.0, 2.0, pyray.ORANGE)
        pyray.end_mode_3d()
        pyray.end_drawing()

    logger.info("Goodbye!")
    pyray.close_window()
