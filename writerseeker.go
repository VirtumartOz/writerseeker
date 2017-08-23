package writerseeker

import (
	"bytes"
	"errors"
	"io"
)

// WriterSeeker is an in-memory io.WriteSeeker implementation
type WriterSeeker struct {
	buf []byte
	pos int
}

func (ws *WriterSeeker) Write(p []byte) (n int, err error) {
	minCap := ws.pos + len(p)
	if minCap > cap(ws.buf) { // Make sure buf has enough capacity:
		buf2 := make([]byte, len(ws.buf), minCap+len(p)) // add some extra
		copy(buf2, ws.buf)
		ws.buf = buf2
	}
	if minCap > len(ws.buf) {
		ws.buf = ws.buf[:minCap]
	}
	copy(ws.buf[ws.pos:], p)
	ws.pos += len(p)
	return len(p), nil
}

func (ws *WriterSeeker) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = ws.pos + offs
	case io.SeekEnd:
		newPos = len(ws.buf) - offs
	}
	if newPos < 0 {
		return 0, errors.New("negative result pos")
	}
	ws.pos = newPos
	return int64(newPos), nil
}

func (ws *WriterSeeker) Reader() io.Reader {
	return bytes.NewReader(ws.buf)
}