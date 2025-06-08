from loguru import logger
import PyEvdi
import pyray

from libunreal import EvdiDisplaySpec

def render_loop(display_metadata: EvdiDisplaySpec, cards: list[PyEvdi.Card]):
    while not pyray.window_should_close():
        pyray.begin_drawing()
        pyray.clear_background(pyray.BLACK)
        pyray.draw_text("Hello world from Python!", 190, 200, 96, pyray.WHITE)
        pyray.end_drawing()

    logger.info("Goodbye!")
    pyray.close_window()
