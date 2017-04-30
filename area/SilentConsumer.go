package area

import (
	"github.com/dertseha/jellui/area/events"
)

// SilentConsumer is an event consumer always returning true.
func SilentConsumer(*Area, events.Event) bool {
	return true
}
