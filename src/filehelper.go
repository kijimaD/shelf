package src

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const IDFormat = "20060102T150405"

// "{id}_{base}.{ext}"
type Fullname string

// 20240303T201703
func formatID(t time.Time) string {
	return t.Format(IDFormat)
}

// Fullnameの文字列からFullnameを生成する
func LoadFullname(full string) (Fullname, error) {
	pattern := `(?P<id>\w{15})_(?P<name>.*)\.`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	matches := re.FindAllStringSubmatch(full, -1)
	for _, match := range matches {
		fmt.Printf("ID: %s, Name: %s\n", match[re.SubexpIndex("id")], match[re.SubexpIndex("name")])
	}

	return "", nil
}

// IDなしファイル名からFullnameを生成する
func GenFullname(origin string, t time.Time) (Fullname, error) {
	// 拡張子がなければエラー
	ext := filepath.Ext(origin)
	if ext == "" {
		return "", fmt.Errorf("拡張子がない")
	}

	id := formatID(t)
	base := strings.TrimSuffix(origin, ext)
	result := fmt.Sprintf("%s_%s%s", id, base, ext)

	return Fullname(result), nil
}

func (full *Fullname) ID() {

}
