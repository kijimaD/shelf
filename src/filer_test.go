package shelf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestGenerateViews(t *testing.T) {
	// 残らなくて便利なので他のテストでもt.TempDir()を使っていきたい
	tempdir := t.TempDir()
	{
		fname := filepath.Join(tempdir, "20010101T010101_abc.pdf")
		assert.NoError(t, os.WriteFile(fname, []byte(""), 0666))
	}
	{
		fname := filepath.Join(tempdir, "20010101T010101.toml")
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
		fname := filepath.Join(tempdir, "20010101T010102_def.pdf")
		assert.NoError(t, os.WriteFile(fname, []byte(""), 0666))

	}
	{
		fname := filepath.Join(tempdir, "20010101T010102.toml")
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
