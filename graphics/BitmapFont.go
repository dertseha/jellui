package graphics

// BitmapFont describes a font based on pixels.
type BitmapFont interface {
	// Height specifies the height of the font.
	Height() int
	// Stride specifies the offset to skip in the bitmap to get to the next scanline.
	Stride() int
	// Char returns an entry into the bitmap for given rune. The returned width specifies how many
	// pixels are associated with the given rune, in pixels.
	Char(r rune) (bitmap []byte, width int)
}
