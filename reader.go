package gomat

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type MatFile struct {
	Header *MatHeader
	r      io.Reader
	w      io.Writer
}

const (
	headerLen                = 128
	headerTextLen            = 116
	headerSubsystemOffsetLen = 8
	headerFlagLen            = 4
)

type MatHeader struct {
	Level    string
	Platform string
	Created  time.Time
	Endian   binary.ByteOrder
}

func (mh *MatHeader) String() string {
	return fmt.Sprintf("MATLAB %s MAT-file, Platform: %s, Create on: %s", mh.Level, mh.Platform, mh.Created.Format(time.ANSIC))
}

func createBufioReader(r io.Reader, size int) (*bufio.Reader, error) {
	buf := make([]byte, size)
	n, err := r.Read(buf)
	if n == 0 {
		// error check
		return nil, errors.New("0 bytes read")
	}
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(bytes.NewBuffer(buf)), nil
}

func (mh *MatHeader) readHeader(r io.Reader) error {
	// read text field
	if err := mh.readTextField(r); err != nil {
		return err
	}
	// read subsystem data offset field
	if err := mh.readSubsystemDataOffsetField(r); err != nil {
		return err
	}
	// read flag fields
	if err := mh.readFlagFields(r); err != nil {
		return err
	}
	return nil
}

func (mh *MatHeader) readTextField(rd io.Reader) error {
	r, err := createBufioReader(rd, headerTextLen)
	if err != nil {
		return err
	}
	// Every mat file starts with "MATLAB"
	if prefix, err := r.ReadBytes(' '); err != nil {
		return err
	} else if !bytes.Equal(prefix, []byte("MATLAB ")) {
		// error must be declared as constant
		return errors.New("not a valid .mat file")
	}

	// Read Level
	// This must be 5.0
	level, err := r.ReadString(' ')
	if err != nil {
		return err
	}
	mh.Level = strings.TrimSpace(level)
	if mh.Level != "5.0" {
		// error must be declared as constant
		return errors.New("not level 5 files")
	}

	// Read Platform
	if _, err := r.Discard(len("MAT-file, Platform: ")); err != nil {
		return err
	}
	platform, err := r.ReadString(' ')
	if err != nil {
		return err
	}
	mh.Platform = strings.TrimRight(platform, ", ")

	// Read Created Day
	if _, err := r.Discard(len("Created on: ")); err != nil {
		return err
	}
	date := make([]byte, 24)
	if _, err = r.Read(date); err != nil {
		return err
	}
	if mh.Created, err = time.Parse(time.ANSIC, strings.TrimSpace(string(date))); err != nil {
		return err
	}

	return nil
}

func (mh *MatHeader) readSubsystemDataOffsetField(r io.Reader) error {
	return nil
}

func (mh *MatHeader) readFlagFields(r io.Reader) error {
	return nil
}

func NewMatFileFromReader(r io.Reader) (*MatFile, error) {
	f := &MatFile{r: r}
	err := f.readHeader()
	return f, err
}

func (mf *MatFile) readHeader() error {
	h := &MatHeader{}
	h.readHeader(mf.r)
	mf.Header = h

	return nil
}
