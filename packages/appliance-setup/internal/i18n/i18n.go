package i18n

import (
	"embed"
	"encoding/json"
)

//go:embed *.json
var i8nFiles embed.FS
var i18nFilesMap = make(map[string]I18n)

var defaultLangCode = "en_US.UTF-8"

func init() {
	files, err := i8nFiles.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		data, err := i8nFiles.ReadFile(file.Name())
		if err != nil {
			panic(err)
		}
		var i18n I18n
		err = json.Unmarshal(data, &i18n)
		if err != nil {
			panic(err)
		}
		i18nFilesMap[i18n.LangCode] = i18n
	}
}

type I18n struct {
	LangCode string            `json:"langCode"`
	LangName string            `json:"langName"`
	Strings  map[string]string `json:"strings"`
}

func GetLangName(langCode string) string {
	return i18nFilesMap[langCode].LangName
}

func Get(langCode string, key string) string {
	if _, ok := i18nFilesMap[langCode]; !ok {
		return key
	}
	if _, ok := i18nFilesMap[langCode].Strings[key]; !ok {
		if langCode == defaultLangCode {
			return key
		}
		return Get(defaultLangCode, key)
	}
	return i18nFilesMap[langCode].Strings[key]
}

func GetLanguagesCodes() []string {
	langs := []string{}
	for _, lang := range i18nFilesMap {
		langs = append(langs, lang.LangCode)
	}
	return langs
}
