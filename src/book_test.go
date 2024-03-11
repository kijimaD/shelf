package shelf

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestNewBookID(t *testing.T) {
	assert.Equal(t, BookID("20011111T111111123456789"), NewBookID(time.Date(2001, 11, 11, 11, 11, 11, 123456789, time.UTC)))
	assert.Equal(t, BookID("20010203T040506123456789"), NewBookID(time.Date(2001, 2, 3, 4, 5, 6, 123456789, time.UTC)))
	// 桁埋めされている
	assert.Equal(t, BookID("20010101T010101123456789"), NewBookID(time.Date(2001, 1, 1, 1, 1, 1, 123456789, time.UTC)))
}

func TestExecuteShelfRegexp(t *testing.T) {
	_, err := executeShelfRegexp("20010203T040506000000000_hello.pdf")
	assert.NoError(t, err)
}

func TestNewBook(t *testing.T) {
	t.Run("取得できる", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000_*.pdf")
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
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		id, err := b.GetID()
		assert.NoError(t, err)
		assert.Equal(t, BookID("20010101T010101000000000"), id)
	})
}

func TestGetFullPath(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		fullpath := b.GetFullPath()
		assert.Contains(t, fullpath, "/tmp/20010101T010101000000000_")
		assert.Contains(t, fullpath, ".pdf")
	})
}

func TestMetaPath(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		metapath, err := b.MetaPath()
		assert.NoError(t, err)
		assert.Equal(t, "/tmp/20010101T010101000000000.toml", metapath)
	})
}

func TestLoadMetaData(t *testing.T) {
	t.Run("メタデータを読み込める", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		metafile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000.toml")
		assert.NoError(t, err)
		defer os.Remove(metafile.Name())
		// CreateTempで作るファイルには末尾に番号がつくので、リネーム
		newPath := filepath.Join(filepath.Dir(metafile.Name()), "20010101T010101000000000.toml")
		assert.NoError(t, os.Rename(metafile.Name(), newPath))
		defer os.Remove(newPath)

		meta := Meta{
			Title: GetPtr("hello"),
			TODO:  GetPtr(TODOTypeNONE),
			Tags:  GetPtr([]string{"new"}),
		}
		assert.NoError(t, toml.NewEncoder(metafile).Encode(meta))

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		loadedMeta, err := b.GetMetaData()
		assert.NoError(t, err)
		assert.Equal(t, meta, *loadedMeta)
	})
	t.Run("メタファイルが空の場合エラーを返す", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		metafile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000.toml")
		assert.NoError(t, err)
		defer os.Remove(metafile.Name())
		// CreateTempで作るファイルには末尾に番号がつくので、リネーム
		newPath := filepath.Join(filepath.Dir(metafile.Name()), "20010101T010101000000000.toml")
		assert.NoError(t, os.Rename(metafile.Name(), newPath))
		defer os.Remove(newPath)

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		_, err = b.GetMetaData()
		assert.Error(t, err)
	})
	t.Run("メタファイルがないときはエラーを返す", func(t *testing.T) {
		tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000_*.pdf")
		assert.NoError(t, err)
		defer os.Remove(tempfile.Name())

		b, err := NewBook(*tempfile)
		assert.NoError(t, err)
		_, err = b.GetMetaData()
		assert.Error(t, err)
	})
}

func TestExtractPDFTitle(t *testing.T) {
	tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000_*.pdf")
	assert.NoError(t, err)
	defer os.Remove(tempfile.Name())

	srcfile, err := os.Open("../fixture/example.pdf")
	assert.NoError(t, err)
	defer srcfile.Close()
	content, err := io.ReadAll(srcfile)
	assert.NoError(t, err)
	_, err = tempfile.Write(content)
	assert.NoError(t, err)

	b, err := NewBook(*tempfile)
	assert.NoError(t, err)

	_, err = b.ExtractPDFTitle()
	assert.NoError(t, err)
}

func TestGenerateShelfPath(t *testing.T) {
	tempfile, err := os.CreateTemp(os.TempDir(), "test_*.pdf")
	assert.NoError(t, err)
	defer os.Remove(tempfile.Name())

	path := generateShelfPath(*tempfile, time.Date(2001, 11, 11, 11, 11, 11, 0, time.UTC))
	assert.Contains(t, path, "/tmp/20011111T111111000000000_test")
	assert.Contains(t, path, ".pdf")
}

func TestWriteBlankMetaFile(t *testing.T) {
	tempfile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000_*.pdf")
	assert.NoError(t, err)
	defer os.Remove(tempfile.Name())

	// PDFタイトルを取得する関係で、中身を実際のPDFにしておく
	srcfile, err := os.Open("../fixture/example.pdf")
	assert.NoError(t, err)
	defer srcfile.Close()
	content, err := io.ReadAll(srcfile)
	assert.NoError(t, err)
	_, err = tempfile.Write(content)
	assert.NoError(t, err)

	// 中身なしのメタファイル
	metafile, err := os.CreateTemp(os.TempDir(), "20010101T010101000000000.toml")
	assert.NoError(t, err)
	defer os.Remove(metafile.Name())

	b, err := NewBook(*tempfile)
	assert.NoError(t, err)
	assert.NoError(t, b.writeBlankMetaFile(metafile))

	// 再度開いて中身を確認する
	metafile, err = os.Open(metafile.Name())
	assert.NoError(t, err)
	metacontent, err := io.ReadAll(metafile)
	assert.NoError(t, err)
	expect := `title = "example"
todo = "NONE"
tags = ["new"]
`
	assert.Equal(t, expect, string(metacontent))
}
