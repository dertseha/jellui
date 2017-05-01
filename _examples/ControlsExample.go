package main

import (
	"fmt"
	"os"
	//"runtime/pprof"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/dertseha/jellui"
	"github.com/dertseha/jellui/area"
	"github.com/dertseha/jellui/area/events"
	"github.com/dertseha/jellui/controls"
	"github.com/dertseha/jellui/env/native"
	"github.com/dertseha/jellui/font"
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

	app := newControlsTestApplication()

	native.Run(app, "ControlsExample", 30.0, deferrer)
}

type controlsTestApplication struct {
	glWindow jellui.OpenGlWindow
	gl       opengl.OpenGl

	projectionMatrix mgl.Mat4

	mouseX, mouseY float32
	mouseButtons   uint32

	uiFontPainter    graphics.TextPainter
	largeFontPainter graphics.TextPainter
	uiTextPalette    *graphics.PaletteTexture
	rectRenderer     *graphics.RectangleRenderer
	uiTextRenderer   *graphics.BitmapTextureRenderer

	rootArea *area.Area
}

func newControlsTestApplication() *controlsTestApplication {
	return &controlsTestApplication{}
}

func (app *controlsTestApplication) Init(glWindow jellui.OpenGlWindow) {
	app.setWindow(glWindow)
	app.initOpenGl()
	app.setDebugOpenGl()

	app.initGraphics()
	app.initInterface()

	app.onWindowResize(glWindow.Size())
}

func (app *controlsTestApplication) setWindow(glWindow jellui.OpenGlWindow) {
	app.glWindow = glWindow
	app.gl = glWindow.OpenGl()

	glWindow.OnRender(app.render)
	glWindow.OnResize(app.onWindowResize)

	glWindow.OnMouseMove(app.onMouseMove)
	glWindow.OnMouseButtonDown(app.onMouseButtonDown)
	glWindow.OnMouseButtonUp(app.onMouseButtonUp)
	glWindow.OnMouseScroll(app.onMouseScroll)
}

