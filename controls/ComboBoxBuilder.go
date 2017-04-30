package controls

import (
	"github.com/dertseha/jellui/area"
	"github.com/dertseha/jellui/area/events"
	"github.com/dertseha/jellui/graphics"
)

// ComboBoxBuilder is a builder for ComboBox instances.
type ComboBoxBuilder struct {
	areaBuilder  *area.AreaBuilder
	rectRenderer *graphics.RectangleRenderer
	labelBuilder *LabelBuilder

	selectionChangeHandler SelectionChangeHandler

	items []ComboBoxItem
}

// NewComboBoxBuilder returns a new ComboBoxBuilder instance.
func NewComboBoxBuilder(labelBuilder *LabelBuilder, rectRenderer *graphics.RectangleRenderer) *ComboBoxBuilder {
	builder := &ComboBoxBuilder{
		areaBuilder:            area.NewAreaBuilder(),
		rectRenderer:           rectRenderer,
		labelBuilder:           labelBuilder,
		selectionChangeHandler: func(ComboBoxItem) {}}

	return builder
}

// Build creates a new ComboBox instance from the current parameters.
func (builder *ComboBoxBuilder) Build() *ComboBox {
	box := &ComboBox{
		labelBuilder:           builder.labelBuilder,
		rectRenderer:           builder.rectRenderer,
		selectionChangeHandler: builder.selectionChangeHandler,
		items: builder.items}

	builder.areaBuilder.OnRender(box.onRender)
	builder.areaBuilder.OnEvent(events.MouseButtonDownEventType, box.onMouseDown)
	builder.areaBuilder.OnEvent(events.MouseButtonUpEventType, area.SilentConsumer)
	builder.areaBuilder.OnEvent(events.MouseButtonClickedEventType, area.SilentConsumer)
	builder.areaBuilder.OnEvent(events.MouseScrollEventType, area.SilentConsumer)
	box.area = builder.areaBuilder.Build()

	builder.labelBuilder.SetParent(box.area)
	builder.labelBuilder.SetTop(area.NewOffsetAnchor(box.area.Top(), 0))
	builder.labelBuilder.SetBottom(area.NewOffsetAnchor(box.area.Bottom(), 0))

	hintRight := area.NewOffsetAnchor(box.area.Right(), 0)
	hintLeft := area.NewOffsetAnchor(hintRight, -25)
	builder.labelBuilder.SetLeft(hintLeft)
	builder.labelBuilder.SetRight(hintRight)
	box.hintLabel = builder.labelBuilder.Build()
	box.hintLabel.SetText("...")

	builder.labelBuilder.SetLeft(area.NewOffsetAnchor(box.area.Left(), 4))
	builder.labelBuilder.SetRight(area.NewOffsetAnchor(hintLeft, -4))
	builder.labelBuilder.AlignedHorizontallyBy(LeftAligner)
	box.selectedLabel = builder.labelBuilder.Build()

	return box
}

// SetParent sets the parent area.
func (builder *ComboBoxBuilder) SetParent(parent *area.Area) *ComboBoxBuilder {
	builder.areaBuilder.SetParent(parent)
	return builder
}

// SetLeft sets the left anchor. Default: ZeroAnchor
func (builder *ComboBoxBuilder) SetLeft(value area.Anchor) *ComboBoxBuilder {
	builder.areaBuilder.SetLeft(value)
	return builder
}

// SetTop sets the top anchor. Default: ZeroAnchor
func (builder *ComboBoxBuilder) SetTop(value area.Anchor) *ComboBoxBuilder {
	builder.areaBuilder.SetTop(value)
	return builder
}

// SetRight sets the right anchor. Default: ZeroAnchor
func (builder *ComboBoxBuilder) SetRight(value area.Anchor) *ComboBoxBuilder {
	builder.areaBuilder.SetRight(value)
	return builder
}

// SetBottom sets the bottom anchor. Default: ZeroAnchor
func (builder *ComboBoxBuilder) SetBottom(value area.Anchor) *ComboBoxBuilder {
	builder.areaBuilder.SetBottom(value)
	return builder
}

// WithItems sets the list of contained items.
func (builder *ComboBoxBuilder) WithItems(items []ComboBoxItem) *ComboBoxBuilder {
	builder.items = make([]ComboBoxItem, len(items))
	copy(builder.items, items)
	return builder
}

// WithSelectionChangeHandler sets the handler for a selection change.
func (builder *ComboBoxBuilder) WithSelectionChangeHandler(handler SelectionChangeHandler) *ComboBoxBuilder {
	builder.selectionChangeHandler = handler
	return builder
}
