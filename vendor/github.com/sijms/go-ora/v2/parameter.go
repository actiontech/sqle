package go_ora

import (
	"database/sql/driver"
	"github.com/sijms/go-ora/v2/network"
	"math"
	"strings"
	"time"
)

type OracleType int
type ParameterDirection int
type NVarChar string
type TimeStamp time.Time

//func (n *NVarChar) ConvertValue(v interface{}) (driver.Value, error) {
//	return driver.Value(string(*n)), nil
//}

func (n *NVarChar) Value() (driver.Value, error) {
	return driver.Value(string(*n)), nil
}

const (
	Input  ParameterDirection = 1
	Output ParameterDirection = 2
	InOut  ParameterDirection = 3
	RetVal ParameterDirection = 9
)

//internal enum BindDirection
//{
//Output = 16,
//Input = 32,
//InputOutput = 48,
//}

//go:generate stringer -type=OracleType

const (
	NCHAR            OracleType = 1
	NUMBER           OracleType = 2
	SB1              OracleType = 3
	SB2              OracleType = 3
	SB4              OracleType = 3
	FLOAT            OracleType = 4
	NullStr          OracleType = 5
	VarNum           OracleType = 6
	LONG             OracleType = 8
	VARCHAR          OracleType = 9
	ROWID            OracleType = 11
	DATE             OracleType = 12
	VarRaw           OracleType = 15
	BFloat           OracleType = 21
	BDouble          OracleType = 22
	RAW              OracleType = 23
	LongRaw          OracleType = 24
	UINT             OracleType = 68
	LongVarChar      OracleType = 94
	LongVarRaw       OracleType = 95
	CHAR             OracleType = 96
	CHARZ            OracleType = 97
	IBFloat          OracleType = 100
	IBDouble         OracleType = 101
	REFCURSOR        OracleType = 102
	OCIXMLType       OracleType = 108
	XMLType          OracleType = 109
	OCIRef           OracleType = 110
	OCIClobLocator   OracleType = 112
	OCIBlobLocator   OracleType = 113
	OCIFileLocator   OracleType = 114
	ResultSet        OracleType = 116
	OCIString        OracleType = 155
	OCIDate          OracleType = 156
	TimeStampDTY     OracleType = 180
	TimeStampTZ_DTY  OracleType = 181
	IntervalYM_DTY   OracleType = 182
	IntervalDS_DTY   OracleType = 183
	TimeTZ           OracleType = 186
	TIMESTAMP        OracleType = 187
	TimeStampTZ      OracleType = 188
	IntervalYM       OracleType = 189
	IntervalDS       OracleType = 190
	UROWID           OracleType = 208
	TimeStampLTZ_DTY OracleType = 231
	TimeStampeLTZ    OracleType = 232
)

type ParameterType int

const (
	Number ParameterType = 1
	String ParameterType = 2
)

type ParameterInfo struct {
	Name                 string
	TypeName             string
	Direction            ParameterDirection
	IsNull               bool
	AllowNull            bool
	ColAlias             string
	DataType             OracleType
	IsXmlType            bool
	Flag                 uint8
	Precision            uint8
	Scale                uint8
	MaxLen               int
	MaxCharLen           int
	MaxNoOfArrayElements int
	ContFlag             int
	ToID                 []byte
	Version              int
	CharsetID            int
	CharsetForm          int
	BValue               []byte
	Value                driver.Value
	OutputVarPtr         interface{}
	getDataFromServer    bool
	oaccollid            int
	cusType              *customType
}

