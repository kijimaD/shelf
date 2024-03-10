package shelf

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestGenerateViews(t *testing.T) {
	// 残らなくて便利なので他のテストでもt.TempDir()を使っていきたい
	tempdir := t.TempDir()

	examplefile, err := os.Open("../fixture/example.pdf")
	assert.NoError(t, err)
	defer examplefile.Close()
	example, err := io.ReadAll(examplefile)
	assert.NoError(t, err)
	{
		fname := filepath.Join(tempdir, "20010101T010101123456789_abc.pdf")
		f, err := os.Create(fname)
		assert.NoError(t, err)
		_, err = f.Write(example)
	}
	{
		fname := filepath.Join(tempdir, "20010101T010101123456789.toml")
		f, err := os.Create(fname)
		assert.NoError(t, err)
		meta := Meta{
			Title: GetPtr("hello1"),
			TODO:  GetPtr(TODOTypeNONE),
			Tags:  GetPtr([]string{"new"}),
		}
		assert.NoError(t, toml.NewEncoder(f).Encode(meta))
	}
	{
		fname := filepath.Join(tempdir, "20010101T010102123456789_def.pdf")
		f, err := os.Create(fname)
		assert.NoError(t, err)
		_, err = f.Write(example)
	}
	{
		fname := filepath.Join(tempdir, "20010101T010102123456789.toml")
		f, err := os.Create(fname)
		assert.NoError(t, err)
		meta := Meta{
			Title: GetPtr("hello2"),
			TODO:  GetPtr(TODOTypeNONE),
			Tags:  GetPtr([]string{"new"}),
		}
		assert.NoError(t, toml.NewEncoder(f).Encode(meta))
	}

	views := GenerateViews(tempdir)
	assert.Equal(t, 2, len(views))
	assert.Equal(t, "hello1", *views[0].Meta.Title)
	assert.Equal(t, "hello2", *views[1].Meta.Title)
}
