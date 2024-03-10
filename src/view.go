// ファイルとメタファイルを収集する
package shelf

import (
	"log"
	"os"
	"path/filepath"
)

type View struct {
	FilePath  string
	Meta      Meta
	Thumbnail string // base64エンコード
}

// ディレクトリから、メタデータを収集する
func GenerateViews(dirpath string) []View {
	views := []View{}
	files, _ := os.ReadDir(dirpath)

	for _, file := range files {
		if filepath.Ext(file.Name()) != DocExtension {
			continue
		}

		f, err := os.Open(filepath.Join(dirpath, file.Name()))
		if err != nil {
			log.Println(err)
			continue
		}
		if f == nil {
			log.Println("ファイルがnilだった")
			continue
		}
		book, err := NewBook(*f)
		if err != nil {
			log.Println(err)
			continue
		}
		meta, err := book.GetMetaData()
		if err != nil {
			log.Println(err)
			continue
		}
		if meta == nil {
			log.Println("メタ情報がなかった")
			continue
		}
		base64Str, err := book.ExtractImageBase64()
		if err != nil {
			log.Println(err)
			continue
		}

		view := View{
			FilePath:  book.GetFullPath(),
			Meta:      *meta,
			Thumbnail: base64Str,
		}
		views = append(views, view)
	}

	return views
}
