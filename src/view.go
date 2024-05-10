// ファイルとメタファイルを収集する
package shelf

import (
	"log"
	"os"
	"path/filepath"
	"slices"
)

type View struct {
	ID       BookID
	FilePath string
	Meta     Meta
}

// メタデータを収集する
func GenerateViews(dirpath string, metas Metas) []View {
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
		book := NewBook(*f)

		id, err := book.GetID()
		if err != nil {
			log.Fatal(err)
		}
		view := View{
			ID:       id,
			FilePath: book.GetFullPath(),
			Meta:     metas[string(id)],
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

// |    discard     | return |
// |<------->|cursor|<------>|
func SkipCursor(id string, views []View) []View {
	ids := []string{}

	for _, view := range views {
		ids = append(ids, string(view.ID))
		if len(views) > 0 {
			// 先頭を削除
			views = views[1:]
		}
		if BookID(id) == view.ID {
			break
		}
	}

	return views
}

func UniqTags(views []View) []string {
	uniq := []string{}
	m := make(map[string]bool)

	for _, view := range views {
		for _, tag := range *view.Meta.Tags {
			if tag == "" {
				continue
			}
			if !m[tag] {
				m[tag] = true
				uniq = append(uniq, tag)
			}
		}

	}

	return uniq
}
