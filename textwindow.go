package main

import (
	"github.com/golang/freetype/truetype"
	"github.com/tadeuszjt/geom/32"
	"github.com/tadeuszjt/gfx"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
	"image"
	"strings"
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
	rect        geom.Rect
	tabSize     int
	cursorLine  int
	cursorChar  int
	scrollLine  int
	lines       []string
	face        font.Face
	texID       *gfx.TexID
	drawBoarder bool
}

func NewTextWindow(w *gfx.Win, rect geom.Rect) *TextWindow {
	face := truetype.NewFace(trueTypeFont, &truetype.Options{
		Size:    float64(textDefaultSize),
		DPI:     textDPI,
		Hinting: font.HintingFull,
	})

	return &TextWindow{
		rect:    rect,
		tabSize: 4,
		lines:   []string{""},
		face:    face,
	}
}

func (t *TextWindow) Free(w *gfx.Win) {
	if t.texID != nil {
		w.FreeTexture(*t.texID)
	}
}

func (t *TextWindow) Scroll(num int) {
	t.scrollLine += num
	if t.scrollLine < 0 {
		t.scrollLine = 0
	} else if t.scrollLine > len(t.lines) {
		t.scrollLine = len(t.lines)
	}
}

func (t *TextWindow) NumLinesPerPage() int {
	return int(t.rect.Height()) / t.LineHeight()
}

func (t *TextWindow) GetLine(idx int) string {
	return t.lines[idx]
}

func (t *TextWindow) GetTabSize() int {
	return t.tabSize
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

func (t *TextWindow) LineHeight() int {
	return (t.face.Metrics().Ascent + t.face.Metrics().Descent).Ceil()
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

	charWidth, _ := t.face.GlyphAdvance(' ')
	lineHeight := t.LineHeight()

	for i := t.scrollLine; i < len(t.lines); i++ {
		s := strings.ReplaceAll(t.lines[i], "\t", strings.Repeat(" ", t.tabSize))
		drawer.DrawString(s)
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
	cursorPos := charPos(t.GetLine(t.cursorLine), t.cursorChar, t.tabSize)

	gfx.DrawRect(
		texCanvas,
		nil,
		gfx.Colour{0, 1, 0, 0.5},
		geom.MakeRect(
			float32(charWidth.Mul(fixed.I(cursorPos)).Ceil()),
			float32(lineHeight*(t.cursorLine-t.scrollLine)),
			float32(charWidth.Ceil()),
			float32(lineHeight),
		),
		geom.RectOrigin(1, 1),
	)

	if t.drawBoarder {
		colour := Solarized.Cyan
		gfx.DrawRect(texCanvas, nil, colour, geom.MakeRect(0, 0, t.rect.Width(), 1), geom.RectOrigin(1, 1))
		gfx.DrawRect(texCanvas, nil, colour, geom.MakeRect(0, t.rect.Height()-1, t.rect.Width(), 1), geom.RectOrigin(1, 1))
		gfx.DrawRect(texCanvas, nil, colour, geom.MakeRect(0, 0, 1, t.rect.Height()), geom.RectOrigin(1, 1))
		gfx.DrawRect(texCanvas, nil, colour, geom.MakeRect(t.rect.Width()-1, 0, 1, t.rect.Height()), geom.RectOrigin(1, 1))
	}
}

func (t *TextWindow) DrawOn(c gfx.Canvas) {
	gfx.DrawRect(c, t.texID, gfx.White, t.rect, geom.RectOrigin(1, 1))
}
