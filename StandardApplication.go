package jellui

import (
	"fmt"
	"os"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/dertseha/jellui/area"
	"github.com/dertseha/jellui/area/events"
	"github.com/dertseha/jellui/controls"
	"github.com/dertseha/jellui/env"
	"github.com/dertseha/jellui/font"
	"github.com/dertseha/jellui/graphics"
	"github.com/dertseha/jellui/input"
	"github.com/dertseha/jellui/opengl"
)

// StandardApplication is a typical implementation of jellui framework.
// It implements the graphics.Context and controls.Factory interfaces.
type StandardApplication struct {
	glWindow env.OpenGlWindow
	gl       opengl.OpenGl

	projectionMatrix mgl.Mat4

	mouseX, mouseY float32
	mouseButtons   uint32

	uiFontPainter        graphics.TextPainter
	uiTextScale          float32
	uiTextPalette        map[int][4]byte
	uiTextPaletteTexture *graphics.PaletteTexture
	rectRenderer         *graphics.RectangleRenderer
	uiTextRenderer       *graphics.BitmapTextureRenderer

	rootArea *area.Area

	uiSetup func(*StandardApplication, *area.Area)
}

// NewStandardApplication returns an instance of a standard application. The provided uiSetup
// function will be called to initialize the UI.
func NewStandardApplication(uiSetup func(*StandardApplication, *area.Area)) *StandardApplication {
	return &StandardApplication{
		uiSetup:     uiSetup,
		uiTextScale: 1.0,
		uiTextPalette: map[int][4]byte{
			0: {0x00, 0x00, 0x00, 0x00},
			1: {0x80, 0x94, 0x54, 0xFF},
			2: {0x00, 0x00, 0x00, 0xC0}}}
}

// SetCursorVisible controls whether the mouse cursor shall be visible.
func (app *StandardApplication) SetCursorVisible(visible bool) {
	app.glWindow.SetCursorVisible(visible)
}

// SetFullScreen sets the full screen state of the window.
func (app *StandardApplication) SetFullScreen(on bool) {
	app.glWindow.SetFullScreen(on)
}

// SetUITextScale sets the default scale factor for any created labels.
func (app *StandardApplication) SetUITextScale(scale float32) {
	app.uiTextScale = scale
}

// SetUITextPalette sets the palette for UI texts.
func (app *StandardApplication) SetUITextPalette(palette map[int][4]byte) {
	for key, color := range palette {
		app.uiTextPalette[key] = color
	}
	app.uiTextPaletteTexture.Update()
}

// Init implements the Application interface.
func (app *StandardApplication) Init(glWindow env.OpenGlWindow) {
	app.setWindow(glWindow)
	app.initOpenGl()
	app.setDebugOpenGl()

	app.initGraphics()
	app.initInterface()

	app.onWindowResize(glWindow.Size())
}

func (app *StandardApplication) setWindow(glWindow env.OpenGlWindow) {
	app.glWindow = glWindow
	app.gl = glWindow.OpenGl()

	glWindow.OnRender(app.render)
	glWindow.OnResize(app.onWindowResize)

	glWindow.OnMouseMove(app.onMouseMove)
	glWindow.OnMouseButtonDown(app.onMouseButtonDown)
	glWindow.OnMouseButtonUp(app.onMouseButtonUp)
	glWindow.OnMouseScroll(app.onMouseScroll)
}

