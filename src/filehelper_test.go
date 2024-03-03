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
		result, err := GenFullname("aaa.pdf", date)
		assert.NoError(t, err)
		assert.Equal(t, "20010101T010101_aaa.pdf", string(result))
	}
	{
		result, err := GenFullname("aaa.epub", date)
		assert.NoError(t, err)
		assert.Equal(t, "20010101T010101_aaa.epub", string(result))
	}
	{
		result, err := GenFullname("あああ.epub", date)
		assert.NoError(t, err)
		assert.Equal(t, "20010101T010101_あああ.epub", string(result))
	}
	{
		_, err := GenFullname("aaa", date)
		assert.Error(t, err)
	}
}

func TestLoadFullname(t *testing.T) {
	_, err := LoadFullname("20010101T010101_aaa.pdf")
	assert.NoError(t, err)
}
