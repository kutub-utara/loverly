package i18n

import (
	"context"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n"
)

const fileSuffix = ".all.json"

type Translator struct {
	definitions     map[string]i18n.TranslateFunc
	defaultLanguage string
}

var translator *Translator

var basepath = func() string {
	// See https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
	_, f, _, _ := runtime.Caller(0)
	// Return the project root directory path.
	return filepath.Dir(filepath.Dir(f))
}()

func Init(ctx context.Context, filePath, appFilePath, defaultLang string) error {
	translator = &Translator{
		defaultLanguage: defaultLang,
		definitions:     make(map[string]i18n.TranslateFunc),
	}

	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(basepath, filePath)
	}

	files, err := filepath.Glob(path.Join(filePath, "*"+fileSuffix))
	if err != nil {
		return fmt.Errorf("loading translation file from path: %s err: %w", filePath, err)
	}

	appTransFile, err := filepath.Glob(path.Join(appFilePath, "*"+fileSuffix))
	if err != nil {
		return fmt.Errorf("loading translation file from path: %s err: %w", appFilePath, err)
	}

	files = append(files, appTransFile...)

	for _, file := range files {
		i18n.MustLoadTranslationFile(file)
		lang := strings.TrimSuffix(filepath.Base(file), fileSuffix)
		translator.definitions[lang] = i18n.MustTfunc(lang)
	}

	if _, ok := translator.definitions[translator.defaultLanguage]; !ok {
		return fmt.Errorf("tranlation is missing for default language: %s. check translation file: %s is present",
			translator.defaultLanguage, filepath.Join(filePath, translator.defaultLanguage+fileSuffix))
	}

	return nil
}

func Translate(language, key string, args ...interface{}) string {
	return translator.translate(language, key, args...)
}

func (tlr *Translator) translate(language, key string, args ...interface{}) string {
	if tFunc, ok := tlr.definitions[tlr.mapLanguage(language)]; ok {
		return tFunc(key, args...)
	}
	return key
}

func (tlr *Translator) mapLanguage(lang string) string {
	switch lang {
	case "id", "id-ID", "id_ID", "in-ID", "in_ID", "in-id", "in_id":
		return "id-ID"
	case "en-id", "en_id", "en-ID", "en_ID", "en-TH", "en_TH", "en-VN", "en_VN":
		return "en-ID"
	case "th-TH", "th_TH":
		return "th-TH"
	case "vi-VN", "vi_VN", "vi-ID", "vi_ID":
		return "vi-VN"
	default:
		log.Printf("translator not found: fallback to default language %s", lang)
		return tlr.defaultLanguage // use DefaultLanguage
	}
}

func Message(lang, key string, args ...interface{}) string {
	return Translate(lang, fmt.Sprintf("%s_message", strings.ToLower(key)), args...)
}

func Title(lang, key string, args ...interface{}) string {
	return Translate(lang, getTitleKey(key), args...)
}

func SubTitle(lang, key string, args ...interface{}) string {
	return Translate(lang, fmt.Sprintf("%s_subtitle", strings.ToLower(key)), args...)
}

func HasTitle(lang, key string) bool {
	return Title(lang, key) != getTitleKey(key)
}

func GetSupportedLocale(lang string) string {
	return translator.mapLanguage(lang)
}

func getTitleKey(key string) string {
	return fmt.Sprintf("%s_title", strings.ToLower(key))
}
