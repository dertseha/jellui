package env

import (
	"github.com/dertseha/jellui"
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

	CallRender            jellui.RenderCallback
	CallResize            jellui.ResizeCallback
	CallOnMouseMove       jellui.MouseMoveCallback
	CallOnMouseButtonUp   jellui.MouseButtonCallback
	CallOnMouseButtonDown jellui.MouseButtonCallback
	CallOnMouseScroll     jellui.MouseScrollCallback
	CallModifier          jellui.ModifierCallback
	CallKey               jellui.KeyCallback
	CallCharCallback      jellui.CharCallback
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
func (window *AbstractOpenGlWindow) OnRender(callback jellui.RenderCallback) {
	window.CallRender = callback
}

// OnResize implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnResize(callback jellui.ResizeCallback) {
	window.CallResize = callback
}

// OnMouseMove implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseMove(callback jellui.MouseMoveCallback) {
	window.CallOnMouseMove = callback
}

// OnMouseButtonDown implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseButtonDown(callback jellui.MouseButtonCallback) {
	window.CallOnMouseButtonDown = callback
}

// OnMouseButtonUp implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseButtonUp(callback jellui.MouseButtonCallback) {
	window.CallOnMouseButtonUp = callback
}

// OnMouseScroll implements the OpenGlWindow interface.
func (window *AbstractOpenGlWindow) OnMouseScroll(callback jellui.MouseScrollCallback) {
	window.CallOnMouseScroll = callback
}

// OnKey implements the OpenGlWindow interface
func (window *AbstractOpenGlWindow) OnKey(callback jellui.KeyCallback) {
	window.CallKey = callback
}

// OnModifier implements the OpenGlWindow interface
func (window *AbstractOpenGlWindow) OnModifier(callback jellui.ModifierCallback) {
	window.CallModifier = callback
}

// OnCharCallback implements the OpenGlWindow interface
func (window *AbstractOpenGlWindow) OnCharCallback(callback jellui.CharCallback) {
	window.CallCharCallback = callback
}
