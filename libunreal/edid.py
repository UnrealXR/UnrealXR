from loguru import logger

import dataclasses
import uuid

@dataclasses.dataclass
class EvdiDisplaySpec:
    edid: bytes
    max_width: int
    max_height: int
    max_refresh_rate: int
    linux_drm_card: str
    linux_drm_connector: str

calculated_msft_payload_size = 22+4

def calculate_checksum(block):
    checksum = (-sum(block[:-1])) & 0xFF
    return checksum

def patch_edid_to_be_specialized(edid_data: bytes | bytearray) -> bytes | bytearray:
    mutable_edid = bytearray(edid_data)
    is_enhanced_mode = len(edid_data) > 128

    found_extension_base = 0
    extension_base_existed = False

    if is_enhanced_mode:
        for i in range(128, len(edid_data), 128):
            if edid_data[i] == 0x02:
                logger.warning("Detected existing ANSI CTA data section. Patching in place but untested! Please report any issues you discover")

                if edid_data[i+1] != 0x03:
                    logger.warning("Incompatible version detected for ANSI CTA data section in EDID")

                found_extension_base = i
                extension_base_existed = True

        if found_extension_base == 0:
            found_extension_base = len(edid_data)
            mutable_edid.extend([0]*128)
    else:
        mutable_edid.extend([0]*128)
        found_extension_base = 128

    generated_uuid = uuid.uuid4()

    mutable_edid[found_extension_base] = 0x02
    mutable_edid[found_extension_base+1] = 0x03

    if extension_base_existed and mutable_edid[found_extension_base+2] != calculated_msft_payload_size and mutable_edid[found_extension_base+2] != 0:
        # We try our best to move our data into place
        current_base = mutable_edid[found_extension_base+2]
        mutable_edid[found_extension_base+2] = calculated_msft_payload_size+1

        mutable_edid[found_extension_base+4:found_extension_base+current_base-1] = [0]*(current_base-1)
        mutable_edid[found_extension_base+calculated_msft_payload_size:found_extension_base+127] = mutable_edid[found_extension_base+current_base:found_extension_base+127]
    else:
        mutable_edid[found_extension_base+2] = calculated_msft_payload_size

    if not extension_base_existed:
        mutable_edid[126] += 1
        mutable_edid[127] = calculate_checksum(mutable_edid[:128])

        mutable_edid[found_extension_base+3] = 0 # We don't know any of these properties

    # Implemented using https://learn.microsoft.com/en-us/windows-hardware/drivers/display/specialized-monitors-edid-extension
    # VST & Length
    mutable_edid[found_extension_base+4] = 0x3 << 5 | 0x15 # 0x3: vendor specific tag; 0x15: length
    # Assigned IEEE OUI
    mutable_edid[found_extension_base+5] = 0x5C
    mutable_edid[found_extension_base+6] = 0x12
    mutable_edid[found_extension_base+7] = 0xCA
    # Actual data
    mutable_edid[found_extension_base+8] = 0x2 # Using version 0x2 for better compatibility
    mutable_edid[found_extension_base+9] = 0x7 # Using VR tag for better compatibility even though it probably doesn't matter
    mutable_edid[found_extension_base+10:found_extension_base+10+16] = generated_uuid.bytes

    mutable_edid[found_extension_base+127] = calculate_checksum(mutable_edid[found_extension_base:found_extension_base+127])

    if isinstance(edid_data, bytes):
        return bytes(mutable_edid)
    else:
        return mutable_edid
