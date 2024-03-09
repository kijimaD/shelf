package shelf

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBookID(t *testing.T) {
	assert.Equal(t, BookID("20011111T111111"), NewBookID(time.Date(2001, 11, 11, 11, 11, 11, 0, time.UTC)))
	assert.Equal(t, BookID("20010203T040506"), NewBookID(time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)))
	// 桁埋めされている
	assert.Equal(t, BookID("20010101T010101"), NewBookID(time.Date(2001, 1, 1, 1, 1, 1, 0, time.UTC)))
}

func TestExecuteShelfRegexp(t *testing.T) {
	_, err := executeShelfRegexp("20010203T040506_hello.pdf")
	assert.NoError(t, err)
}

func TestNewBook(t *testing.T) {
	t.Run("取得できる", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		_, err = NewBook(*tempfile)
		assert.NoError(t, err)
	})
	t.Run("無効なパスだとエラーを返す", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "INVALID_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		_, err = NewBook(*tempfile)
		assert.Error(t, err)
	})
}

func TestGetID(t *testing.T) {
	t.Run("IDを取得できる", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		id, err := b.GetID()
		assert.NoError(t, err)
		assert.Equal(t, BookID("20010101T010101"), id)
	})
}

func TestGetFullPath(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		fullpath, err := b.GetFullPath()
		assert.Contains(t, fullpath, "/tmp/20010101T010101_")
		assert.Contains(t, fullpath, ".pdf")
	})
}

func TestMetaPath(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		metapath, err := b.MetaPath()
		assert.NoError(t, err)
		assert.Equal(t, "/tmp/20010101T010101.toml", metapath)
	})
}
