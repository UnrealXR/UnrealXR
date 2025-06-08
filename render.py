from loguru import logger
import PyEvdi
import pyray

from libunreal import EvdiDisplaySpec

def render_loop(display_metadata: EvdiDisplaySpec, cards: list[PyEvdi.Card]):
    pyray.init_window(display_metadata.max_width, display_metadata.max_height, "UnrealXR")

    while not pyray.window_should_close():
        # Implement fullscreen toggle
        if pyray.is_key_pressed(pyray.KeyboardKey.KEY_F11):
            display = pyray.get_current_monitor()

            if pyray.is_window_fullscreen():
                pyray.set_window_size(display_metadata.max_width, display_metadata.max_height)
            else:
                pyray.set_window_size(pyray.get_monitor_width(display), pyray.get_monitor_height(display))

            pyray.toggle_fullscreen()
        # Ctrl-C to quit
        elif pyray.is_key_down(pyray.KeyboardKey.KEY_LEFT_CONTROL) and pyray.is_key_down(pyray.KeyboardKey.KEY_C):
            break

        pyray.begin_drawing()
        pyray.clear_background(pyray.BLACK)
        pyray.draw_text("Hello world", 190, 200, 20, pyray.VIOLET)
        pyray.end_drawing()

    logger.info("Goodbye!")
    pyray.close_window()
