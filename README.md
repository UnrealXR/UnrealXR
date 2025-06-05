# UnrealXR

UnrealXR is a display multiplexer for the original Nreal Air (other devices may be supported).

## Development

1. Clone this repository with submodules: `git clone --recurse-submodules https://git.terah.dev/imterah/unrealxr.git`
2. Initialize the development shell: `nix-shell`
3. Create a virtual environment: `python3 -m venv .venv; source .venv/bin/activate`
4. Install Python dependencies: `pip install -r requirements.txt`
4. Build libevdi: `cd evdi/library; make -j$(nproc); cd ../..`
5. Build pyevdi: `cd evdi/pyevdi; make -j$(nproc); make install; cd ../..`
