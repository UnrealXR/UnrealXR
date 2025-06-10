from tempfile import TemporaryDirectory
from dataclasses import dataclass
from os import path, environ
from typing import Callable
from loguru import logger
from shutil import which
from time import sleep
from enum import Enum
import subprocess
import threading
import socket
import struct
import signal
import atexit
import os

class MCUCommandTypes(Enum):
    ROLL = 0
    PITCH = 1
    YAW = 2
    GENERIC_MESSAGE = 3
    BRIGHTNESS_UP = 4
    BRIGHTNESS_DOWN = 5

@dataclass
class MCUCallbackWrapper:
    OnRollUpdate: Callable[[float], None]
    OnPitchUpdate: Callable[[float], None]
    OnYawUpdate: Callable[[float], None]
    OnTextMessageRecieved: Callable[[str], None]
    OnBrightnessUp: Callable[[int], None]
    OnBrightnessDown: Callable[[int], None]

vendor_to_driver_table: dict[str, str] = {
    "MRG": "xreal_ar_driver",
}

def find_executable_path_from_driver_name(driver_name) -> str:
    # First try the normal driver path
    try:
        driver_path = path.join("drivers", driver_name)

        file = open(driver_path)
        file.close()

        return driver_path
    except OSError:
        # Then search the system path
        driver_path = which(driver_name)

        if driver_path == "":
            raise OSError("Could not find driver executable in driver directory or in PATH")

        return driver_path

def start_mcu_event_listener(driver_vendor: str, events: MCUCallbackWrapper):
    driver_executable = find_executable_path_from_driver_name(vendor_to_driver_table[driver_vendor])

    created_temp_dir = TemporaryDirectory()
    sock_path = path.join(created_temp_dir.name, "mcu_socket")

    def on_socket_event(sock: socket.socket):
        while True:
            message_type = sock.recv(1)

            if message_type[0] == MCUCommandTypes.ROLL.value:
                roll_data = sock.recv(4)
                roll_value = struct.unpack("!f", roll_data)[0]

                if not isinstance(roll_value, float):
                    logger.warning("Expected roll value to be a float but got other type instead")
                    continue

                events.OnRollUpdate(roll_value)
            elif message_type[0] == MCUCommandTypes.PITCH.value:
                pitch_data = sock.recv(4)
                pitch_value = struct.unpack("!f", pitch_data)[0]

                if not isinstance(pitch_value, float):
                    logger.warning("Expected pitch value to be a float but got other type instead")
                    continue

                events.OnPitchUpdate(pitch_value)
            elif message_type[0] == MCUCommandTypes.YAW.value:
                yaw_data = sock.recv(4)
                yaw_value = struct.unpack("!f", yaw_data)[0]

                if not isinstance(yaw_value, float):
                    logger.warning("Expected yaw value to be a float but got other type instead")
                    continue

                events.OnYawUpdate(yaw_value)
            elif message_type[0] == MCUCommandTypes.GENERIC_MESSAGE.value:
                length_bytes = sock.recv(4)

                msg_len = struct.unpack("!I", length_bytes)[0]
                msg_bytes = sock.recv(msg_len)

                msg = msg_bytes.decode("utf-8", errors="replace")
                events.OnTextMessageRecieved(msg)
            elif message_type[0] == MCUCommandTypes.BRIGHTNESS_UP.value:
                brightness_bytes = sock.recv(1)
                events.OnBrightnessUp(int.from_bytes(brightness_bytes, byteorder='big'))
            elif message_type[0] == MCUCommandTypes.BRIGHTNESS_DOWN.value:
                brightness_bytes = sock.recv(1)
                events.OnBrightnessDown(int.from_bytes(brightness_bytes, byteorder='big'))
            else:
                logger.warning(f"Unknown message type recieved: {str(message_type[0])}")

    def start_socket_handout():
        server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        server.bind(sock_path)
        server.listen(1)

        sock: socket.socket | None = None

        while True:
            sock, _ = server.accept()

            threaded_connection_processing = threading.Thread(target=on_socket_event, args=(sock,), daemon=True)
            threaded_connection_processing.start()

        created_temp_dir.cleanup()

    def spawn_child_process():
        custom_env = environ.copy()
        custom_env["UNREALXR_NREAL_DRIVER_SOCK"] = sock_path

        process = subprocess.Popen([driver_executable], env=custom_env)

        def kill_child():
            if process.pid is None:
                pass
            else:
                os.kill(process.pid, signal.SIGTERM)

        atexit.register(kill_child)

    threaded_socket_handout = threading.Thread(target=start_socket_handout, daemon=True)
    threaded_socket_handout.start()

    sleep(0.01) # Give the socket server time to initialize

    threaded_child_process = threading.Thread(target=spawn_child_process, daemon=True)
    threaded_child_process.start()
