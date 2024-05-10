package shelf

import (
	"io"
	"os"
	"path/filepath"
	"testing"

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
		fname := filepath.Join(tempdir, "20010101T010102123456789_def.pdf")
		f, err := os.Create(fname)
		assert.NoError(t, err)
		_, err = f.Write(example)
	}

	views := GenerateViews(tempdir, Metas{
		"20010101T010101123456789": Meta{
			Title: GetPtr("hello1"),
			Tags:  GetPtr([]string{"one"}),
		},
		"20010101T010102123456789": Meta{
			Title: GetPtr("hello2"),
			Tags:  GetPtr([]string{"two"}),
		},
	})
	assert.Equal(t, 2, len(views))
	assert.Equal(t, "hello1", *views[0].Meta.Title)
	assert.Equal(t, "hello2", *views[1].Meta.Title)

	{
		searchOne := FilterViewsByTag("one", views)
		assert.Equal(t, 1, len(searchOne))
		assert.Equal(t, "hello1", *searchOne[0].Meta.Title)
	}
	{
		searchTwo := FilterViewsByTag("two", views)
		assert.Equal(t, 1, len(searchTwo))
		assert.Equal(t, "hello2", *searchTwo[0].Meta.Title)
	}
	{
		search := FilterViewsByTag("NOT FOUND", views)
		assert.Equal(t, 0, len(search))
	}
}
