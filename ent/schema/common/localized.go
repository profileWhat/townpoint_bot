package common

import (
	"fmt"
	"strings"
)

const (
	En          = "en"
	Ru          = "ru"
	Sr          = "sr"
	DefaultLang = En
)

type Localized struct {
	En string `json:"en"`
	Ru string `json:"ru,omitempty"`
	Sr string `json:"sr,omitempty"`
}

type Translated interface {
	GetTranslated() Localized
}

func GetLanguages() []string {
	languages := make([]string, 0)
	languages = append(languages,
		En,
		Ru,
		Sr,
	)

	return languages
}

func IsLangExist(lang string) bool {
	langs := GetLanguages()
	for _, cLang := range langs {
		if strings.EqualFold(lang, cLang) {
			return true
		}
	}

	return false
}

func (l Localized) GetTranslate(lang string) (string, error) {
	switch {
	case strings.EqualFold(lang, En):
		return l.En, nil
	case strings.EqualFold(lang, Ru):
		return l.Ru, nil
	case strings.EqualFold(lang, Sr):
		return l.Sr, nil
	default:
		return "", fmt.Errorf("not supported language: %v", lang)
	}
}

func (l Localized) GetTranslateWithDefault(lang string) string {
	switch {
	case strings.EqualFold(lang, Ru):
		return l.Ru
	case strings.EqualFold(lang, Sr):
		return l.Sr
	default:
		return l.En
	}
}

func (l Localized) FillValue(lang string, value string) Localized {
	newLocalized := l
	switch {
	case strings.EqualFold(lang, En):
		newLocalized.En = value
	case strings.EqualFold(lang, Ru):
		newLocalized.Ru = value
	case strings.EqualFold(lang, Sr):
		newLocalized.Sr = value
	}

	return newLocalized
}
