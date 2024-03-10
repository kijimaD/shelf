package shelf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
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

func TestLoadMetaData(t *testing.T) {
	t.Run("メタデータを読み込める", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		metafile, err := os.CreateTemp(os.TempDir(), "20010101T010101.toml")
		assert.NoError(t, err)
		defer os.Remove(metafile.Name())
		// CreateTempで作るファイルには末尾に番号がつくので、リネーム
		newPath := filepath.Join(filepath.Dir(metafile.Name()), "20010101T010101.toml")
		assert.NoError(t, os.Rename(metafile.Name(), newPath))
		defer os.Remove(newPath)

		meta := Meta{
			Title: "hello",
			TODO:  TODOTypeNONE,
			Tags:  []string{"new"},
		}
		assert.NoError(t, toml.NewEncoder(metafile).Encode(meta))

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		assert.NoError(t, b.LoadMetaData())
		assert.Equal(t, meta, b.Meta) // 読み込めている
	})
	t.Run("メタファイルがないときはエラーを返す", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		assert.Error(t, b.LoadMetaData())
	})
}

func TestExtractPDFTitle(t *testing.T) {
	tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101_*.pdf")
	assert.NoError(t, err)
	defer os.Remove(tempfile.Name())

	srcfile, err := os.Open("../example.pdf")
	assert.NoError(t, err)
	defer srcfile.Close()
	content, err := ioutil.ReadAll(srcfile)
	assert.NoError(t, err)

	_, err = tempfile.Write(content)
	assert.NoError(t, err)

	b, err := NewBook(*tempfile)
	assert.NoError(t, err)

	_, err = b.ExtractPDFTitle()
	assert.NoError(t, err)
}
