package v1

type CheckDbConnectable struct {
	// DB Service type
	// Required: true
	// example: MySQL
	DBType string `json:"db_type"  example:"mysql" validate:"required"`
	// DB Service admin user
	// Required: true
	// example: root
	User string `json:"user"  example:"root" valid:"required"`
	// DB Service host
	// Required: true
	// example: 127.0.0.1
	Host string `json:"host"  example:"10.10.10.10" valid:"required,ip_addr|uri|hostname|hostname_rfc1123"`
	// DB Service port
	// Required: true
	// example: 3306
	Port string `json:"port"  example:"3306" valid:"required,port"`
	// DB Service admin password
	// Required: true
	// example: 123456
	Password string `json:"password"  example:"123456"`
	// DB Service Custom connection parameters
	// Required: false
	AdditionalParams []*AdditionalParam `json:"additional_params" from:"additional_params"`
}

type AdditionalParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
