package main

type TextBuffer interface {
	GetLine(idx int) string
	SetLine(idx int, str string)
	InsertLine(idx int, str string)
	RemoveLine(idx int)
	GetCursor() (line int, char int)
	SetCursor(line, char int)
	GetTabSize() int
	NumLines() int
}

func TextBufferString(buf TextBuffer) (s string) {
	for i := 0; i < buf.NumLines(); i++ {
		s = s + buf.GetLine(i) + "\n"
	}
	return
}

func TextBufferClear(buf TextBuffer) {
	for i := buf.NumLines() - 1; i > 0; i-- {
		buf.RemoveLine(i)
	}
	buf.SetLine(0, "")
}

func TextBufferAppend(buf TextBuffer, str string) {
	buf.InsertLine(buf.NumLines(), str)
}

func TextBufferEnter(buf TextBuffer) {
	cl, cc := buf.GetCursor()
	line := buf.GetLine(cl)
	buf.SetLine(cl, line[:cc])
	buf.InsertLine(cl+1, line[cc:])
	buf.SetCursor(cl+1, 0)
}

func TextBufferBackspace(buf TextBuffer) {
	cl, cc := buf.GetCursor()
	line := buf.GetLine(cl)

	if cc <= 0 {
		if cl > 0 {
			prev := buf.GetLine(cl - 1)
			buf.SetLine(cl-1, prev+line[cc:])
			buf.SetCursor(cl-1, len(prev))
			buf.RemoveLine(cl)
		}
		return
	}

	buf.SetLine(cl, line[:cc-1]+line[cc:])
	buf.SetCursor(cl, cc-1)
}

func TextBufferInsertChar(buf TextBuffer, r rune) {
	cl, cc := buf.GetCursor()
	line := buf.GetLine(cl)
	buf.SetLine(cl, line[:cc]+string(r)+line[cc:])
	buf.SetCursor(cl, cc+1)
}

func TextBufferMoveCursor(buf TextBuffer, DLine, DChar int) {
	curLine, curChar := buf.GetCursor()
	numLines := buf.NumLines()
	tabSize := buf.GetTabSize()

	newLine := curLine + DLine
	if newLine < 0 {
		newLine = 0
	} else if newLine >= numLines {
		newLine = numLines - 1
	}

    curPos := charPos(buf.GetLine(curLine), curChar, tabSize)
    newLineStr := buf.GetLine(newLine)

    newChar := 0
    for ; newChar < len(newLineStr); newChar++ {
        if charPos(newLineStr, newChar, tabSize) >= curPos {
            break
        } else if charPos(newLineStr, newChar+1, tabSize) > curPos {
            break
        }
    }

	newChar += DChar
    if newChar > len(newLineStr) {
		newChar = len(newLineStr)
	} else if newChar < 0 {
		newChar = 0
	}

	buf.SetCursor(newLine, newChar)
}

func charPos(str string, charIdx, tabSize int) int {
    pos := 0
    for i, r := range str {
        if i >= charIdx {
            break
        }
        if r == '\t' {
            pos += tabSize
        } else {
            pos += 1
        }
    }

    return pos
}
