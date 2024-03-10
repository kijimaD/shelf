package shelf

import (
	"fmt"
	"os"
	"time"
)

// 入力パスに対して、ID付きファイル名にリネーム + メタファイルの作成
func Register(originfile *os.File) (*Book, error) {
	newpath := generateShelfPath(*originfile, time.Now())
	err := os.Rename(originfile.Name(), newpath)
	if err != nil {
		return nil, err
	}

	// リネーム後のファイルを開く
	shelfFile, err := os.Open(newpath)
	if err != nil {
		return nil, err
	}
	if shelfFile == nil {
		return nil, fmt.Errorf("新しいshelfFileが発見できなかった")
	}
	book, err := NewBook(*shelfFile)
	metapath, err := book.MetaPath()
	if err != nil {
		return nil, err
	}
	metaFile, err := os.Create(metapath)
	if err != nil {
		return nil, err
	}
	err = book.writeBlankMetaFile(metaFile)
	if err != nil {
		return nil, err
	}

	return book, nil
}
