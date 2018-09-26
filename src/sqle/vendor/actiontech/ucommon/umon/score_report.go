package umon

type ScoreReport struct {
	MysqlId string `json:"mysql_id"`
	MysqlAlias string `json:"mysql_alias"`
	MysqlPort string `json:"mysql_port"`
	MysqlVersion string `json:"mysql_version"`

	NonStandardCharacterDatabase []map[string]string `json:"non_standard_character_database"`
	NonStandardCharacterDatabaseError string `json:"non_standard_character_database_error"`

	NonStandardCharacterTable []map[string]string `json:"non_standard_character_table"`
	NonStandardCharacterTableError string `json:"non_standard_character_table_error"`

	NoUseInnodbTable []map[string]string `json:"no_use_innodb_table"`
	NoUseInnodbTableError string `json:"no_use_innodb_table_error"`

	NoPrimarykeyTable []map[string]string `json:"no_primary_key_table"`
	NoPrimarykeyTableError string `json:"no_primary_key_table_error"`

	RedundancyIndexsTable []map[string]string `json:"redundancy_indexs_table"`
	RedundancyIndexsTableError string `json:"redundancy_indexs_table_error"`
}
