module git.terah.dev/UnrealXR/unrealxr/app

go 1.24.3

replace git.terah.dev/UnrealXR/unrealxr/ardriver => ../ardriver

replace git.terah.dev/UnrealXR/unrealxr/edidpatcher => ../edidpatcher

require (
	git.terah.dev/UnrealXR/raylib-go/raylib v0.55.2-0.20250623002739-1468af2636e1
	git.terah.dev/UnrealXR/unrealxr/ardriver v0.0.0-00010101000000-000000000000
	git.terah.dev/UnrealXR/unrealxr/edidpatcher v0.0.0-00010101000000-000000000000
	git.terah.dev/imterah/goevdi/libevdi v0.1.0-evdi1.14.10
	github.com/anoopengineer/edidparser v0.0.0-20240602223913-86ca9ed3d2b0
	github.com/charmbracelet/log v0.4.2
	github.com/goccy/go-yaml v1.18.0
	github.com/kirsle/configdir v0.0.0-20170128060238-e45d2f54772f
	github.com/tebeka/atexit v0.3.0
	github.com/urfave/cli/v3 v3.3.8
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/colorprofile v0.2.3-0.20250311203215-f60798e515dc // indirect
	github.com/charmbracelet/lipgloss v1.1.0 // indirect
	github.com/charmbracelet/x/ansi v0.8.0 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.13-0.20250311204145-2c3ea96c31dd // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/muesli/termenv v0.16.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect
	golang.org/x/sys v0.33.0 // indirect
)
