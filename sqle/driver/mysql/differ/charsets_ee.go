//go:build enterprise
// +build enterprise

package differ

// collationIsDefault returns true if the supplied collation is the default
// collation for the supplied charset in flavor. Results not guaranteed to be
// accurate for invalid input (i.e. mismatched collation/charset or combinations
// that do not exist in the supplied flavor), or in MariaDB 11.2+ since its
// @@character_set_collations feature permits arbitrary overrides.
func collationIsDefault(collation, charset string, flavor Flavor) bool {
	// Handle special cases manually
	if collation == "utf8mb4_0900_ai_ci" || collation == "utf8mb3_general_ci" || collation == "utf8_general_ci" {
		return true
	} else if collation == "utf8mb4_general_ci" {
		return !flavor.MinMySQL(8)
	} else {
		return baseDefaultCollation[charset] == collation
	}
}

// SELECT CONCAT('"', character_set_name, '": "', default_collate_name, '",') from information_schema.character_sets ORDER BY character_set_name;
var baseDefaultCollation = map[string]string{
	"armscii8": "armscii8_general_ci",
	"ascii":    "ascii_general_ci",
	"big5":     "big5_chinese_ci",
	"binary":   "binary",
	"cp1250":   "cp1250_general_ci",
	"cp1251":   "cp1251_general_ci",
	"cp1256":   "cp1256_general_ci",
	"cp1257":   "cp1257_general_ci",
	"cp850":    "cp850_general_ci",
	"cp852":    "cp852_general_ci",
	"cp866":    "cp866_general_ci",
	"cp932":    "cp932_japanese_ci",
	"dec8":     "dec8_swedish_ci",
	"eucjpms":  "eucjpms_japanese_ci",
	"euckr":    "euckr_korean_ci",
	"gb18030":  "gb18030_chinese_ci", // added in MySQL 5.7
	"gb2312":   "gb2312_chinese_ci",
	"gbk":      "gbk_chinese_ci",
	"geostd8":  "geostd8_general_ci",
	"greek":    "greek_general_ci",
	"hebrew":   "hebrew_general_ci",
	"hp8":      "hp8_english_ci",
	"keybcs2":  "keybcs2_general_ci",
	"koi8r":    "koi8r_general_ci",
	"koi8u":    "koi8u_general_ci",
	"latin1":   "latin1_swedish_ci",
	"latin2":   "latin2_general_ci",
	"latin5":   "latin5_turkish_ci",
	"latin7":   "latin7_general_ci",
	"macce":    "macce_general_ci",
	"macroman": "macroman_general_ci",
	"sjis":     "sjis_japanese_ci",
	"swe7":     "swe7_swedish_ci",
	"tis620":   "tis620_thai_ci",
	"ucs2":     "ucs2_general_ci",
	"ujis":     "ujis_japanese_ci",
	"utf16":    "utf16_general_ci",
	"utf16le":  "utf16le_general_ci", // added in MySQL 5.6, also present in MariaDB
	"utf32":    "utf32_general_ci",
	"utf8":     "utf8_general_ci",    // removed in MySQL 8.0.29 and MariaDB 10.6
	"utf8mb3":  "utf8mb3_general_ci", // added in MySQL 8.0.29 (with default of "utf8_general_ci" until 8.0.30) and MariaDB 10.6
	"utf8mb4":  "utf8mb4_general_ci", // default changes to "utf8mb4_0900_ai_ci" in MySQL 8
}
