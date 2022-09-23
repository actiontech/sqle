package config

var Version string

// 在登录校验通过后, Login接口会将登陆成功的用户名写入echo.Context
const LoginUserNameKey = "login_user_name"