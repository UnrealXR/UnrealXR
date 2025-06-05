import dataclasses

@dataclasses.dataclass
class EvdiDisplaySpec:
    edid: bytes
    max_width: int
    max_height: int
    max_refresh_rate: int