func (app *controlsTestApplication) initOpenGl() {
	app.gl.Disable(opengl.DEPTH_TEST)
	app.gl.Enable(opengl.BLEND)
	app.gl.BlendFunc(opengl.SRC_ALPHA, opengl.ONE_MINUS_SRC_ALPHA)
	app.gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func (app *controlsTestApplication) setDebugOpenGl() {
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

func (app *controlsTestApplication) initGraphics() {
	uiTextPalette := map[int][4]byte{
		0: {0x00, 0x00, 0x00, 0x00},
		1: {0x80, 0x94, 0x54, 0xFF},
		2: {0x00, 0x00, 0x00, 0xC0},

		90: {0x80, 0x54, 0x94, 0xFF},
		92: {0x70, 0x44, 0x84, 0xFF},
		94: {0x60, 0x34, 0x74, 0xFF},
		95: {0x50, 0x24, 0x64, 0xFF},
		98: {0x40, 0x14, 0x54, 0x20},
		/*
			90: {0x80, 0x94, 0x54, 0xFF},
			92: {0x70, 0x84, 0x44, 0xFF},
			94: {0x60, 0x74, 0x34, 0xFF},
			95: {0x50, 0x64, 0x24, 0xFF},
			98: {0x40, 0x54, 0x14, 0x20},
		*/
		/*
			90: {0x80, 0x94, 0x54, 0xFF},
			92: {0x7C, 0x90, 0x50, 0xFF},
			94: {0x78, 0x8C, 0x4C, 0xFF},
			95: {0x74, 0x88, 0x48, 0xFF},
			98: {0x70, 0x84, 0x44, 0x20},
		*/
	}
	app.uiTextPalette = app.NewPaletteTexture(func(index int) (byte, byte, byte, byte) {
		entry := uiTextPalette[index]
		return entry[0], entry[1], entry[2], entry[3]
	})
	viewMatrix := mgl.Ident4()
	uiRenderContext := graphics.NewBasicRenderContext(app.gl, &app.projectionMatrix, &viewMatrix)
	app.uiTextRenderer = graphics.NewBitmapTextureRenderer(uiRenderContext, app.uiTextPalette)

	app.uiFontPainter = graphics.NewBitmapTextPainter(font.SmallShock, 0x02)
	app.largeFontPainter = graphics.NewBitmapTextPainter(font.ColorHeadingShock, 0x00)

	app.rectRenderer = graphics.NewRectangleRenderer(app.gl, &app.projectionMatrix)
}

func (app *controlsTestApplication) initInterface() {
	rootBuilder := area.NewAreaBuilder()

	rootBuilder.SetRight(area.NewAbsoluteAnchor(0.0))
	rootBuilder.SetBottom(area.NewAbsoluteAnchor(0.0))
	rootBuilder.OnRender(func(area *area.Area) {
		app.rectRenderer.Fill(area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
			/*graphics.RGBA(0.125, 0.25, 0.45, 1.0) */ graphics.RGBA(0.0, 0.0, 0.0, 1.0))
	})

	app.rootArea = rootBuilder.Build()

	lastBottom := app.rootArea.Top()
	{
		labelBuilder := app.ForLabel()
		labelBuilder.SetParent(app.rootArea)
		labelBuilder.SetRight(app.rootArea.Right())
		labelBuilder.SetTop(lastBottom)
		lastBottom = area.NewOffsetAnchor(lastBottom, 20)
		labelBuilder.SetBottom(lastBottom)
		label1 := labelBuilder.Build()
		label1.SetText("The quick brown fox jumps over the lazy dog 0123456789 :")
	}
	{
		fullScreen := false

		buttonBuilder := app.ForTextButton()
		buttonBuilder.SetParent(app.rootArea)
		buttonBuilder.SetRight(app.rootArea.Right())
		buttonBuilder.SetTop(lastBottom)
		lastBottom = area.NewOffsetAnchor(lastBottom, 20)
		buttonBuilder.SetBottom(lastBottom)
		buttonBuilder.OnAction(func() {
			fmt.Printf("Button click!\n")
			fullScreen = !fullScreen
			app.glWindow.SetFullScreen(fullScreen)
		})
		buttonBuilder.WithText("The Button")
		buttonBuilder.Build()
	}
	{
		boxBuilder := app.ForComboBox()
		boxBuilder.SetParent(app.rootArea)
		boxBuilder.SetRight(app.rootArea.Right())
		boxBuilder.SetTop(lastBottom)
		lastBottom = area.NewOffsetAnchor(lastBottom, 20)
		boxBuilder.SetBottom(lastBottom)
		boxBuilder.WithItems([]controls.ComboBoxItem{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5", "Item 6", "Item 7", "Item 8"})
		boxBuilder.WithSelectionChangeHandler(func(item controls.ComboBoxItem) {
			fmt.Printf("Selected: %v\n", item)
		})
		boxBuilder.Build()
	}
	{
		labelBuilder := app.ForLabel()
		labelBuilder.SetParent(app.rootArea)
		labelBuilder.SetRight(app.rootArea.Right())
		labelBuilder.SetTop(lastBottom)
		lastBottom = area.NewOffsetAnchor(lastBottom, 20)
		labelBuilder.SetBottom(lastBottom)
		labelBuilder.WithTextPainter(app.largeFontPainter)
		labelBuilder.SetScale(1.0)
		label1 := labelBuilder.Build()
		label1.SetText("The quick brown fox jumps over the lazy dog 0123456789 :")
	}
}

func (app *controlsTestApplication) onWindowResize(width int, height int) {
	app.projectionMatrix = mgl.Ortho2D(0.0, float32(width), float32(height), 0.0)
	app.gl.Viewport(0, 0, int32(width), int32(height))

	app.rootArea.Right().RequestValue(float32(width))
	app.rootArea.Bottom().RequestValue(float32(height))
}

func (app *controlsTestApplication) render() {
	gl := app.gl

	gl.Clear(opengl.COLOR_BUFFER_BIT)
	app.rootArea.Render()
}

func (app *controlsTestApplication) onMouseMove(x float32, y float32) {
	app.mouseX, app.mouseY = x, y
	app.rootArea.DispatchPositionalEvent(events.NewMouseMoveEvent(x, y, 0, 0))
}

func (app *controlsTestApplication) onMouseButtonDown(mouseButton uint32, modifier input.Modifier) {
	app.mouseButtons |= mouseButton
	app.rootArea.DispatchPositionalEvent(events.NewMouseButtonEvent(events.MouseButtonDownEventType,
		app.mouseX, app.mouseY, 0, app.mouseButtons, mouseButton))
}

func (app *controlsTestApplication) onMouseButtonUp(mouseButton uint32, modifier input.Modifier) {
	app.mouseButtons &= ^mouseButton
	app.rootArea.DispatchPositionalEvent(events.NewMouseButtonEvent(events.MouseButtonUpEventType,
		app.mouseX, app.mouseY, 0, app.mouseButtons, mouseButton))
}

func (app *controlsTestApplication) onMouseScroll(dx float32, dy float32) {
	app.rootArea.DispatchPositionalEvent(events.NewMouseScrollEvent(
		app.mouseX, app.mouseY, 0, app.mouseButtons, dx, dy))
}

// ForGraphics implements the Context interface.
func (app *controlsTestApplication) ForGraphics() graphics.Context {
	return app
}

// RectangleRenderer implements the graphics.Context interface.
func (app *controlsTestApplication) RectangleRenderer() *graphics.RectangleRenderer {
	return app.rectRenderer
}

// Texturize implements the graphics.Context interface.
func (app *controlsTestApplication) Texturize(bmp *graphics.Bitmap) *graphics.BitmapTexture {
	return graphics.NewBitmapTexture(app.gl, bmp.Width, bmp.Height, bmp.Pixels)
}

// UITextPainter implements the graphics.Context interface.
func (app *controlsTestApplication) UITextPainter() graphics.TextPainter {
	return app.uiFontPainter
}

// UITextRenderer implements the graphics.Context interface.
func (app *controlsTestApplication) UITextRenderer() *graphics.BitmapTextureRenderer {
	return app.uiTextRenderer
}

// NewPaletteTexture implements the graphics.Context interface.
func (app *controlsTestApplication) NewPaletteTexture(colorProvider graphics.ColorProvider) *graphics.PaletteTexture {
	return graphics.NewPaletteTexture(app.gl, colorProvider)
}

// ForLabel implements the controls.Factory interface.
func (app *controlsTestApplication) ForLabel() *controls.LabelBuilder {
	builder := controls.NewLabelBuilder(app.uiFontPainter, app.Texturize, app.uiTextRenderer)
	builder.SetScale(2.0)
	return builder
}

// ForTextButton implements the controls.Factory interface.
func (app *controlsTestApplication) ForTextButton() *controls.TextButtonBuilder {
	return controls.NewTextButtonBuilder(app.ForLabel(), app.rectRenderer)
}

// ForComboBox implements the controls.Factory interface.
func (app *controlsTestApplication) ForComboBox() *controls.ComboBoxBuilder {
	return controls.NewComboBoxBuilder(app.ForLabel(), app.rectRenderer)
}

// ForSlider implements the controls.Factory interface.
func (app *controlsTestApplication) ForSlider() *controls.SliderBuilder {
	return controls.NewSliderBuilder(app.ForLabel(), app.rectRenderer)
}
