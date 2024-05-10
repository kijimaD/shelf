package shelf

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	tempdir := t.TempDir()
	fname := filepath.Join(tempdir, "test.pdf")
	tempfile, err := os.Create(fname)
	assert.NoError(t, err)
	defer os.Remove(tempfile.Name())

	samplefile, err := os.Open("../fixture/example.pdf")
	assert.NoError(t, err)
	defer samplefile.Close()
	content, err := io.ReadAll(samplefile)
	assert.NoError(t, err)
	_, err = tempfile.Write(content)
	assert.NoError(t, err)

	book, err := Register(tempfile)
	assert.NoError(t, err)

	// 存在確認
	_, err = os.Stat(book.GetFullPath())
	assert.NoError(t, err)
	// 非存在確認
	_, err = os.Stat(tempfile.Name())
	assert.Error(t, err)
}
