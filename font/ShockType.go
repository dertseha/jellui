package font

import "github.com/dertseha/jellui/font/text"

var cp850 = text.Codepage850()

// ShockType describes a simple bitmap font, providing only a limited set of characters.
type ShockType struct {
	height int
	stride int
	bitmap []byte

	monochrome bool

	firstCharacter int
	lastCharacter  int
	glyphXOffsets  []int
}

// Monochrome returns true if the type will map only 0x00 or 0x01
func (shock ShockType) Monochrome() bool {
	return shock.monochrome
}

// Height specifies the height of the font.
func (shock ShockType) Height() int {
	return shock.height
}

// Stride specifies the offset to skip in the bitmap to get to the next scanline.
func (shock ShockType) Stride() int {
	return shock.stride
}

// Char returns an entry into the bitmap for given rune. The returned width specifies how many
// pixels are associated with the given rune, in pixels.
func (shock ShockType) Char(r rune) (bitmap []byte, width int) {
	cpIndex := int(cp850.Encode(string(r))[0])
	startOffset := 0

	if (cpIndex >= shock.firstCharacter) && (cpIndex <= shock.lastCharacter) {
		glyphIndex := cpIndex - shock.firstCharacter
		startOffset = shock.glyphXOffsets[glyphIndex]
		width = shock.glyphXOffsets[glyphIndex+1] - startOffset
	}

	return shock.bitmap[startOffset:], width
}
