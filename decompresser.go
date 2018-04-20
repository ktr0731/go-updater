package updater

import (
	"archive/tar"
	"compress/gzip"
	"io"
)

type Decompresser func(io.Reader) (io.Reader, error)

var (
	DefaultDecompresser = TarGZIPDecompresser
	TarDecompresser     = func(r io.Reader) (io.Reader, error) {
		tr := tar.NewReader(r)
		_, err := tr.Next()
		return tr, err
	}
	TarGZIPDecompresser = func(r io.Reader) (io.Reader, error) {
		gr, err := gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		return TarDecompresser(gr)
	}
)
