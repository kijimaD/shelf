// ファイルとメタファイルを収集する
package shelf

import (
	"fmt"
	"io/ioutil"
)

type Filer struct {
	FilePath string
	Meta     Meta
}

// ディレクトリから、メタデータを収集する
func GenerateFiler(dirpath string) []Filer {
	filers := []Filer{}
	files, _ := ioutil.ReadDir(dirpath)

	for _, f := range files {
		// IDのついたパターンにマッチするものを拾う
		// 本体とメタファイルが存在するものを拾う
		// IDの一覧を取得して、両方アクセスする
		pathobj, err := NewFullname(f.Name())
		if err != nil {
			continue
		}
		if pathobj.ExistMetaFile() {
			fmt.Println(dirpath)
			fmt.Println("存在する")
		}

	}

	return filers
}
