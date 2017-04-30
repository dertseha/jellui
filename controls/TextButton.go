package controls

import (
	"github.com/dertseha/jellui/area"
	"github.com/dertseha/jellui/area/events"
	"github.com/dertseha/jellui/graphics"
	"github.com/dertseha/jellui/input"
)

// ActionHandler is the callback for a firing button (actionable).
type ActionHandler func()

// TextButton is a button with a text label on it.
type TextButton struct {
	area         *area.Area
	rectRenderer *graphics.RectangleRenderer

	label     *Label
	labelLeft area.Anchor
	labelTop  area.Anchor

	actionHandler ActionHandler

	idleColor     graphics.Color
	preparedColor graphics.Color

	prepared bool
	color    graphics.Color
}

// Dispose releases all resources.
func (button *TextButton) Dispose() {
	button.label.Dispose()
	button.area.Remove()
}

func (button *TextButton) onRender(area *area.Area) {
	button.rectRenderer.Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(), button.color)
}

func (button *TextButton) onMouseDown(area *area.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseButtonEvent)

	if mouseEvent.Buttons() == input.MousePrimary {
		area.RequestFocus()
		button.prepare()
		consumed = true
	}

	return
}

func (button *TextButton) onMouseUp(area *area.Area, event events.Event) (consumed bool) {
	mouseEvent := event.(*events.MouseButtonEvent)

	if button.area.HasFocus() && mouseEvent.AffectedButtons() == input.MousePrimary {
		area.ReleaseFocus()
		button.unprepare()
		if button.contains(mouseEvent) {
			button.callHandler()
		}
		consumed = true
	}

	return
}

func (button *TextButton) prepare() {
	if !button.prepared {
		button.color = button.preparedColor
		button.labelLeft.RequestValue(button.labelLeft.Value() + 5)
		button.labelTop.RequestValue(button.labelTop.Value() + 2)
		button.prepared = true
	}
}

func (button *TextButton) unprepare() {
	if button.prepared {
		button.color = button.idleColor
		button.labelLeft.RequestValue(button.labelLeft.Value() - 5)
		button.labelTop.RequestValue(button.labelTop.Value() - 2)
		button.prepared = false
	}
}

func (button *TextButton) contains(event events.PositionalEvent) bool {
	x, y := event.Position()

	return (x >= button.area.Left().Value()) && (x < button.area.Right().Value()) &&
		(y >= button.area.Top().Value()) && (y < button.area.Bottom().Value())
}

func (button *TextButton) callHandler() {
	button.actionHandler()
}
