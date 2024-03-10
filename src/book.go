package shelf

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/unidoc/unipdf/v3/model"
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
// baseは人が見て識別する用で、プログラム側からは参照しない
type Book struct {
	File os.File
	Meta Meta
}

func NewBook(file os.File) (*Book, error) {
	book := Book{File: file}

	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	_, err = executeShelfRegexp(fileinfo.Name())
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (b *Book) LoadMetaData() error {
	metapath, err := b.MetaPath()
	if err != nil {
		return err
	}
	metafile, err := os.Open(metapath)
	if err != nil {
		return err
	}
	defer metafile.Close()

	meta := Meta{}
	_, err = toml.NewDecoder(metafile).Decode(&meta)
	if err != nil {
		return err
	}
	b.Meta = meta

	return nil
}

func (b *Book) GetID() (BookID, error) {
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

func (b *Book) GetFullPath() (string, error) {
	// このやり方は正しくなさそう
	// UnixライクOS限定
	fd := int(b.File.Fd())
	fullPath, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", fd))
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

func (b *Book) MetaPath() (string, error) {
	fullpath, err := b.GetFullPath()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(fullpath)
	id, _ := b.GetID()
	filename := fmt.Sprintf("%s.toml", id)

	return filepath.Join(dir, filename), nil
}

func (b *Book) ExtractPDFTitle() (string, error) {
	pdfReader, err := model.NewPdfReader(&b.File)
	if err != nil {
		return "", err
	}

	pdfInfo, err := pdfReader.GetPdfInfo()
	if err != nil {
		return "", err
	}
	if pdfInfo.Title != nil {
		return pdfInfo.Title.Decoded(), nil
	}

	return "", nil
}

// 初期生成用にブランクのメタファイルを生成して書き込む
func (b *Book) writeBlankMetaFile(w io.Writer) error {
	title := "PDF Title"
	extract, err := b.ExtractPDFTitle()
	if err != nil {
		log.Println("PDFタイトルを取得できなかった")
	}
	if extract != "" {
		title = extract
	}

	meta := Meta{
		Title: title,
		TODO:  TODOTypeNONE,
		Tags:  []string{"new"},
	}
	err = toml.NewEncoder(w).Encode(meta)
	if err != nil {
		return err
	}

	return nil
}

// IDを返す
func executeShelfRegexp(raw string) (string, error) {
	re, err := regexp.Compile(ShelfRegexp)
	if err != nil {
		return "", err
	}
	matches := re.FindAllStringSubmatch(filepath.Base(raw), -1)
	if len(matches) < 1 {
		return "", fmt.Errorf("shelfフォーマットにマッチしなかった: %s", raw)
	}

	return matches[0][re.SubexpIndex("id")], nil
}
