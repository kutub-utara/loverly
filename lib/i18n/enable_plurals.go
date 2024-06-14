package i18n

import (
	"strings"

	"github.com/nicksnyder/go-i18n/i18n/language"
)

var langPreferences = []string{"en-ID", "id-ID", "th-TH", "vi-VN"}

func init() {
	enablePluralsForAllSupportedLanguages(langPreferences)
}

func enablePluralsForAllSupportedLanguages(languages []string) {
	var languageTags []string

	pluralSet := map[language.Plural]struct{}{
		language.One:   {},
		language.Two:   {},
		language.Other: {},
	}

	for _, l := range languages {
		lang := strings.ToLower(strings.Split(l, "-")[0])
		languageTags = append(languageTags, lang)
	}

	language.RegisterPluralSpec(languageTags, &language.PluralSpec{
		Plurals: pluralSet,
		PluralFunc: func(ops *language.Operands) language.Plural {
			if ops.NequalsAny(1) {
				return language.One
			}
			if ops.NequalsAny(2) {
				return language.Two
			}
			return language.Other
		},
	})
}
