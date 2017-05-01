package main

import (
	"fmt"
	"os"
	//"runtime/pprof"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/dertseha/jellui"
	"github.com/dertseha/jellui/area"
	"github.com/dertseha/jellui/area/events"
	"github.com/dertseha/jellui/env/native"
	"github.com/dertseha/jellui/graphics"
	"github.com/dertseha/jellui/input"
	"github.com/dertseha/jellui/opengl"
)

func main() {
	deferrer := make(chan func(), 100)
	defer close(deferrer)

	/*
		f, err := os.Create("profile")
		if err != nil {
			fmt.Println(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	*/

	app := newAreaTestApplication()

	native.Run(app, "BasicAreaExample", 30.0, deferrer)
}

type areaTestApplication struct {
	glWindow jellui.OpenGlWindow
	gl       opengl.OpenGl

	projectionMatrix mgl.Mat4

	mouseX, mouseY float32
	mouseButtons   uint32

	rootArea     *area.Area
	rectRenderer *graphics.RectangleRenderer
}

func newAreaTestApplication() *areaTestApplication {
	return &areaTestApplication{}
}

func (app *areaTestApplication) Init(glWindow jellui.OpenGlWindow) {
	app.setWindow(glWindow)
	app.initOpenGl()
	app.setDebugOpenGl()

	app.rectRenderer = graphics.NewRectangleRenderer(app.gl, &app.projectionMatrix)

	app.initInterface()

	app.onWindowResize(glWindow.Size())
}

func (app *areaTestApplication) setWindow(glWindow jellui.OpenGlWindow) {
	app.glWindow = glWindow
	app.gl = glWindow.OpenGl()

	glWindow.OnRender(app.render)
	glWindow.OnResize(app.onWindowResize)

	glWindow.OnMouseMove(app.onMouseMove)
	glWindow.OnMouseButtonDown(app.onMouseButtonDown)
	glWindow.OnMouseButtonUp(app.onMouseButtonUp)
}

