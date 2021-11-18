package go_ora

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/sijms/go-ora/v2/trace"
	"io"
	"reflect"
	"time"

	"github.com/sijms/go-ora/v2/network"
)

// Compile time Sentinels for implemented Interfaces.
var _ = driver.Rows((*DataSet)(nil))
var _ = driver.RowsColumnTypeDatabaseTypeName((*DataSet)(nil))
var _ = driver.RowsColumnTypeLength((*DataSet)(nil))
var _ = driver.RowsColumnTypeNullable((*DataSet)(nil))

// var _ = driver.RowsColumnTypePrecisionScale((*DataSet)(nil))
// var _ = driver.RowsColumnTypeScanType((*DataSet)(nil))
// var _ = driver.RowsNextResultSet((*DataSet)(nil))

type Row []driver.Value

type DataSet struct {
	ColumnCount     int
	RowCount        int
	UACBufferLength int
	MaxRowSize      int
	Cols            []ParameterInfo
	Rows            []Row
	currentRow      Row
	lasterr         error
	index           int
	parent          StmtInterface
}

// load Loading dataset information from network session
func (dataSet *DataSet) load(session *network.Session) error {
	_, err := session.GetByte()
	if err != nil {
		return err
	}
	columnCount, err := session.GetInt(2, true, true)
	if err != nil {
		return err
	}
	num, err := session.GetInt(4, true, true)
	if err != nil {
		return err
	}
	columnCount += num * 0x100
	if columnCount > dataSet.ColumnCount {
		dataSet.ColumnCount = columnCount
	}
	if len(dataSet.currentRow) != dataSet.ColumnCount {
		dataSet.currentRow = make(Row, dataSet.ColumnCount)
	}
	dataSet.RowCount, err = session.GetInt(4, true, true)
	if err != nil {
		return err
	}
	dataSet.UACBufferLength, err = session.GetInt(2, true, true)
	if err != nil {
		return err
	}
	bitVector, err := session.GetDlc()
	if err != nil {
		return err
	}
	dataSet.setBitVector(bitVector)
	_, err = session.GetDlc()
	return nil
}

// setBitVector bit vector is an array of bit that define which column need to be read
// from network session
func (dataSet *DataSet) setBitVector(bitVector []byte) {
	index := dataSet.ColumnCount / 8
	if dataSet.ColumnCount%8 > 0 {
		index++
	}
	if len(bitVector) > 0 {
		for x := 0; x < len(bitVector); x++ {
			for i := 0; i < 8; i++ {
				if (x*8)+i < dataSet.ColumnCount {
					dataSet.Cols[(x*8)+i].getDataFromServer = bitVector[x]&(1<<i) > 0
				}
			}
		}
	} else {
		for x := 0; x < len(dataSet.Cols); x++ {
			dataSet.Cols[x].getDataFromServer = true
		}
	}

}

func (dataSet *DataSet) Close() error {
	return nil
}

// Next_ act like Next in sql package return false if no other rows in dataset
func (dataSet *DataSet) Next_() bool {
	err := dataSet.Next(dataSet.currentRow)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false
		}
		dataSet.lasterr = err
		return false
	}

	return true
}

// Scan act like scan in sql package return row values to dest variable pointers
func (dataSet *DataSet) Scan(dest ...interface{}) error {
	if dataSet.lasterr != nil {
		return dataSet.lasterr
	}
	if len(dest) != len(dataSet.currentRow) {
		return fmt.Errorf("go-ora: expected %d destination arguments in Scan, not %d",
			len(dataSet.currentRow), len(dest))
	}
	for i, col := range dataSet.currentRow {
		if dest[i] == nil {
			return fmt.Errorf("go-ora: argument %d is nil", i)
		}
		destTyp := reflect.TypeOf(dest[i])
		if destTyp.Kind() != reflect.Ptr {
			return errors.New("go-ora: argument in scan should be passed as pointers")
		}

		switch col.(type) {
		case string:
			switch destTyp.Elem().Kind() {
			case reflect.String:
				reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(col))
			default:
				return fmt.Errorf("go-ora: column %d require type string", i)
			}
		case int64:
			switch destTyp.Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				reflect.ValueOf(dest[i]).Elem().SetInt(reflect.ValueOf(col).Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				reflect.ValueOf(dest[i]).Elem().SetUint(uint64(reflect.ValueOf(col).Int()))
			case reflect.Float32, reflect.Float64:
				reflect.ValueOf(dest[i]).Elem().SetFloat(float64(reflect.ValueOf(col).Int()))
			default:
				return fmt.Errorf("go-ora: column %d require an integer", i)
			}
		case float64:
			switch destTyp.Elem().Kind() {
			case reflect.Float32, reflect.Float64:
				reflect.ValueOf(dest[i]).Elem().SetFloat(reflect.ValueOf(col).Float())
			default:
				return fmt.Errorf("go-ora: column %d require type float", i)
			}
		case time.Time:
			if destTyp.Elem() == reflect.TypeOf(time.Time{}) {
				reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(col))
			} else {
				return fmt.Errorf("go-ora: column %d require type time.Time", i)
			}
		case []byte:
			if destTyp.Elem() == reflect.TypeOf([]byte{}) {
				reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(col))
			} else {
				return fmt.Errorf("go-ora: column %d require type []byte", i)
			}
		default:
			if reflect.TypeOf(col) == destTyp.Elem() {
				reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(col))
			} else {
				return fmt.Errorf("go-ora: column %d require type %v", i, reflect.TypeOf(col))
			}
		}
	}
	return nil
}

