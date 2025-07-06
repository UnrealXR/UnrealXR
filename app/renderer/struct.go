package renderer

import (
	"git.lunr.sh/UnrealXR/unrealxr/evdi/libevdi"
)

type EvdiDisplayMetadata struct {
	EvdiNode     *libevdi.EvdiNode
	Rect         *libevdi.EvdiDisplayRect
	Buffer       *libevdi.EvdiBuffer
	EventContext *libevdi.EvdiEventContext
}
