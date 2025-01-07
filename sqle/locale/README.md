在文件`message_zh.go`中新增、修改 `i18n.Message` 后需更新对应语言文件（active.*.toml），更新的脚本写在项目Makefile中了，使用步骤如下：
1. 安装需要的工具，已安装则跳过: \
   `make install_i18n_tool`
2. 根据i18n.Message生成中文语言包文件(active.zh.toml): \
   `make extract_i18n`
3. 如果只需要支持中文，后续就不需要执行了
4. 生成待翻译的临时文件(translate.en.toml): \
   `make start_trans_i18n`
5. 人工介入将 translate.en.toml 文件中的文本翻译替换
6. 根据翻译好的文本生成英文语言包文件(active.en.toml): \
   `make end_trans_i18n`