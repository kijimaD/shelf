package shelf

import (
	"os"
	"time"
)

// 入力パスに対して、ID付きファイル名にリネーム + メタファイルの作成
func Register(path string) error {
	fullname, err := NewFullnameByRaw(path, time.Now())
	if err != nil {
		return err
	}
	file, err := os.Create(fullname.MetaPath())
	if err != nil {
		return err
	}
	defer file.Close()
	fullname.blankMetafile(file)

	_, err = fullname.rename()
	if err != nil {
		return err
	}

	return nil
}
