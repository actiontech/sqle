//go:build release
// +build release

package license

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/utils"
)

var relAse = utils.NewEncryptor(config.RelAesKey)

var (
	DELIMITER        = ";;"
	genEncoding      = base64.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_~")
	collectEncoding  = base64.NewEncoding("012345ghijklmnopq6789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefrstuvwxyz_~")
	ExpireDateFormat = "2006-01-02"
)

const (
	// must be 16 byte
	cipherKey = "ActionTech--SQLE"

	// must be 16 byte
	cipherText = "a1b2c3d4e5f6g7h8"
)

var (
	ErrInvalidLicense     = errors.New("invalid license")
	ErrLicenseEmpty       = errors.New("license is empty")
	ErrCollectLicenseInfo = errors.New("collect license info error")
)

type LimitOfEachType map[string] /*db type*/ LimitOfType

type LimitOfType struct {
	DBType string `json:"db_type"`
	Count  int    `json:"count"`
}

type LicensePermission struct {
	WorkDurationDay            int             // How long SQLE is authorized to run
	Version                    string          // SQLE version
	UserCount                  int             // The number of user
	NumberOfInstanceOfEachType LimitOfEachType // Instance limit
}

type LicenseContent struct {
	Permission          LicensePermission
	HardwareSign        string
	ClusterHardwareSign map[string]string // v2.2303.0 版本引入
}
type LicenseStatus struct {
	WorkDurationHour int // 实际的运行时间，加密存在许可证内容里
}

type License struct {
	LicenseContent
	LicenseStatus
}

/*
License Format:
	license: [LicenseStatus.WorkDurationHour]~~[LicenseContent]

	[WorkDurationHour] is number

	[LicenseContent]: [LicenseDesc];;[Permission];;[HardwareSign];;[ClusterHardwareSign]

Example:
	10~~This license is for: &{WorkDurationDay:60 Version:演示环境 UserCount:10 NumberOfInstanceOfEachType:map[custom:{DBType:custom Count:3} mysql:{DBType:mysql Count:3}]};;1_XBm2N8t7coUEuhg7J5V8o9AYlhUfq2AmndctDHCxz9u~GyOKyJW0e~sVDuQVbkaKzAZQvpsGBqB~liD7svsTvbzD3ZHfdvEtSPkoYSnk2nxrYJLrW0wmzTVIicDWg1Dp2MICEK9T09Od3Xn1u4XWO7e182mzrHqncLOGKXJKlSrCsL_kWY6o6w8pWKL1Xdzduyq4uLdXuL9E6oOzyUMF3rYlnOhvoOwdoE;;9S~ViK_ZoRx8045cLM5pTZXCCpDEY_yxjfaLYGBMMOKyWpgc
*/

func (l *License) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("failed to unmarshal license value:", value))
	}
	licenseInfo, err := relAse.AesDecrypt(string(data))
	if err != nil {
		return err
	}
	separate := strings.Index(licenseInfo, "~~")
	if separate == -1 {
		return fmt.Errorf("failed to unmarshal license value: %s", data)
	}
	l.WorkDurationHour, _ = strconv.Atoi(licenseInfo[:separate])

	license := licenseInfo[separate+2:]
	return l.LicenseContent.Decode(license)
}

func (l License) Value() (driver.Value, error) {
	licenseContentStr, err := l.LicenseContent.Encode()
	if nil != err {
		return nil, err
	}
	return relAse.AesEncrypt(fmt.Sprintf("%v~~%v", l.WorkDurationHour, licenseContentStr))
}

func (l *LicenseContent) Encode() (text string, err error) {
	permissionStr, err := json.Marshal(l.Permission)
	if nil != err {
		return "", err
	}
	encodedPermissionStr, err := encode(string(permissionStr))
	if nil != err {
		return "", err
	}

	encodedHardwareSignStr, err := encode(l.HardwareSign)
	if nil != err {
		return "", err
	}
	clusterHardwareSignStr, err := json.Marshal(l.ClusterHardwareSign)
	if nil != err {
		return "", err
	}

	encodedClusterHardwareSignStr, err := encode(string(clusterHardwareSignStr))
	if nil != err {
		return "", err
	}

	ret := make([]string, 0, 4)
	licenseInfo := fmt.Sprintf("This license is for: %+v", l.Permission)
	ret = append(ret, licenseInfo)
	ret = append(ret, encodedPermissionStr)
	ret = append(ret, encodedHardwareSignStr)
	ret = append(ret, encodedClusterHardwareSignStr)
	return strings.Join(ret, DELIMITER), nil
}

