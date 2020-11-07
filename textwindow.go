package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/tadeuszjt/geom/32"
	"github.com/tadeuszjt/gfx"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
	"image"
)

const (
	textDPI         = 72
	textDefaultSize = 14
)

var (
	trueTypeFont *truetype.Font
)

func init() {
	trueTypeFont, _ = truetype.Parse(gomono.TTF)
}

type TextWindow struct {
	rect          geom.Rect
	tabSize       int
	cursorLine    int
	cursorChar    int
	lines         []string
	face          font.Face
	texID         *gfx.TexID
}

func NewTextWindow(w *gfx.Win, rect geom.Rect) *TextWindow {
	face := truetype.NewFace(trueTypeFont, &truetype.Options{
		Size:    float64(textDefaultSize),
		DPI:     textDPI,
		Hinting: font.HintingFull,
	})

	return &TextWindow{
		rect:       rect,
		tabSize:    4,
		cursorLine: 0,
		cursorChar: 0,
		lines:      []string{""},
		face:       face,
	}
}

func (t *TextWindow) Free(w *gfx.Win) {
	if t.texID != nil {
		w.FreeTexture(*t.texID)
	}
}

func (t *TextWindow) GetLine(idx int) string {
	return t.lines[idx]
}

func (t *TextWindow) SetLine(idx int, str string) {
	t.lines[idx] = str

}

func (t *TextWindow) GetCursor() (line int, char int) {
	return t.cursorLine, t.cursorChar
}

func (t *TextWindow) SetCursor(line, char int) {
	t.cursorChar = char
	t.cursorLine = line
}

func (t *TextWindow) InsertLine(lindIdx int, str string) {
	t.lines = append(t.lines, str)
	copy(t.lines[lindIdx+1:], t.lines[lindIdx:])
	t.lines[lindIdx] = str
}

func (t *TextWindow) NumLines() int {
	return len(t.lines)
}

func (t *TextWindow) RemoveLine(lindIdx int) {
	t.lines = append(t.lines[:lindIdx], t.lines[lindIdx+1:]...)
}

func (t *TextWindow) GetRect() geom.Rect {
	return t.rect
}

func (t *TextWindow) SetRect(w *gfx.Win, r geom.Rect) {
	t.rect = r
    if t.texID != nil {
        w.FreeTexture(*t.texID)
        t.texID = nil
    }
}

func (t *TextWindow) Redraw(w *gfx.Win) {
	width := int(t.rect.Width())
	height := int(t.rect.Height())

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(Solarized.Base00),
		Face: t.face,
		Dot:  fixed.Point26_6{X: 0, Y: t.face.Metrics().Ascent},
	}

	lineHeight := (t.face.Metrics().Ascent + t.face.Metrics().Descent).Ceil()
	charWidth, _ := t.face.GlyphAdvance(' ')

	for i := range t.lines {
		drawer.DrawString(t.lines[i])
		drawer.Dot.X = 0
		drawer.Dot.Y = drawer.Dot.Y + fixed.I(lineHeight)
	}

	if t.texID == nil {
		id := w.LoadTextureFromPixels(width, height, img.Pix)
		t.texID = &id
	} else {
		w.SetTexturePixels(*t.texID, 0, 0, width, height, img.Pix)
	}

	texCanvas := w.GetTextureCanvas(*t.texID)
	gfx.DrawRect(
		texCanvas,
		nil,
		gfx.Colour{0, 1, 0, 0.5},
		geom.MakeRect(
			float32(charWidth.Mul(fixed.I(t.cursorChar)).Ceil()),
			float32(lineHeight*t.cursorLine),
			float32(charWidth.Ceil()),
			float32(lineHeight),
		),
		geom.RectOrigin(1, 1),
	)
}

func (t *TextWindow) DrawOn(c gfx.Canvas) {
	gfx.DrawRect(c, t.texID, gfx.White, t.rect, geom.RectOrigin(1, 1))
}