func (app *StandardApplication) initOpenGl() {
	app.gl.Disable(opengl.DEPTH_TEST)
	app.gl.Enable(opengl.BLEND)
	app.gl.BlendFunc(opengl.SRC_ALPHA, opengl.ONE_MINUS_SRC_ALPHA)
	app.gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func (app *StandardApplication) setDebugOpenGl() {
	builder := opengl.NewDebugBuilder(app.gl)

	builder.OnError(func(name string, errorCodes []uint32) {
		errorStrings := make([]string, len(errorCodes))
		for index, errorCode := range errorCodes {
			errorStrings[index] = opengl.ErrorString(errorCode)
		}
		fmt.Fprintf(os.Stderr, "!!: [%-20s] %v -> %v\n", name, errorCodes, errorStrings)
	})

	app.gl = builder.Build()
}

func (app *StandardApplication) initGraphics() {
	app.uiTextPaletteTexture = app.NewPaletteTexture(func(index int) (byte, byte, byte, byte) {
		entry := app.uiTextPalette[index]
		return entry[0], entry[1], entry[2], entry[3]
	})
	viewMatrix := mgl.Ident4()
	uiRenderContext := graphics.NewBasicRenderContext(app.gl, &app.projectionMatrix, &viewMatrix)
	app.uiTextRenderer = graphics.NewBitmapTextureRenderer(uiRenderContext, app.uiTextPaletteTexture)

	app.uiFontPainter = graphics.NewBitmapTextPainter(font.SmallShock, 0x02)

	app.rectRenderer = graphics.NewRectangleRenderer(app.gl, &app.projectionMatrix)
}

func (app *StandardApplication) initInterface() {
	rootBuilder := area.NewAreaBuilder()

	rootBuilder.SetRight(area.NewAbsoluteAnchor(0.0))
	rootBuilder.SetBottom(area.NewAbsoluteAnchor(0.0))

	app.rootArea = rootBuilder.Build()
	app.uiSetup(app, app.rootArea)
}

func (app *StandardApplication) onWindowResize(width int, height int) {
	app.projectionMatrix = mgl.Ortho2D(0.0, float32(width), float32(height), 0.0)
	app.gl.Viewport(0, 0, int32(width), int32(height))

	app.rootArea.Right().RequestValue(float32(width))
	app.rootArea.Bottom().RequestValue(float32(height))
}

func (app *StandardApplication) render() {
	gl := app.gl

	gl.Clear(opengl.COLOR_BUFFER_BIT)
	app.rootArea.Render()
}

func (app *StandardApplication) onMouseMove(x float32, y float32) {
	app.mouseX, app.mouseY = x, y
	app.rootArea.DispatchPositionalEvent(events.NewMouseMoveEvent(x, y, 0, 0))
}

func (app *StandardApplication) onMouseButtonDown(mouseButton uint32, modifier input.Modifier) {
	app.mouseButtons |= mouseButton
	app.rootArea.DispatchPositionalEvent(events.NewMouseButtonEvent(events.MouseButtonDownEventType,
		app.mouseX, app.mouseY, 0, app.mouseButtons, mouseButton))
}

func (app *StandardApplication) onMouseButtonUp(mouseButton uint32, modifier input.Modifier) {
	app.mouseButtons &= ^mouseButton
	app.rootArea.DispatchPositionalEvent(events.NewMouseButtonEvent(events.MouseButtonUpEventType,
		app.mouseX, app.mouseY, 0, app.mouseButtons, mouseButton))
}

func (app *StandardApplication) onMouseScroll(dx float32, dy float32) {
	app.rootArea.DispatchPositionalEvent(events.NewMouseScrollEvent(
		app.mouseX, app.mouseY, 0, app.mouseButtons, dx, dy))
}

// RectangleRenderer implements the graphics.Context interface.
func (app *StandardApplication) RectangleRenderer() *graphics.RectangleRenderer {
	return app.rectRenderer
}

// Texturize implements the graphics.Context interface.
func (app *StandardApplication) Texturize(bmp *graphics.Bitmap) *graphics.BitmapTexture {
	return graphics.NewBitmapTexture(app.gl, bmp.Width, bmp.Height, bmp.Pixels)
}

// UITextPainter implements the graphics.Context interface.
func (app *StandardApplication) UITextPainter() graphics.TextPainter {
	return app.uiFontPainter
}

// UITextRenderer implements the graphics.Context interface.
func (app *StandardApplication) UITextRenderer() *graphics.BitmapTextureRenderer {
	return app.uiTextRenderer
}

// NewPaletteTexture implements the graphics.Context interface.
func (app *StandardApplication) NewPaletteTexture(colorProvider graphics.ColorProvider) *graphics.PaletteTexture {
	return graphics.NewPaletteTexture(app.gl, colorProvider)
}

// ForLabel implements the controls.Factory interface.
func (app *StandardApplication) ForLabel() *controls.LabelBuilder {
	builder := controls.NewLabelBuilder(app.uiFontPainter, app.Texturize, app.uiTextRenderer)
	builder.SetScale(app.uiTextScale)
	return builder
}

// ForTextButton implements the controls.Factory interface.
func (app *StandardApplication) ForTextButton() *controls.TextButtonBuilder {
	return controls.NewTextButtonBuilder(app.ForLabel(), app.rectRenderer)
}

// ForComboBox implements the controls.Factory interface.
func (app *StandardApplication) ForComboBox() *controls.ComboBoxBuilder {
	return controls.NewComboBoxBuilder(app.ForLabel(), app.rectRenderer)
}

// ForSlider implements the controls.Factory interface.
func (app *StandardApplication) ForSlider() *controls.SliderBuilder {
	return controls.NewSliderBuilder(app.ForLabel(), app.rectRenderer)
}
