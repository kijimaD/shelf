package src

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const IDFormat = "20060102T150405"

func formatID(t time.Time) string {
	return t.Format(IDFormat)
}

// "{id}_{base}.{ext}"
type Fullname struct {
	id   string // ID
	base string // 任意の文字列
	ext  string // 拡張子
}

// フル文字列からFullnameを生成する
func NewFullname(full string) (*Fullname, error) {
	ext := filepath.Ext(full)
	if ext == "" {
		return nil, fmt.Errorf("拡張子がない")
	}

	pattern := `^(?P<id>\w{15})_(?P<base>.*)\.\w*$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	matches := re.FindAllStringSubmatch(full, -1)
	if len(matches) < 1 {
		return nil, fmt.Errorf("マッチしなかった")
	}
	id := matches[0][re.SubexpIndex("id")]
	base := matches[0][re.SubexpIndex("base")]

	return &Fullname{id: id, base: base, ext: ext}, nil
}

// IDなしファイル名からFullnameを生成する
func NewFullnameByRaw(origin string, t time.Time) (*Fullname, error) {
	// 拡張子がなければエラー
	ext := filepath.Ext(origin)
	if ext == "" {
		return nil, fmt.Errorf("拡張子がない")
	}

	id := formatID(t)
	base := strings.TrimSuffix(origin, ext)

	return &Fullname{id: id, base: base, ext: ext}, nil
}

func (f *Fullname) String() string {
	return fmt.Sprintf("%s_%s%s", f.id, f.base, f.ext)
}

func (f *Fullname) MetaFilename() string {
	return fmt.Sprintf("%s_%s.toml", f.id, f.base)
}

func (f *Fullname) writeMetafile(book Book, w io.Writer) error {
	err := toml.NewEncoder(w).Encode(book)
	if err != nil {
		return err
	}

	return nil
}

func (f *Fullname) touchMetafile(w io.Writer) error {
	book := Book{
		Title: f.base,
		TODO:  TODOTypeNONE,
		Tags:  []string{"new"},
	}
	err := f.writeMetafile(book, w)
	if err != nil {
		return err
	}

	return nil
}
