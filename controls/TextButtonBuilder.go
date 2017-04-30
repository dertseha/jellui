package controls

import (
	"github.com/dertseha/jellui/area"
	"github.com/dertseha/jellui/area/events"
	"github.com/dertseha/jellui/graphics"
)

// TextButtonBuilder is a builder for TextButton instances.
type TextButtonBuilder struct {
	areaBuilder  *area.AreaBuilder
	rectRenderer *graphics.RectangleRenderer

	idleColor     graphics.Color
	preparedColor graphics.Color

	labelBuilder *LabelBuilder
	text         string

	actionHandler ActionHandler
}

// NewTextButtonBuilder returns a new TextButtonBuilder instance.
func NewTextButtonBuilder(labelBuilder *LabelBuilder, rectRenderer *graphics.RectangleRenderer) *TextButtonBuilder {
	builder := &TextButtonBuilder{
		areaBuilder:   area.NewAreaBuilder(),
		rectRenderer:  rectRenderer,
		idleColor:     graphics.RGBA(0.31, 0.56, 0.34, 0.8),
		preparedColor: graphics.RGBA(0.31, 0.56, 0.34, 0.95),
		labelBuilder:  labelBuilder,
		text:          "",
		actionHandler: func() {}}

	return builder
}

// Build creates a new TextButton instance from the current parameters.
func (builder *TextButtonBuilder) Build() *TextButton {
	button := &TextButton{
		rectRenderer:  builder.rectRenderer,
		idleColor:     builder.idleColor,
		preparedColor: builder.preparedColor,
		color:         builder.idleColor,
		actionHandler: builder.actionHandler}

	builder.areaBuilder.OnRender(button.onRender)
	builder.areaBuilder.OnEvent(events.MouseButtonDownEventType, button.onMouseDown)
	builder.areaBuilder.OnEvent(events.MouseButtonUpEventType, button.onMouseUp)
	button.area = builder.areaBuilder.Build()

	button.labelLeft = area.NewOffsetAnchor(button.area.Left(), 0)
	button.labelTop = area.NewOffsetAnchor(button.area.Top(), 0)

	builder.labelBuilder.SetParent(button.area)
	builder.labelBuilder.SetLeft(button.labelLeft)
	builder.labelBuilder.SetTop(button.labelTop)
	builder.labelBuilder.SetRight(area.NewOffsetAnchor(button.area.Right(), 0))
	builder.labelBuilder.SetBottom(area.NewOffsetAnchor(button.area.Bottom(), 0))

	button.label = builder.labelBuilder.Build()
	button.label.SetText(builder.text)

	return button
}

// SetParent sets the parent area.
func (builder *TextButtonBuilder) SetParent(parent *area.Area) *TextButtonBuilder {
	builder.areaBuilder.SetParent(parent)
	return builder
}

// SetLeft sets the left anchor. Default: ZeroAnchor
func (builder *TextButtonBuilder) SetLeft(value area.Anchor) *TextButtonBuilder {
	builder.areaBuilder.SetLeft(value)
	return builder
}

// SetTop sets the top anchor. Default: ZeroAnchor
func (builder *TextButtonBuilder) SetTop(value area.Anchor) *TextButtonBuilder {
	builder.areaBuilder.SetTop(value)
	return builder
}

// SetRight sets the right anchor. Default: ZeroAnchor
func (builder *TextButtonBuilder) SetRight(value area.Anchor) *TextButtonBuilder {
	builder.areaBuilder.SetRight(value)
	return builder
}

// SetBottom sets the bottom anchor. Default: ZeroAnchor
func (builder *TextButtonBuilder) SetBottom(value area.Anchor) *TextButtonBuilder {
	builder.areaBuilder.SetBottom(value)
	return builder
}

// WithText sets the label text to be used for the new button.
func (builder *TextButtonBuilder) WithText(value string) *TextButtonBuilder {
	builder.text = value
	return builder
}

// OnAction sets the action handler of the new button.
func (builder *TextButtonBuilder) OnAction(handler ActionHandler) *TextButtonBuilder {
	builder.actionHandler = handler
	return builder
}

// WithIdleColor sets the idle background color.
func (builder *TextButtonBuilder) WithIdleColor(color graphics.Color) *TextButtonBuilder {
	builder.idleColor = color
	return builder
}

// WithPreparedColor sets the background color for the prepared state.
func (builder *TextButtonBuilder) WithPreparedColor(color graphics.Color) *TextButtonBuilder {
	builder.preparedColor = color
	return builder
}
