package shelf

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const ShelfRegexp = `^(?P<id>\w{15})_(?P<base>.*)\.\w*$`

// ファイル名の一部となるID
// 読み書きで使う
type BookID string

func NewBookID(t time.Time) BookID {
	return BookID(t.Format(IDFormat))
}

// shelfフォーマットを満たすファイル
// "{id}_{base}.{ext}"
type Book struct {
	File os.File
	Meta Meta
}

func NewBook(file os.File) (*Book, error) {
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	_, err = executeShelfRegexp(fileinfo.Name())
	if err != nil {
		return nil, err
	}

	// TODO: 用意する
	book := Book{
		File: file,
		Meta: Meta{Title: "example", TODO: TODOTypeNONE, Tags: []string{}},
	}

	return &book, nil
}

func (b Book) GetID() (BookID, error) {
	fileinfo, err := b.File.Stat()
	if err != nil {
		return "", err
	}

	id, err := executeShelfRegexp(fileinfo.Name())
	if err != nil {
		return "", err
	}

	return BookID(id), nil
}

func (b Book) GetFullPath() (string, error) {
	// このやり方は正しくなさそう
	// UnixライクOS限定
	fd := int(b.File.Fd())
	fullPath, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", fd))
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

func (b Book) MetaPath() (string, error) {
	fullpath, err := b.GetFullPath()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(fullpath)
	id, _ := b.GetID()
	filename := fmt.Sprintf("%s.toml", id)

	return filepath.Join(dir, filename), nil
}

func executeShelfRegexp(raw string) (string, error) {
	re, err := regexp.Compile(ShelfRegexp)
	if err != nil {
		return "", err
	}
	matches := re.FindAllStringSubmatch(filepath.Base(raw), -1)
	if len(matches) < 1 {
		return "", fmt.Errorf("shelfフォーマットにマッチしなかった")
	}

	return matches[0][re.SubexpIndex("id")], nil
}
