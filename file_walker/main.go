package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nguyenthenguyen/docx"
)

var rule map[string]string

func walkFunc(path string, info os.FileInfo, err error) error {
	old_path := path
	if err != nil {
		return err // Если по какой-то причине мы получили ошибку, проигнорируем эту итерацию
	}

	if info.IsDir() {
		return nil // Проигнорируем директории
	}

	if !strings.HasSuffix(path, ".docx") {
		for old, new := range rule {
			path = strings.ReplaceAll(path, old, new)
		}

		if strings.Compare(old_path, path) != 0 {
			os.Rename(old_path, path)
		}

		return nil
	}

	r, err := docx.ReadDocxFile(path)
	if err != nil {
		return nil
	}

	docx1 := r.Editable()

	for old, new := range rule {
		docx1.Replace(old, new, -1)
	}

	if strings.Contains(strings.ToLower(path), "основы программирования") || strings.Contains(strings.ToLower(path), "информатика") {
		for imageIndex := 1; imageIndex <= docx1.ImagesLen(); imageIndex++ {
			docx1.ReplaceImage("word/media/image"+strconv.Itoa(imageIndex)+".png", "../assets/гуриков.png")
			docx1.ReplaceImage("word/media/image"+strconv.Itoa(imageIndex)+".jpg", "../assets/гуриков.jpg")
			docx1.ReplaceImage("word/media/image"+strconv.Itoa(imageIndex)+".jpeg", "../assets/гуриков.jpg")
		}
	}

	for old, new := range rule {
		path = strings.ReplaceAll(path, old, new)
	}

	path, _ = strings.CutSuffix(path, ".docx")

	docx1.WriteToFile(path + "_src.docx")

	r.Close()

	os.Remove(old_path)

	return nil
}

var pathDoc string

func init() {
	flag.StringVar(&pathDoc, "Labs", "./МТУСИ", "path to lab directory")
}

func main() {
	flag.Parse()
	file, err := os.Open("./dict.json")
	if err != nil {
		panic(err)
	}

	info, err := file.Stat()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, info.Size())
	reader := bufio.NewReader(file)
	_, err = reader.Read(buf)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(buf, &rule)
	if err := filepath.Walk(pathDoc, walkFunc); err != nil {
		fmt.Printf("Да это какая-то ошибка: %v\n", err)
	}
}
