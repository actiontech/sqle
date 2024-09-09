package i18nPkg

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"golang.org/x/text/language"
)

func ConvertStr2I18nAsDefaultLang(s string) I18nStr {
	return map[language.Tag]string{DefaultLang: s}
}

func ConvertStrMap2I18nStr(s map[string]string) (I18nStr, error) {
	if len(s) == 0 {
		// todo i18n old plugin
		return nil, nil
	}
	if _, exist := s[DefaultLang.String()]; !exist {
		return nil, fmt.Errorf("must contains DefaultLang:%s", DefaultLang.String())
	}
	i := make(I18nStr, len(s))
	for lang, v := range s {
		langTag, err := language.Parse(lang)
		if err != nil {
			return nil, fmt.Errorf("parse key:%s err:%v", lang, err)
		}
		i[langTag] = v
	}
	return i, nil
}

type I18nStr map[language.Tag]string

// GetStrInLang if the lang not exists, return DefaultLang
func (s *I18nStr) GetStrInLang(lang language.Tag) string {
	if s == nil || *s == nil {
		return ""
	}
	if str, exist := (*s)[lang]; exist {
		return str
	}
	return (*s)[DefaultLang]
}

func (s *I18nStr) SetStrInLang(lang language.Tag, str string) {
	if *s == nil {
		*s = map[language.Tag]string{lang: str}
	} else {
		(*s)[lang] = str
	}
	return
}

func (s *I18nStr) StrMap() map[string]string {
	if s == nil || *s == nil {
		return map[string]string{}
	}
	m := make(map[string]string, len(*s))
	for langTag, v := range *s {
		m[langTag.String()] = v
	}
	return m
}

func (s *I18nStr) Copy() I18nStr {
	if s == nil || *s == nil {
		return nil
	}
	i := make(I18nStr, len(*s))
	for k, v := range *s {
		i[k] = v
	}
	return i
}

// Value impl sql. driver.Valuer interface
func (s I18nStr) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

// Scan impl sql.Scanner interface
func (s *I18nStr) Scan(input interface{}) error {
	if input == nil {
		return nil
	}
	if data, ok := input.([]byte); !ok {
		return fmt.Errorf("I18nStr Scan input is not bytes")
	} else {
		return json.Unmarshal(data, s)
	}
}
