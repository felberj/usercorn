package loader

import (
	"bytes"
	"io"
)

var cgcMagic = []byte{0x7f, 0x43, 0x47, 0x43}
var elfMagicReader = bytes.NewReader(elfMagic)

func MatchCgc(r io.ReaderAt) bool {
	return bytes.Equal(getMagic(r), cgcMagic)
}

type FakeCgcReader struct {
	io.ReaderAt
	first bool
}

func (f *FakeCgcReader) ReadAt(p []byte, off int64) (int, error) {
	n := 0
	if off < 4 && f.first {
		f.first = false
		n, _ = elfMagicReader.ReadAt(p, off)
		if n == len(p) {
			return n, nil
		}
		p = p[n:]
		off = 4
	}
	n1, err := f.ReaderAt.ReadAt(p, off)
	return n1 + n, err
}

type CgcLoader struct {
	Loader
}

func (c *CgcLoader) OS() string {
	return "cgc"
}

func NewCgcLoader(r io.ReaderAt, arch string) (Loader, error) {
	l, err := NewElfLoader(&FakeCgcReader{r, true}, arch, NoOSHint)
	return &CgcLoader{l}, err
}
