package workwx

//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/apis.md ./apis.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/chat_info.md ./chat_info.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/dept_info.md ./dept_info.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/external_contact.md ./external_contact.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/kf.md ./kf.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/user_info.md ./user_info.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/oa.md ./oa.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/rx_msg.md ./rx_msg.md.go
//go:generate go run --tags sdkcodegen ./internal/errcodegen ./errcodes/mod.go
