package graphics

type bitmapTextPainter struct {
	font         BitmapFont
	outlineValue byte
}

type charBitmap struct {
	bitmap []byte
	width  int
}

// NewBitmapTextPainter returns a new text painter for the given bitmap font.
func NewBitmapTextPainter(font BitmapFont, outlineValue byte) TextPainter {
	return &bitmapTextPainter{
		font:         font,
		outlineValue: outlineValue}
}

func (painter *bitmapTextPainter) Paint(text string) TextBitmap {
	var bmp TextBitmap
	characterLines := painter.mapCharacters(text)

	bmp.lineHeight = painter.font.Height() + 1
	for _, line := range characterLines {
		lineWidth := 2
		lineOffsets := []int{0}
		for characterOffset, character := range line {
			lineWidth += character.width
			lineOffsets = append(lineOffsets, lineOffsets[characterOffset]+character.width)
		}
		bmp.offsets = append(bmp.offsets, lineOffsets)
		if bmp.Width < lineWidth {
			bmp.Width = lineWidth
		}
	}
	bmp.Height = painter.font.Height()*len(characterLines) + 1 + len(characterLines)
	bmp.Pixels = make([]byte, bmp.Width*bmp.Height)
	for lineIndex, line := range characterLines {
		outStartY := 1 + lineIndex + painter.font.Height()*lineIndex
		outStartX := 1
		for _, character := range line {
			for y := 0; y < painter.font.Height(); y++ {
				inX := painter.font.Stride() * y
				copy(bmp.Pixels[bmp.Width*(outStartY+y)+outStartX:], character.bitmap[inX:inX+character.width])
			}
			outStartX += character.width
		}
	}
	if painter.outlineValue > 0 {
		painter.outline(bmp.Bitmap)
	}

	return bmp
}

func (painter *bitmapTextPainter) mapCharacters(text string) [][]charBitmap {
	lines := [][]charBitmap{}
	curLine := []charBitmap{}

	for _, character := range text {
		if character == '\n' {
			lines = append(lines, curLine)
			curLine = []charBitmap{}
		} else {
			bitmap, width := painter.font.Char(character)
			curLine = append(curLine, charBitmap{bitmap, width})
		}
	}
	lines = append(lines, curLine)

	return lines
}

func (painter *bitmapTextPainter) outline(bmp Bitmap) {
	perimeter := func(index, limit int) (values []int) {
		if index > 0 {
			values = append(values, -1)
		}
		values = append(values, 0)
		if index < (limit - 1) {
			values = append(values, 1)
		}
		return
	}

	for pixelOffset, pixelValue := range bmp.Pixels {
		if pixelValue == 0 {
			lines := perimeter(pixelOffset/bmp.Width, bmp.Height)
			columns := perimeter(pixelOffset%bmp.Width, bmp.Width)
			isNeighbour := false

			for _, lineOffset := range lines {
				for _, columnOffset := range columns {
					if !isNeighbour && (bmp.Pixels[pixelOffset+lineOffset*bmp.Width+columnOffset] == 1) {
						isNeighbour = true
					}
				}
			}
			if isNeighbour {
				bmp.Pixels[pixelOffset] = painter.outlineValue
			}
		}
	}
}
