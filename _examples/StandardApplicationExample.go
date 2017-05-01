package main

import (
	"github.com/dertseha/jellui"
	"github.com/dertseha/jellui/area"
	"github.com/dertseha/jellui/area/events"
	"github.com/dertseha/jellui/env/native"
	"github.com/dertseha/jellui/graphics"
	"github.com/dertseha/jellui/input"
)

func main() {
	deferrer := make(chan func(), 100)
	defer close(deferrer)

	app := jellui.NewStandardApplication(initInterface)

	native.Run(app, "StandardApplicationExample", 30.0, deferrer)
}

func initInterface(app *jellui.StandardApplication, rootArea *area.Area) {
	var window *area.Area

	{
		windowBuilder := area.NewAreaBuilder()
		windowBuilder.SetParent(rootArea)

		windowHorizontalCenter := area.NewRelativeAnchor(rootArea.Left(), rootArea.Right(), 0.5)
		windowVerticalCenter := area.NewRelativeAnchor(rootArea.Top(), rootArea.Bottom(), 0.5)

		windowBuilder.SetLeft(area.NewOffsetAnchor(windowHorizontalCenter, -70.0))
		windowBuilder.SetRight(area.NewOffsetAnchor(windowHorizontalCenter, 70.0))
		windowBuilder.SetTop(area.NewOffsetAnchor(windowVerticalCenter, -8.0))
		windowBuilder.SetBottom(area.NewOffsetAnchor(windowVerticalCenter, 8.0))

		lastGrabX, lastGrabY := float32(0.0), float32(0.0)

		windowBuilder.OnEvent(events.MouseButtonDownEventType, func(area *area.Area, event events.Event) bool {
			buttonEvent := event.(*events.MouseButtonEvent)
			if buttonEvent.Buttons() == input.MousePrimary {
				area.RequestFocus()
				lastGrabX, lastGrabY = buttonEvent.Position()
			}
			return true
		})
		windowBuilder.OnEvent(events.MouseButtonUpEventType, func(area *area.Area, event events.Event) bool {
			buttonEvent := event.(*events.MouseButtonEvent)
			if buttonEvent.AffectedButtons() == input.MousePrimary {
				area.ReleaseFocus()
			}
			return true
		})
		windowBuilder.OnEvent(events.MouseMoveEventType, func(area *area.Area, event events.Event) bool {
			moveEvent := event.(*events.MouseMoveEvent)
			if area.HasFocus() {
				newX, newY := moveEvent.Position()
				windowHorizontalCenter.RequestValue(windowHorizontalCenter.Value() + (newX - lastGrabX))
				windowVerticalCenter.RequestValue(windowVerticalCenter.Value() + (newY - lastGrabY))
				lastGrabX, lastGrabY = newX, newY
			}
			return true
		})
		windowBuilder.OnRender(func(area *area.Area) {
			app.RectangleRenderer().Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
				graphics.RGBA(0.6, 0.2, 0.8, 0.4))
		})

		window = windowBuilder.Build()
	}
	{
		labelBuilder := app.ForLabel()
		labelBuilder.SetParent(window)
		labelBuilder.SetLeft(area.NewOffsetAnchor(window.Left(), 0))
		labelBuilder.SetRight(area.NewOffsetAnchor(window.Right(), 0))
		labelBuilder.SetTop(area.NewOffsetAnchor(window.Top(), 0))
		labelBuilder.SetBottom(area.NewOffsetAnchor(window.Bottom(), 0))
		labelBuilder.Build().SetText("Example Text. Drag it around.")
	}
}
