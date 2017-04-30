package graphics

// Context is a provider of graphic utilities.
type Context interface {
	RectangleRenderer() *RectangleRenderer
	Texturize(bmp *Bitmap) *BitmapTexture
	UITextPainter() TextPainter
	UITextRenderer() *BitmapTextureRenderer

	NewPaletteTexture(colorProvider ColorProvider) *PaletteTexture
}
