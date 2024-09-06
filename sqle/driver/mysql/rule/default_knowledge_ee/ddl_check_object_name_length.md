样例说明：

```
CREATE TABLE this_is_a_64_character_table_name_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx (
    this_is_a_64_character_column_name_xxxxxxxxxxxxxxxxxxxxxxxxxxxxx INT,
    INDEX this_is_a_64_character_index_name_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx (this_is_a_64_character_column_name_xxxxxxxxxxxxxxxxxxxxxxxxxxxxx)
);-- 不能超过阈值
```