func (app *areaTestApplication) initOpenGl() {
	app.gl.Disable(opengl.DEPTH_TEST)
	app.gl.Enable(opengl.BLEND)
	app.gl.BlendFunc(opengl.SRC_ALPHA, opengl.ONE_MINUS_SRC_ALPHA)
	app.gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func (app *areaTestApplication) setDebugOpenGl() {
	builder := opengl.NewDebugBuilder(app.gl)

	/*
	   builder.OnEntry(func(name string, param ...interface{}) {
	      fmt.Fprintf(os.Stderr, "GL: [%-20s] %v ", name, param)
	   })
	   builder.OnExit(func(name string, result ...interface{}) {
	      fmt.Fprintf(os.Stderr, "-> %v\n", result)
	   })
	*/
	builder.OnError(func(name string, errorCodes []uint32) {
		errorStrings := make([]string, len(errorCodes))
		for index, errorCode := range errorCodes {
			errorStrings[index] = opengl.ErrorString(errorCode)
		}
		fmt.Fprintf(os.Stderr, "!!: [%-20s] %v -> %v\n", name, errorCodes, errorStrings)
	})

	app.gl = builder.Build()
}

func (app *areaTestApplication) initInterface() {
	rootBuilder := area.NewAreaBuilder()

	rootBuilder.SetRight(area.NewAbsoluteAnchor(0.0))
	rootBuilder.SetBottom(area.NewAbsoluteAnchor(0.0))
	rootBuilder.OnRender(func(area *area.Area) {
		app.rectRenderer.Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
			graphics.RGBA(0.125, 0.25, 0.45, 1.0))
	})

	app.rootArea = rootBuilder.Build()

	//
	mainPanelBuilder := area.NewAreaBuilder()
	mainPanelBuilder.SetParent(app.rootArea)

	mainPanelRight := area.NewOffsetAnchor(app.rootArea.Right(), -5.0)
	mainPanelBuilder.SetRight(mainPanelRight)
	mainPanelBuilder.SetLeft(area.NewOffsetAnchor(mainPanelRight, -20.0))
	mainPanelTop := area.NewRelativeAnchor(app.rootArea.Top(), app.rootArea.Bottom(), 0.1)
	mainPanelBuilder.SetTop(mainPanelTop)
	mainPanelBuilder.SetBottom(app.rootArea.Bottom())
	mainPanelBuilder.OnRender(func(area *area.Area) {
		app.rectRenderer.Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
			graphics.RGBA(0.45, 0.06, 0.25, 0.5))
	})
	mainPanelBuilder.Build()

	//
	horizontalCenter := area.NewRelativeAnchor(app.rootArea.Left(), app.rootArea.Right(), 0.5)
	verticalCenter := area.NewRelativeAnchor(app.rootArea.Top(), app.rootArea.Bottom(), 0.5)

	minPanelWidth := float32(50.0)
	minPanelHeight := float32(30.0)

	centerPanelBuilder := area.NewAreaBuilder()
	centerPanelBuilder.SetParent(app.rootArea)

	centerPanelBuilder.SetLeft(area.NewOffsetAnchor(horizontalCenter, minPanelWidth/-2.0))
	centerPanelBuilder.SetRight(area.NewOffsetAnchor(horizontalCenter, minPanelWidth/2.0))
	centerPanelBuilder.SetTop(area.NewOffsetAnchor(verticalCenter, minPanelHeight/-2.0))
	centerPanelBuilder.SetBottom(area.NewOffsetAnchor(verticalCenter, minPanelHeight/2.0))

	centerPanelBuilder.OnRender(func(area *area.Area) {
		app.rectRenderer.Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
			graphics.RGBA(0.25, 0.0, 0.25, 0.75))
	})
	centerPanelBuilder.Build()

	//
	sidePanelBuilder := area.NewAreaBuilder()
	sidePanelBuilder.SetParent(app.rootArea)

	sidePanelLeft := area.NewOffsetAnchor(app.rootArea.Left(), 10.0)
	sidePanelBuilder.SetLeft(sidePanelLeft)
	sidePanelBuilder.SetTop(area.NewOffsetAnchor(app.rootArea.Top(), 10.0))
	sidePanelBuilder.SetBottom(area.NewOffsetAnchor(app.rootArea.Bottom(), -10.0))

	minRight := area.NewOffsetAnchor(sidePanelLeft, 200.0)
	maxRight := area.NewOffsetAnchor(sidePanelLeft, 400.0)
	sidePanelBuilder.SetRight(area.NewLimitedAnchor(minRight, maxRight, area.NewRelativeAnchor(app.rootArea.Left(), app.rootArea.Right(), 0.4)))

	sidePanelBuilder.OnRender(func(area *area.Area) {
		app.rectRenderer.Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
			graphics.RGBA(0.0, 0.33, 0.25, 0.75))
	})
	sidePanelBuilder.Build()

	{
		windowBuilder := area.NewAreaBuilder()
		windowBuilder.SetParent(app.rootArea)

		windowHorizontalCenter := area.NewOffsetAnchor(app.rootArea.Left(), 200.0)
		windowVerticalCenter := area.NewRelativeAnchor(app.rootArea.Top(), app.rootArea.Bottom(), 0.5)

		windowBuilder.SetLeft(area.NewOffsetAnchor(windowHorizontalCenter, -50.0))
		windowBuilder.SetRight(area.NewOffsetAnchor(windowHorizontalCenter, 50.0))
		windowBuilder.SetTop(area.NewOffsetAnchor(windowVerticalCenter, -50.0))
		windowBuilder.SetBottom(area.NewOffsetAnchor(windowVerticalCenter, 50.0))

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
			app.rectRenderer.Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
				graphics.RGBA(1.0, 1.0, 1.0, 0.8))
		})

		windowBuilder.Build()
	}
}

func (app *areaTestApplication) onWindowResize(width int, height int) {
	app.projectionMatrix = mgl.Ortho2D(0.0, float32(width), float32(height), 0.0)
	app.gl.Viewport(0, 0, int32(width), int32(height))

	app.rootArea.Right().RequestValue(float32(width))
	app.rootArea.Bottom().RequestValue(float32(height))
}

func (app *areaTestApplication) render() {
	gl := app.gl

	gl.Clear(opengl.COLOR_BUFFER_BIT)
	app.rootArea.Render()
}

func (app *areaTestApplication) onMouseMove(x float32, y float32) {
	app.mouseX, app.mouseY = x, y
	app.rootArea.DispatchPositionalEvent(events.NewMouseMoveEvent(x, y, 0, 0))
}

func (app *areaTestApplication) onMouseButtonDown(mouseButton uint32, modifier input.Modifier) {
	app.mouseButtons |= mouseButton
	app.rootArea.DispatchPositionalEvent(events.NewMouseButtonEvent(events.MouseButtonDownEventType,
		app.mouseX, app.mouseY, 0, app.mouseButtons, mouseButton))
}

func (app *areaTestApplication) onMouseButtonUp(mouseButton uint32, modifier input.Modifier) {
	app.mouseButtons &= ^mouseButton
	app.rootArea.DispatchPositionalEvent(events.NewMouseButtonEvent(events.MouseButtonUpEventType,
		app.mouseX, app.mouseY, 0, app.mouseButtons, mouseButton))
}
