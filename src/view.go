// ファイルとメタファイルを収集する
package shelf

import (
	"log"
	"os"
	"path/filepath"
	"slices"
)

type View struct {
	FilePath string
	Meta     Meta
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
		view := View{
			FilePath: book.GetFullPath(),
			Meta:     *meta,
		}
		views = append(views, view)
	}

	return views
}

func FilterViewsByTag(tag string, views []View) []View {
	newviews := []View{}

	for _, view := range views {
		if slices.Contains(*view.Meta.Tags, tag) {
			newviews = append(newviews, view)
		}
	}

	return newviews
}

func UniqTags(views []View) []string {
	uniq := []string{}
	m := make(map[string]bool)

	for _, view := range views {
		for _, tag := range *view.Meta.Tags {
			if !m[tag] {
				m[tag] = true
				uniq = append(uniq, tag)
			}
		}

	}

	return uniq
}
