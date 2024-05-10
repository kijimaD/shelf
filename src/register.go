package shelf

import (
	"fmt"
	"os"
	"time"
)

var (
	ErrAlreadyFormatted = fmt.Errorf("すでにフォーマットを満たしている")
)

// 入力パスに対して、ID付きファイル名にリネーム
func Register(originfile *os.File) (*Book, error) {
	path := ""

	// すでにフォーマットを満たしていたら(=エラーが出なければ)終了
	ok := ValidFormat(*originfile)
	if ok {
		// すでにフォーマットを満たす
		path = originfile.Name()
	} else {
		// フォーマットを満たさない
		path = generateShelfPath(*originfile, time.Now())
		err := os.Rename(originfile.Name(), path)
		if err != nil {
			return nil, err
		}
	}
	shelfFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if shelfFile == nil {
		return nil, fmt.Errorf("新しいshelfFileが発見できなかった")
	}
	book := NewBook(*shelfFile)

	return book, nil
}
