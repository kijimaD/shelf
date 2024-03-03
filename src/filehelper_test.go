package src

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatID(t *testing.T) {
	assert.Equal(t, "20011111T111111", formatID(time.Date(2001, 11, 11, 11, 11, 11, 0, time.UTC)))
	assert.Equal(t, "20010203T040506", formatID(time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)))
	// 桁埋めされている
	assert.Equal(t, "20010101T010101", formatID(time.Date(2001, 1, 1, 1, 1, 1, 0, time.UTC)))
}

func TestFullname(t *testing.T) {
	date := time.Date(2001, 1, 1, 1, 1, 1, 0, time.UTC)
	{
		result, err := NewFullnameByRaw("aaa.pdf", date)
		assert.NoError(t, err)
		assert.Equal(t, "20010101T010101_aaa.pdf", result.String())
	}
	{
		result, err := NewFullnameByRaw("aaa.epub", date)
		assert.NoError(t, err)
		assert.Equal(t, "20010101T010101_aaa.epub", result.String())
	}
	{
		result, err := NewFullnameByRaw("/any/dir/aaa.epub", date)
		assert.NoError(t, err)
		assert.Equal(t, "20010101T010101_aaa.epub", result.String())
	}
	{
		result, err := NewFullnameByRaw("あああ.epub", date)
		assert.NoError(t, err)
		assert.Equal(t, "20010101T010101_あああ.epub", result.String())
	}
	{
		_, err := NewFullnameByRaw("aaa", date)
		assert.Error(t, err)
	}
}

func TestLoadFullname(t *testing.T) {
	{
		_, err := NewFullname("20010101T010101_aaa.pdf")
		assert.NoError(t, err)
	}
	{
		_, err := NewFullname("20010101T010101_あああ.pdf")
		assert.NoError(t, err)
	}
	{
		_, err := NewFullname("20010101T010101_aaa")
		assert.Error(t, err)
	}
	{
		_, err := NewFullname("INVALID_20010101T010101_aaa.pdf")
		assert.Error(t, err)
	}
}

func TestMetaFilename(t *testing.T) {
	{
		fullname, err := NewFullname("20010101T010101_aaa.pdf")
		assert.NoError(t, err)
		result := fullname.MetaFilename()
		assert.Equal(t, "20010101T010101.toml", result)
	}
}

func TestOriginalFilename(t *testing.T) {
	date := time.Date(2001, 1, 1, 1, 1, 1, 0, time.UTC)
	{
		fullname, err := NewFullnameByRaw("aaa.pdf", date)
		assert.NoError(t, err)
		result := fullname.OriginalFilename()
		assert.Equal(t, "aaa.pdf", result)
	}
	{
		fullname, err := NewFullnameByRaw("aaa/bbb.pdf", date)
		assert.NoError(t, err)
		result := fullname.OriginalFilename()
		assert.Equal(t, "aaa/bbb.pdf", result)
	}
	{
		fullname, err := NewFullnameByRaw("aaa/bbb.pdf", date)
		assert.NoError(t, err)
		result := fullname.OriginalFilename()
		assert.Equal(t, "aaa/bbb.pdf", result)
	}
	{
		fullname, err := NewFullnameByRaw("/aaa/bbb.pdf", date)
		assert.NoError(t, err)
		result := fullname.OriginalFilename()
		assert.Equal(t, "/aaa/bbb.pdf", result)
	}
}

func TestWriteMetafile(t *testing.T) {
	fullname, err := NewFullname("20010101T010101_aaa.pdf")
	assert.NoError(t, err)

	book := Book{
		Title: fullname.base,
		TODO:  TODOTypeNONE,
		Tags:  []string{"new"},
	}
	buf := bytes.Buffer{}
	assert.NoError(t, fullname.writeMetafile(book, &buf))
	expect := `title = "aaa"
todo = "NONE"
tags = ["new"]
`
	assert.Equal(t, expect, buf.String())
}

func TestBlankMetafile(t *testing.T) {
	fullname, err := NewFullname("20010101T010101_abc.pdf")
	assert.NoError(t, err)

	buf := bytes.Buffer{}
	assert.NoError(t, fullname.blankMetafile(&buf))
	expect := `title = "abc"
todo = "NONE"
tags = ["new"]
`
	assert.Equal(t, expect, buf.String())
}

func TestRename(t *testing.T) {
	tmpfile, err := os.CreateTemp(os.TempDir(), "test-*.pdf")
	assert.NoError(t, err)

	date := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
	{
		fullname, err := NewFullnameByRaw(tmpfile.Name(), date)
		assert.NoError(t, err)
		newpath, err := fullname.rename()
		assert.NoError(t, err)

		// 存在確認
		_, err = os.Stat(newpath)
		assert.NoError(t, err)

		// 非存在確認
		_, err = os.Stat(tmpfile.Name())
		assert.Error(t, err)
	}
}
