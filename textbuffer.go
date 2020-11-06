package main

type TextBuffer interface {
    GetLine(idx int) string
    SetLine(idx int, str string)
    InsertLine(idx int, str string)
    RemoveLine(idx int)
    GetCursor() (line int, char int)
    SetCursor(line, char int)
    NumLines() int
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
    line   := buf.GetLine(cl)

    if cc <= 0 {
        if cl > 0 {
            prev := buf.GetLine(cl-1)
            buf.SetLine(cl-1, prev + line[cc:])
            buf.SetCursor(cl-1, len(prev))
            buf.RemoveLine(cl)
        }
        return
    }

    buf.SetLine(cl, line[:cc-1] + line[cc:])
    buf.SetCursor(cl, cc-1)

}

func TextBufferInsertChar(buf TextBuffer, r rune) {
    cl, cc := buf.GetCursor()
    line := buf.GetLine(cl)
    buf.SetLine(cl, line[:cc] + string(r) + line[cc:])
    buf.SetCursor(cl, cc + 1)
}

func TextBufferMoveCursor(buf TextBuffer, DLine, DChar int) {
    cl, cc := buf.GetCursor()
    nl := buf.NumLines()

    cln := cl + DLine
    if cln < 0 {
        cln = 0
    } else if cln >= nl {
        cln = nl-1
    }

    ccn := cc + DChar
    line := buf.GetLine(cln)
    if ccn < 0 {
        ccn = 0
    } else if ccn > len(line) {
        ccn = len(line)
    }

    buf.SetCursor(cln, ccn)
}
