# Sourced from "https://uefi.org/uefi-pnp-export"
supported_devices: dict[str, dict[str, dict[str, str | int]]] = {
    "MRG": {
        # Quirks section
        "Air": {
            "max_width": 1920,
            "max_height": 1080,
            "max_refresh": 120,
            "sensor_init_delay": 10,
            "z_vector_disabled": True,
        }
    }
}
