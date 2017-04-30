package controls

import (
	"github.com/dertseha/jellui/area"
	"github.com/dertseha/jellui/area/events"
	"github.com/dertseha/jellui/graphics"
)

// SliderBuilder is a builder for Slider instances.
type SliderBuilder struct {
	areaBuilder  *area.AreaBuilder
	rectRenderer *graphics.RectangleRenderer
	labelBuilder *LabelBuilder

	sliderChangeHandler SliderChangeHandler

	valueMin int64
	valueMax int64
}

// NewSliderBuilder returns a new SliderBuilder instance.
func NewSliderBuilder(labelBuilder *LabelBuilder, rectRenderer *graphics.RectangleRenderer) *SliderBuilder {
	builder := &SliderBuilder{
		areaBuilder:         area.NewAreaBuilder(),
		rectRenderer:        rectRenderer,
		labelBuilder:        labelBuilder,
		sliderChangeHandler: func(int64) {}}

	return builder
}

// Build creates a new Slider instance from the current parameters.
func (builder *SliderBuilder) Build() *Slider {
	slider := &Slider{
		rectRenderer:        builder.rectRenderer,
		sliderChangeHandler: builder.sliderChangeHandler,
		valueMin:            builder.valueMin,
		valueMax:            builder.valueMax,
		valueUndefined:      true}

	builder.areaBuilder.OnRender(slider.onRender)
	builder.areaBuilder.OnEvent(events.MouseButtonDownEventType, slider.onMouseButtonDown)
	builder.areaBuilder.OnEvent(events.MouseButtonUpEventType, slider.onMouseButtonUp)
	builder.areaBuilder.OnEvent(events.MouseMoveEventType, slider.onMouseMove)
	builder.areaBuilder.OnEvent(events.MouseScrollEventType, slider.onMouseScroll)
	builder.areaBuilder.OnEvent(events.MouseButtonClickedEventType, area.SilentConsumer)
	slider.area = builder.areaBuilder.Build()

	builder.labelBuilder.SetParent(slider.area)
	builder.labelBuilder.SetLeft(area.NewOffsetAnchor(slider.area.Left(), 4))
	builder.labelBuilder.SetTop(area.NewOffsetAnchor(slider.area.Top(), 0))
	builder.labelBuilder.SetRight(area.NewOffsetAnchor(slider.area.Right(), -4))
	builder.labelBuilder.SetBottom(area.NewOffsetAnchor(slider.area.Bottom(), 0))
	builder.labelBuilder.AlignedHorizontallyBy(LeftAligner)

	slider.valueLabel = builder.labelBuilder.Build()

	return slider
}

// SetParent sets the parent area.
func (builder *SliderBuilder) SetParent(parent *area.Area) *SliderBuilder {
	builder.areaBuilder.SetParent(parent)
	return builder
}

// SetLeft sets the left anchor. Default: ZeroAnchor
func (builder *SliderBuilder) SetLeft(value area.Anchor) *SliderBuilder {
	builder.areaBuilder.SetLeft(value)
	return builder
}

// SetTop sets the top anchor. Default: ZeroAnchor
func (builder *SliderBuilder) SetTop(value area.Anchor) *SliderBuilder {
	builder.areaBuilder.SetTop(value)
	return builder
}

// SetRight sets the right anchor. Default: ZeroAnchor
func (builder *SliderBuilder) SetRight(value area.Anchor) *SliderBuilder {
	builder.areaBuilder.SetRight(value)
	return builder
}

// SetBottom sets the bottom anchor. Default: ZeroAnchor
func (builder *SliderBuilder) SetBottom(value area.Anchor) *SliderBuilder {
	builder.areaBuilder.SetBottom(value)
	return builder
}

// WithSliderChangeHandler sets the handler for a value change.
func (builder *SliderBuilder) WithSliderChangeHandler(handler SliderChangeHandler) *SliderBuilder {
	builder.sliderChangeHandler = handler
	return builder
}

// WithRange sets the allowed range of the slider.
func (builder *SliderBuilder) WithRange(valueMin, valueMax int64) *SliderBuilder {
	builder.valueMin = valueMin
	builder.valueMax = valueMax
	return builder
}
