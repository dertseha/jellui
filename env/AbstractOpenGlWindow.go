package env

import (
	"github.com/dertseha/jellui/input"
)

type keyDeferrer struct {
	window *AbstractOpenGlWindow
}

func (def *keyDeferrer) Key(key input.Key, modifier input.Modifier) {
	def.window.CallKey(key, modifier)
}

func (def *keyDeferrer) Modifier(modifier input.Modifier) {
	def.window.CallModifier(modifier)
}

// AbstractOpenGlWindow implements the common, basic functionality of OpenGlWindow.
type AbstractOpenGlWindow struct {
	keyBuffer *input.StickyKeyBuffer

	CallRender            RenderCallback
	CallResize            ResizeCallback
	CallOnMouseMove       MouseMoveCallback
	CallOnMouseButtonUp   MouseButtonCallback
	CallOnMouseButtonDown MouseButtonCallback
	CallOnMouseScroll     MouseScrollCallback
	CallModifier          ModifierCallback
	CallKey               KeyCallback
	CallCharCallback      CharCallback
}

// InitAbstractOpenGlWindow returns an initialized instance.
func InitAbstractOpenGlWindow() AbstractOpenGlWindow {
	return AbstractOpenGlWindow{
		CallRender:            func() {},
		CallResize:            func(int, int) {},
		CallOnMouseMove:       func(float32, float32) {},
		CallOnMouseButtonUp:   func(uint32, input.Modifier) {},
		CallOnMouseButtonDown: func(uint32, input.Modifier) {},
		CallOnMouseScroll:     func(float32, float32) {},
		CallKey:               func(input.Key, input.Modifier) {},
		CallModifier:          func(input.Modifier) {},
		CallCharCallback:      func(rune) {}}
}

// StickyKeyListener returns an instance of a listener acting as an adapter
// for the key-down/-up callbacks.
func (window *AbstractOpenGlWindow) StickyKeyListener() input.StickyKeyListener {
	return &keyDeferrer{window}
}

// OnRender implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnRender(callback RenderCallback) {
	window.CallRender = callback
}

// OnResize implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnResize(callback ResizeCallback) {
	window.CallResize = callback
}

// OnMouseMove implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseMove(callback MouseMoveCallback) {
	window.CallOnMouseMove = callback
}

// OnMouseButtonDown implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseButtonDown(callback MouseButtonCallback) {
	window.CallOnMouseButtonDown = callback
}

// OnMouseButtonUp implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseButtonUp(callback MouseButtonCallback) {
	window.CallOnMouseButtonUp = callback
}

// OnMouseScroll implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseScroll(callback MouseScrollCallback) {
	window.CallOnMouseScroll = callback
}

// OnKey implements the OpenGlWindow interface
func (window *AbstractOpenGlWindow) OnKey(callback KeyCallback) {
	window.CallKey = callback
}

// OnModifier implements the OpenGlWindow interface
func (window *AbstractOpenGlWindow) OnModifier(callback ModifierCallback) {
	window.CallModifier = callback
}

// OnCharCallback implements the OpenGlWindow interface
func (window *AbstractOpenGlWindow) OnCharCallback(callback CharCallback) {
	window.CallCharCallback = callback
}
