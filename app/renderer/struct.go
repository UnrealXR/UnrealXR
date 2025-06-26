package renderer

import (
	"git.terah.dev/imterah/goevdi/libevdi"
)

type EvdiDisplayMetadata struct {
	EvdiNode            *libevdi.EvdiNode
	Rect                *libevdi.EvdiDisplayRect
	Buffer              *libevdi.EvdiBuffer
	EventContext        *libevdi.EvdiEventContext
}
