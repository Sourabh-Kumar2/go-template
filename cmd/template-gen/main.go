package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

func main() {
	flg, err := parseFlags()
	if err != nil {
		flag.PrintDefaults()
		return
	}
	file, testFile, err := getFileWriters(flg)
	handleErr(err)
	defer file.Close()
	defer testFile.Close()
	templ := template.Must(template.New("http.tmpl").Funcs(getFuncMap()).ParseGlob("templates/*"))
	handleErr(templ.Execute(file, flg.handler))
	handleErr(templ.ExecuteTemplate(testFile, "http_test.tmpl", nil))
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getFuncMap() template.FuncMap {
	return map[string]any{
		"toTitle": func(s string) string {
			r, size := utf8.DecodeRuneInString(s)
			return string(unicode.ToTitle(r)) + s[size:]
		},
	}
}

func getFileWriters(flg *flags) (mainFile *os.File, testFile *os.File, err error) {
	var filepath string
	filepath = path.Join(flg.path, toSnakeCase(flg.handler)) + ".go"
	mainFile, err = createFile(filepath)
	if err != nil {
		return nil, nil, err
	}
	filepath = path.Join(flg.path, toSnakeCase(flg.handler)) + "_test.go"
	testFile, err = createFile(filepath)
	if err != nil {
		return nil, nil, err
	}

	return mainFile, testFile, nil
}

func createFile(filepath string) (*os.File, error) {
	_, err := os.Stat(filepath)
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(filepath)
		if err != nil {
			return nil, fmt.Errorf("cannot create file %q: %w", filepath, err)
		}
		return f, nil
	}
	return nil, fmt.Errorf("cannot override file: %w", err)
}

func toSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