func (l *LicenseContent) Decode(license string) error {
	options := strings.Split(license, DELIMITER)
	if len(options) < 3 {
		return ErrInvalidLicense
	}
	permissionStr, err := decode(options[1])
	if nil != err {
		return err
	}
	permission := &LicensePermission{}
	err = json.Unmarshal([]byte(permissionStr), &permission)
	if nil != err {
		return err
	}
	hardwareSign, err := decode(options[2])
	if nil != err {
		return err
	}
	if len(options) >= 4 {
		clusterHardwareSign := map[string]string{}
		clusterHardwareSignStr, err := decode(options[3])
		if nil != err {
			return err
		}
		err = json.Unmarshal([]byte(clusterHardwareSignStr), &clusterHardwareSign)
		if err != nil {
			return err
		}
		l.ClusterHardwareSign = clusterHardwareSign
	}

	l.Permission = *permission
	l.HardwareSign = hardwareSign
	return nil
}

func encode(str string) (string, error) {
	block, err := aes.NewCipher([]byte(cipherKey))
	if nil != err {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, []byte(cipherText))
	encrypted := make([]byte, len(str))
	encrypter.XORKeyStream(encrypted, []byte(str))
	return genEncoding.EncodeToString(encrypted), nil
}

func decode(str string) (string, error) {
	block, err := aes.NewCipher([]byte(cipherKey))
	if nil != err {
		return "", err
	}
	decrypter := cipher.NewCFBDecrypter(block, []byte(cipherText))
	a, err := genEncoding.DecodeString(str)
	if nil != err {
		return "", err
	}
	decrypted := make([]byte, len(a))
	decrypter.XORKeyStream(decrypted, a)
	return string(decrypted), nil
}

func (l *License) CheckHardwareSignIsMatch(hardwareSign string) error {
	for _, s := range l.ClusterHardwareSign {
		if hardwareSign == s {
			return nil
		}
	}
	if hardwareSign == l.HardwareSign {
		return nil
	}
	return fmt.Errorf("the server is not in the license")
}

func (l *License) CheckLicenseNotExpired() error {
	if l.WorkDurationHour >= l.LicenseContent.Permission.WorkDurationDay*24 {
		return fmt.Errorf("license is expired")
	}
	return nil
}

func (c *License) CheckCanCreateUser(userCount int64) error {
	if userCount+1 > int64(c.Permission.UserCount) {
		return fmt.Errorf("user count reaches the limitation")
	}
	return nil
}

const CustomTypeKey = "custom"

func (l *License) CheckCanCreateInstance(dbType string, usage LimitOfEachType) error {
	// 指定的数据库类型没有超出限制
	max := l.Permission.NumberOfInstanceOfEachType[dbType]
	cur := usage[dbType]
	if cur.Count+1 <= max.Count {
		return nil
	}

	// 指定的数据库类型超出限制，判断是否需要custom 类型的是否超出
	var customUsage int
	for _, count := range usage {
		total, ok := l.Permission.NumberOfInstanceOfEachType[count.DBType]
		if !ok {
			// 如果许可证里没有这个类型的数据库，则全部算custom数量
			customUsage += count.Count
			continue
		}
		// 该数据库类型使用量未超出
		if count.Count <= total.Count {
			continue
		}
		customUsage += count.Count - total.Count
	}

	maxCustom := l.Permission.NumberOfInstanceOfEachType[CustomTypeKey]
	if customUsage+1 <= maxCustom.Count {
		return nil
	}
	return fmt.Errorf("instance count reaches the limitation")
}