// Err return last error
func (dataSet *DataSet) Err() error {
	return dataSet.lasterr
}

// Next implement method need for sql.Rows interface
func (dataSet *DataSet) Next(dest []driver.Value) error {
	hasMoreRows := dataSet.parent.hasMoreRows()
	noOfRowsToFetch := len(dataSet.Rows) // dataSet.parent.noOfRowsToFetch()
	hasBLOB := dataSet.parent.hasBLOB()
	hasLONG := dataSet.parent.hasLONG()
	if !hasMoreRows && noOfRowsToFetch == 0 {
		return io.EOF
	}
	if dataSet.index > 0 && dataSet.index%len(dataSet.Rows) == 0 {
		if hasMoreRows {
			dataSet.Rows = make([]Row, 0, dataSet.parent.noOfRowsToFetch())
			err := dataSet.parent.fetch(dataSet)
			if err != nil {
				return err
			}
			noOfRowsToFetch = len(dataSet.Rows)
			hasMoreRows = dataSet.parent.hasMoreRows()
			dataSet.index = 0
			if !hasMoreRows && noOfRowsToFetch == 0 {
				return io.EOF
			}
		} else {
			return io.EOF
		}
	}
	if hasMoreRows && (hasBLOB || hasLONG) && dataSet.index == 0 {
		if err := dataSet.parent.fetch(dataSet); err != nil {
			return err
		}
	}
	if dataSet.index%noOfRowsToFetch < len(dataSet.Rows) {
		for x := 0; x < len(dataSet.Rows[dataSet.index%noOfRowsToFetch]); x++ {
			dest[x] = dataSet.Rows[dataSet.index%noOfRowsToFetch][x]
		}
		dataSet.index++
		return nil
	}
	return io.EOF
}

//func (dataSet *DataSet) NextRow(args... interface{}) error {
//	var values = make([]driver.Value, len(args))
//	err := dataSet.Next(values)
//	if err != nil {
//		return err
//	}
//	for index, arg := range args {
//		*arg = values[index]
//		//if val, ok := values[index].(t); !ok {
//		//
//		//}
//	}
//	return nil
//}

// Columns return a string array that represent columns names
func (dataSet *DataSet) Columns() []string {
	if len(dataSet.Cols) == 0 {
		return nil
	}
	ret := make([]string, len(dataSet.Cols))
	for x := 0; x < len(dataSet.Cols); x++ {
		ret[x] = dataSet.Cols[x].Name
	}
	return ret
}

func (dataSet DataSet) Trace(t trace.Tracer) {
	for r, row := range dataSet.Rows {
		if r > 25 {
			break
		}
		t.Printf("Row %d", r)
		for c, col := range dataSet.Cols {
			t.Printf("  %-20s: %v", col.Name, row[c])
		}
	}
}

// ColumnTypeDatabaseTypeName return Col DataType name
func (dataSet DataSet) ColumnTypeDatabaseTypeName(index int) string {
	return dataSet.Cols[index].DataType.String()
}

// ColumnTypeLength return length of column type
func (dataSet DataSet) ColumnTypeLength(index int) (length int64, ok bool) {
	switch dataSet.Cols[index].DataType {
	case NCHAR, CHAR:
		return int64(dataSet.Cols[index].MaxCharLen), true
	case NUMBER:
		return int64(dataSet.Cols[index].Precision), true
	}
	return int64(0), false

}

// ColumnTypeNullable return if column allow null or not
func (dataSet DataSet) ColumnTypeNullable(index int) (nullable, ok bool) {
	return dataSet.Cols[index].AllowNull, true
}
