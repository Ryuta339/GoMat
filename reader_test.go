package gomat

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TEST_COMPRESSED_FILE = "testdata/compressed.mat"
)

func TestReadHeader(t *testing.T) {
	tests := []struct {
		filename    string
		description string
		endian      binary.ByteOrder
	}{
		{
			filename:    TEST_COMPRESSED_FILE,
			description: "MATLAB 5.0 MAT-file, Platform: MACI64, Create on: Fri Jul 15 18:55:56 2022",
			endian:      binary.LittleEndian,
		},
	}
	for _, tt := range tests {
		f, err := os.Open(tt.filename)
		assert.NoError(t, err, "read file error")
		defer f.Close()

		mh := &MatHeader{}
		err = mh.readHeader(f)
		assert.NoError(t, err, "read header error")

		assert.Equal(t, tt.description, mh.String(), "Header Description Text is Failed")
	}
}
