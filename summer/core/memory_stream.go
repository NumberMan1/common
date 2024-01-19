package core

import "io"

type MemoryStream struct {
	buffer []byte
	pos    int
}

func NewMemoryStream(size int) *MemoryStream {
	return &MemoryStream{
		buffer: make([]byte, size),
		pos:    0,
	}
}

func (ms *MemoryStream) Write(p []byte) (n int, err error) {
	if ms.pos+len(p) > len(ms.buffer) {
		return 0, io.ErrShortWrite
	}
	copy(ms.buffer[ms.pos:], p)
	ms.pos += len(p)
	return len(p), nil
}

func (ms *MemoryStream) Read(p []byte) (n int, err error) {
	if ms.pos >= len(ms.buffer) {
		return 0, io.EOF
	}
	n = copy(p, ms.buffer[ms.pos:])
	ms.pos += n
	return n, nil
}

func (ms *MemoryStream) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		ms.pos = int(offset)
	case io.SeekCurrent:
		ms.pos += int(offset)
	case io.SeekEnd:
		ms.pos = len(ms.buffer) - int(offset)
	}
	return int64(ms.pos), nil
}
