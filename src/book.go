package shelf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
)

const (
	ShelfRegexp   = `^(?P<id>\w{24})_(?P<base>.*)\.pdf$`
	IDFormat      = "20060102T150405"
	DocExtension  = ".pdf"
	MetaExtension = ".toml"
	MetaPath      = "index.toml"
)

var (
	ErrInputDir    = fmt.Errorf("ファイルを入力を想定している箇所でディレクトリが入力された")
	ErrInputNotPDF = fmt.Errorf("PDF以外が入力された")
)

// ファイル名の一部となるID
// 読み書きで使う
type BookID string

func NewBookID(t time.Time) BookID {
	nano := t.Nanosecond()
	formatted := fmt.Sprintf("%s%09d", t.Format(IDFormat), nano)

	return BookID(formatted)
}

// shelfフォーマットを満たす元ファイル(PDFなど)
// "{id}_{base}.{ext}"
// baseは人が見て識別する用で、プログラム側からは参照しない
type Book struct {
	File os.File
}

// shelfフォーマットを満たすファイルなのを確認して初期化する
func NewBook(file os.File) *Book {
	book := Book{File: file}

	return &book
}

// IDフォーマットを満たすかを判定する
func ValidFormat(file os.File) bool {
	if filepath.Ext(file.Name()) != DocExtension {
		return false
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return false
	}
	if fileinfo.IsDir() {
		return false
	}
	_, err = executeShelfRegexp(fileinfo.Name())
	if err != nil {
		return false
	}

	return true
}

// ファイルのIDを取得する
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

// ファイルのフルパスを取得する
func (b *Book) GetFullPath() string {
	return b.File.Name()
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

// メタファイルにデフォルトのメタ情報を書き込む
func (b *Book) writeBlankMetaFile(w io.Writer) error {
	title := "PDF Title"
	extract, err := b.ExtractPDFTitle()
	if err != nil {
		id, _ := b.GetID()
		log.Printf("PDFタイトルを取得できなかった。id: %s\n", id)
	}
	if extract != "" {
		title = extract
	}

	meta := Meta{
		Title: GetPtr(title),
		TODO:  GetPtr(TODOTypeNONE),
		Tags:  GetPtr([]string{"new"}),
	}
	err = toml.NewEncoder(w).Encode(meta)
	if err != nil {
		return err
	}

	return nil
}

func (b *Book) AppendMeta(metas Metas) error {
	title := "PDF Title"
	extract, err := b.ExtractPDFTitle()
	if err != nil {
		id, _ := b.GetID()
		log.Printf("PDFタイトルを取得できなかった。id: %s\n", id)
	}
	if extract != "" {
		title = extract
	}

	id, err := b.GetID()
	if err != nil {
		return err
	}

	if _, ok := metas[string(id)]; !ok {
		metas[string(id)] = Meta{
			Title: GetPtr(title),
			TODO:  GetPtr(TODOTypeNONE),
			Tags:  GetPtr([]string{"new"}),
		}
	}

	return nil
}

func (b *Book) ExtractImageBase64() (string, error) {
	// PDFファイルを解析
	pdfReader, err := model.NewPdfReader(&b.File)
	if err != nil {
		return "", err
	}

	// ページを抽出
	pageNum := 1
	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		return "", err
	}

	device := render.NewImageDevice()
	image, err := device.Render(page)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, image); err != nil {
		return "", err
	}
	str := buf.Bytes()
	imgBase64Str := base64.StdEncoding.EncodeToString(str)

	return imgBase64Str, nil
}

// shelf対応のパスに変換して返す
func generateShelfPath(file os.File, time time.Time) string {
	id := NewBookID(time)
	fileName := filepath.Base(file.Name())
	base := fileName[:len(fileName)-len(filepath.Ext(fileName))] // 拡張子を除去する
	shelfFileName := fmt.Sprintf("%s_%s.pdf", id, base)
	dir := filepath.Join(filepath.Dir(file.Name()))
	shelfpath := filepath.Join(dir, shelfFileName)

	return shelfpath
}

// IDを返す
func executeShelfRegexp(raw string) (BookID, error) {
	re, err := regexp.Compile(ShelfRegexp)
	if err != nil {
		return "", err
	}
	matches := re.FindAllStringSubmatch(filepath.Base(raw), -1)
	if len(matches) < 1 {
		return "", fmt.Errorf("shelfフォーマットにマッチしなかった: %s", raw)
	}

	return BookID(matches[0][re.SubexpIndex("id")]), nil
}
