package src

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatID(t *testing.T) {
	assert.Equal(t, "20011111T111111", formatID(time.Date(2001, 11, 11, 11, 11, 11, 0, time.UTC)))
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
		assert.Equal(t, "20010101T010101_aaa.toml", result)
	}
}

func TestMakeMetafile(t *testing.T) {
	fullname, err := NewFullname("20010101T010101_aaa.pdf")
	assert.NoError(t, err)
	assert.NoError(t, fullname.makeMetafile())
}
