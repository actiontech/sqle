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
	cipherKey    = "ActionTech--SQLE"
	cipherKeyDMS = "ActionTech---DMS"

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
	KnowledgeBaseDBTypes       []string        // Knowledge base database type
}

type ClusterHardwareSign struct {
	Id   string `json:"id"`        // cluster 标识
	Sign string `json:"signature"` // license 服务签名
}

type LicenseContent struct {
	Permission           LicensePermission
	HardwareSign         string
	ClusterHardwareSigns []ClusterHardwareSign // v2.2303.0 版本引入
	LicenseId            string
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
	block, err := aes.NewCipher([]byte(cipherKey))
	if nil != err {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, []byte(cipherText))

	permissionStr, err := json.Marshal(l.Permission)
	if nil != err {
		return "", err
	}
	encodedPermissionStr, err := encode(string(permissionStr), encrypter)
	if nil != err {
		return "", err
	}

	encodedHardwareSignStr, err := encode(l.HardwareSign, encrypter)
	if nil != err {
		return "", err
	}
	clusterHardwareSignStr, err := json.Marshal(l.ClusterHardwareSigns)
	if nil != err {
		return "", err
	}

	encodedClusterHardwareSignStr, err := encode(string(clusterHardwareSignStr), encrypter)
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
	block, err := aes.NewCipher([]byte(cipherKey))
	if nil != err {
		return err
	}
	decrypter := cipher.NewCFBDecrypter(block, []byte(cipherText))

	options := strings.Split(license, DELIMITER)
	if len(options) < 3 {
		return ErrInvalidLicense
	}
	permissionStr, err := decode(options[1], decrypter)
	if nil != err {
		return err
	}
	permission := &LicensePermission{}
	err = json.Unmarshal([]byte(permissionStr), &permission)
	if nil != err {
		return err
	}
	hardwareSign, err := decode(options[2], decrypter)
	if nil != err {
		return err
	}
	if len(options) >= 4 {
		clusterHardwareSigns := []ClusterHardwareSign{}
		clusterHardwareSignStr, err := decode(options[3], decrypter)
		if nil != err {
			return err
		}
		err = json.Unmarshal([]byte(clusterHardwareSignStr), &clusterHardwareSigns)
		if err != nil {
			return err
		}
		l.ClusterHardwareSigns = clusterHardwareSigns
	}

	l.Permission = *permission
	l.HardwareSign = hardwareSign
	return nil
}

func (l *LicenseContent) DecodeDMSLicense(license string) error {
	block, err := aes.NewCipher([]byte(cipherKeyDMS))
	if nil != err {
		return err
	}
	decrypter := cipher.NewCFBDecrypter(block, []byte(cipherText))

	options := strings.Split(license, DELIMITER)
	if len(options) < 4 {
		return ErrInvalidLicense
	}
	l.LicenseId = options[0]

	permissionStr, err := decode(options[2], decrypter)
	if nil != err {
		return err
	}
	permission := &LicensePermission{}
	err = json.Unmarshal([]byte(permissionStr), &permission)
	if nil != err {
		return err
	}
	hardwareSign, err := decode(options[3], decrypter)
	if nil != err {
		return err
	}
	if len(options) >= 5 {
		clusterHardwareSigns := []ClusterHardwareSign{}
		clusterHardwareSignStr, err := decode(options[4], decrypter)
		if nil != err {
			return err
		}
		err = json.Unmarshal([]byte(clusterHardwareSignStr), &clusterHardwareSigns)
		if err != nil {
			return err
		}
		l.ClusterHardwareSigns = clusterHardwareSigns
	}

	l.Permission = *permission
	l.HardwareSign = hardwareSign
	return nil
}

func encode(str string, encrypter cipher.Stream) (string, error) {
	encrypted := make([]byte, len(str))
	encrypter.XORKeyStream(encrypted, []byte(str))
	return genEncoding.EncodeToString(encrypted), nil
}

func decode(str string, decrypter cipher.Stream) (string, error) {
	a, err := genEncoding.DecodeString(str)
	if nil != err {
		return "", err
	}
	decrypted := make([]byte, len(a))
	decrypter.XORKeyStream(decrypted, a)
	return string(decrypted), nil
}

func (l *License) CheckHardwareSignIsMatch(hardwareSign string) error {
	for _, s := range l.ClusterHardwareSigns {
		if hardwareSign == s.Sign {
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

/*
	CheckCanCreateInstance 验证可添加的数据库实例上限。
	支持对每个类型的数据库实例单独限制上限，也支持配置custom的类型上限。
	custom 指的是通用数据库类型，代表该许可不限制数据库种类。
	例如配置了custom 10, MySQL 10，则可以支持添加 20 个MySQL，或 10 个 MySQL 和其他任意10个数据库类型的实例。
*/
func (l *License) CheckCanCreateInstance(dbType string, usage LimitOfEachType) error {
	// 优先验证指定的数据库类型是否超出限制
	max := l.Permission.NumberOfInstanceOfEachType[dbType]
	cur := usage[dbType]
	if cur.Count+1 <= max.Count {
		return nil
	}

	// 当指定的数据库类型超出限制，则使用 custom 数据库类型的容量。判断添加的数据库实例是否超过 custom 类型的限制
	var customUsage int
	for _, count := range usage {
		limitation, ok := l.Permission.NumberOfInstanceOfEachType[count.DBType]
		if !ok {
			// 如果许可证里没有这个类型的数据库，则全部算custom数量
			customUsage += count.Count
			continue
		}
		// 该数据库类型使用量未超出
		if count.Count <= limitation.Count {
			continue
		}
		customUsage += count.Count - limitation.Count
	}

	maxCustom := l.Permission.NumberOfInstanceOfEachType[CustomTypeKey]
	if customUsage+1 <= maxCustom.Count {
		return nil
	}
	return fmt.Errorf("instance count reaches the limitation")
}

func CheckKnowledgeBaseLicense(content string) error {
	licenseContent := &LicenseContent{}
	err := licenseContent.DecodeDMSLicense(content)
	if err != nil {
		return err
	}
	license := &License{
		LicenseContent: *licenseContent,
	}
	return license.CheckSupportKnowledgeBase()
}

// 检查License是否支持知识库
func (l *License) CheckSupportKnowledgeBase() error {
	if len(l.Permission.KnowledgeBaseDBTypes) == 0 {
		return fmt.Errorf("knowledge base is not supported")
	}
	return nil
}

// 获取License中支持的知识库数据库类型
func (l *License) GetKnowledgeBaseDBTypes() []string {
	return l.Permission.KnowledgeBaseDBTypes
}

func GetDMSLicense(content string) (*License, error) {
	licenseContent := &LicenseContent{}
	err := licenseContent.DecodeDMSLicense(content)
	if err!= nil {
		return nil, err
	}
	license := &License{
		LicenseContent: *licenseContent,
	}
	return license, nil
}
