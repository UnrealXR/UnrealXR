from libunreal.supported_devices import *
from libunreal.mcu_driver import *
from libunreal.edid import *

from sys import platform

if platform == "linux" or platform == "linux2":
    from libunreal.linux import *
else:
    raise OSError("Unsupported operating system")
