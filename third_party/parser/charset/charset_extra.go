package charset

import "strings"

// 1. support utf8mb3, 最新解析器代码是支持了 utf8mb3, 但是解析器库又挪回了TiDB仓库，整个改动太大了。因此目前还是基于老的独立解析器仓库进行定制。除非有其他新SQL不得不支持。
// 2. 支持能成功解析出所有字符集和排序。当前的解析库解析SQL时，如果存在 TiDB 不支持的字符集和排序则会显式的报错，不符合预期。
func InitAllCharset() {
	// 将所有的字符集都放进 `charsets`, 这是个代表当前支持的字符集。
	for _, c := range collations {
		if charset, ok := charsets[c.CharsetName]; ok {
			charset.Collations[c.Name] = c
			if c.IsDefault {
				charset.DefaultCollation = c.Name
			}
		} else {
			charsets[c.CharsetName] = &Charset{
				Name:             c.CharsetName,
				DefaultCollation: c.Name,
				Collations: map[string]*Collation{
					c.Name: c,
				},
			}
		}
	}

	// utf8mb3 直接引用 utf8 的 charset.
	charsets["utf8mb3"] = charsets["utf8"]

	// 将所有 utf8 的字符集排序都使用 utf8mb3 别名重新引用一次，保证通过 utf8mb3 相关字符集排序解析器支持。
	utf8mb3Collations := []*Collation{}
	for _, c := range collations {
		if c.CharsetName == "utf8" {
			aliasName := strings.Replace(c.Name, "utf8_", "utf8mb3_", 1)
			collationsNameMap[aliasName] = c
			utf8mb3Collations = append(utf8mb3Collations, c)
		}
	}
	// collations 原始的记录表，此处的改动是为了保证单测 `TestGetDefaultCollation` 成功。
	collations = append(collations, utf8mb3Collations...)
}