// load get parameter information form network session
func (par *ParameterInfo) load(conn *Connection) error {
	session := conn.session
	par.getDataFromServer = true
	dataType, err := session.GetByte()
	if err != nil {
		return err
	}
	par.DataType = OracleType(dataType)
	par.Flag, err = session.GetByte()
	if err != nil {
		return err
	}
	par.Precision, err = session.GetByte()
	//precision, err := session.GetInt(1, false, false)
	//var scale int
	switch par.DataType {
	case NUMBER:
		fallthrough
	case TimeStampDTY:
		fallthrough
	case TimeStampTZ_DTY:
		fallthrough
	case IntervalDS_DTY:
		fallthrough
	case TIMESTAMP:
		fallthrough
	case TimeStampTZ:
		fallthrough
	case IntervalDS:
		fallthrough
	case TimeStampLTZ_DTY:
		fallthrough
	case TimeStampeLTZ:
		if scale, err := session.GetInt(2, true, true); err != nil {
			return err
		} else {
			if scale == -127 {
				par.Precision = uint8(math.Ceil(float64(par.Precision) * 0.30103))
				par.Scale = 0xFF
			} else {
				par.Scale = uint8(scale)
			}
		}
	default:
		par.Scale, err = session.GetByte()
		//scale, err = session.GetInt(1, false, false)
	}
	//if par.Scale == uint8(-127) {
	//
	//}
	if par.DataType == NUMBER && par.Precision == 0 && (par.Scale == 0 || par.Scale == 0xFF) {
		par.Precision = 38
		par.Scale = 0xFF
	}

	//par.Scale = uint16(scale)
	//par.Precision = uint16(precision)
	par.MaxLen, err = session.GetInt(4, true, true)
	if err != nil {
		return err
	}
	switch par.DataType {
	case ROWID:
		par.MaxLen = 128
	case DATE:
		par.MaxLen = 7
	case IBFloat:
		par.MaxLen = 4
	case IBDouble:
		par.MaxLen = 8
	case TimeStampTZ_DTY:
		par.MaxLen = 13
	case IntervalYM_DTY:
		fallthrough
	case IntervalDS_DTY:
		fallthrough
	case IntervalYM:
		fallthrough
	case IntervalDS:
		par.MaxLen = 11
	}
	par.MaxNoOfArrayElements, err = session.GetInt(4, true, true)
	if err != nil {
		return err
	}
	if session.TTCVersion >= 10 {
		par.ContFlag, err = session.GetInt(8, true, true)
	} else {
		par.ContFlag, err = session.GetInt(4, true, true)
	}
	if err != nil {
		return err
	}
	par.ToID, err = session.GetDlc()
	par.Version, err = session.GetInt(2, true, true)
	if err != nil {
		return err
	}
	par.CharsetID, err = session.GetInt(2, true, true)
	if err != nil {
		return err
	}
	par.CharsetForm, err = session.GetInt(1, false, false)
	if err != nil {
		return err
	}
	par.MaxCharLen, err = session.GetInt(4, true, true)
	if err != nil {
		return err
	}
	if session.TTCVersion >= 8 {
		par.oaccollid, err = session.GetInt(4, true, true)
	}
	num1, err := session.GetInt(1, false, false)
	if err != nil {
		return err
	}
	par.AllowNull = num1 > 0
	_, err = session.GetByte() //  session.GetInt(1, false, false)
	if err != nil {
		return err
	}
	bName, err := session.GetDlc()
	if err != nil {
		return err
	}
	par.Name = session.StrConv.Decode(bName)
	_, err = session.GetDlc()
	bName, err = session.GetDlc()
	if err != nil {
		return err
	}
	par.TypeName = strings.ToUpper(session.StrConv.Decode(bName))
	if par.DataType == XMLType && par.TypeName != "XMLTYPE" {
		for typName, cusTyp := range conn.cusTyp {
			if typName == par.TypeName {
				par.cusType = &cusTyp
			}
		}
	}
	if par.TypeName == "XMLTYPE" {
		par.DataType = XMLType
		par.IsXmlType = true
	}
	if session.TTCVersion < 3 {
		return nil
	}
	_, err = session.GetInt(2, true, true)
	if session.TTCVersion < 6 {
		return nil
	}
	_, err = session.GetInt(4, true, true)
	return nil
}

// write parameter information to network session
func (par *ParameterInfo) write(session *network.Session) error {
	session.PutBytes(uint8(par.DataType), par.Flag, par.Precision, par.Scale)
	session.PutUint(par.MaxLen, 4, true, true)
	session.PutInt(par.MaxNoOfArrayElements, 4, true, true)
	if session.TTCVersion >= 10 {
		session.PutInt(par.ContFlag, 8, true, true)
	} else {
		session.PutInt(par.ContFlag, 4, true, true)
	}
	if par.ToID == nil {
		session.PutBytes(0)
		//session.PutInt(0, 1, false, false)
	} else {
		session.PutInt(len(par.ToID), 4, true, true)
		session.PutClr(par.ToID)
	}
	session.PutUint(par.Version, 2, true, true)
	session.PutUint(par.CharsetID, 2, true, true)
	session.PutBytes(uint8(par.CharsetForm))
	//session.PutUint(par.CharsetForm, 1, false, false)
	session.PutUint(par.MaxCharLen, 4, true, true)
	if session.TTCVersion >= 8 {
		session.PutInt(par.oaccollid, 4, true, true)
	}
	return nil
}
