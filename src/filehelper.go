package src

import (
	"fmt"
	"io"
	"os"
	"path"
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

// TODO: この名前は実態を全く表せてないのでどうにかする
// "{id}_{base}.{ext}"
type Fullname struct {
	id   string // ID
	dir  string // ディレクトリ名
	base string // ファイル名(ディレクトリ名を含まない)
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
	fullname := Fullname{
		id:   matches[0][re.SubexpIndex("id")],
		dir:  filepath.Dir(full),
		base: matches[0][re.SubexpIndex("base")],
		ext:  ext,
	}

	return &fullname, nil
}

// IDなしファイル名からFullnameを生成する
func NewFullnameByRaw(origin string, t time.Time) (*Fullname, error) {
	// 拡張子がなければエラー
	ext := filepath.Ext(origin)
	if ext == "" {
		return nil, fmt.Errorf("拡張子がない")
	}
	fullname := Fullname{
		id:   formatID(t),
		dir:  filepath.Dir(origin),
		base: strings.TrimSuffix(path.Base(origin), ext),
		ext:  ext,
	}

	return &fullname, nil
}

// ファイル名。ディレクトリ名は含まない
// TODO: 名前をどうにかする
func (f *Fullname) String() string {
	return fmt.Sprintf("%s_%s%s", f.id, f.base, f.ext)
}

func (f *Fullname) MetaFilename() string {
	return fmt.Sprintf("%s_%s.toml", f.id, f.base)
}

// 元のパスを返す
func (f *Fullname) OriginalFilename() string {
	return filepath.Join(f.dir, f.base+f.ext)
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

// 同じ階層でIDつきパスにリネームする
func (f *Fullname) rename() (string, error) {
	newpath := filepath.Join(f.dir, f.String())
	err := os.Rename(f.OriginalFilename(), newpath)
	if err != nil {
		return "", err
	}

	return newpath, nil
}
