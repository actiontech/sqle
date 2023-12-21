package dms

var dmsServerAddress string

func GetDMSServerAddress() string {
	return dmsServerAddress
}

func InitDMSServerAddress(addr string) {
	dmsServerAddress = addr
}
