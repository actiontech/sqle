package agentpb

import (
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go-sumtype:decl isQueryActionValue_Kind

func makeValue(value interface{}) (*QueryActionValue, error) {
	// In the future, we may decide to:
	// * dereference pointers;
	// * handle other types of the same kind (like `type String string`);
	// * handle typed nils (like `(*int)(nil)`).
	//
	// We should do it only once we have a clear use-case without (or with too ugly) workarounds.

	// avoid reflection for basic types
	var err error
	switch v := value.(type) {
	case nil:
		return &QueryActionValue{Kind: &QueryActionValue_Nil{Nil: true}}, nil

	case bool:
		return &QueryActionValue{Kind: &QueryActionValue_Bool{Bool: v}}, nil

	case int:
		return &QueryActionValue{Kind: &QueryActionValue_Int64{Int64: int64(v)}}, nil
	case int8:
		return &QueryActionValue{Kind: &QueryActionValue_Int64{Int64: int64(v)}}, nil
	case int16:
		return &QueryActionValue{Kind: &QueryActionValue_Int64{Int64: int64(v)}}, nil
	case int32:
		return &QueryActionValue{Kind: &QueryActionValue_Int64{Int64: int64(v)}}, nil
	case int64:
		return &QueryActionValue{Kind: &QueryActionValue_Int64{Int64: v}}, nil

	case uint:
		return &QueryActionValue{Kind: &QueryActionValue_Uint64{Uint64: uint64(v)}}, nil
	case uint8:
		return &QueryActionValue{Kind: &QueryActionValue_Uint64{Uint64: uint64(v)}}, nil
	case uint16:
		return &QueryActionValue{Kind: &QueryActionValue_Uint64{Uint64: uint64(v)}}, nil
	case uint32:
		return &QueryActionValue{Kind: &QueryActionValue_Uint64{Uint64: uint64(v)}}, nil
	case uint64:
		return &QueryActionValue{Kind: &QueryActionValue_Uint64{Uint64: v}}, nil

	case float32:
		return &QueryActionValue{Kind: &QueryActionValue_Double{Double: float64(v)}}, nil
	case float64:
		return &QueryActionValue{Kind: &QueryActionValue_Double{Double: v}}, nil

	case []byte:
		return &QueryActionValue{Kind: &QueryActionValue_Bytes{Bytes: v}}, nil
	case string:
		// We couldn't encode Go string (that can contain any byte sequence)
		// to protobuf string (that can contain only valid UTF-8 byte sequence).
		// See https://jira.percona.com/browse/SAAS-107.
		return &QueryActionValue{Kind: &QueryActionValue_Bytes{Bytes: []byte(v)}}, nil

	case time.Time:
		ts, err := ptypes.TimestampProto(v)
		if err != nil {
			return nil, errors.Wrap(err, "failed to handle time.Time")
		}
		return &QueryActionValue{Kind: &QueryActionValue_Timestamp{Timestamp: ts}}, nil
	case primitive.Timestamp:
		// https://docs.mongodb.com/manual/reference/bson-types/#timestamps
		// resolution is up to a second; cram I (ordinal) into nanoseconds
		var t time.Time
		if !v.IsZero() {
			t = time.Unix(int64(v.T), int64(v.I))
		}
		ts, err := ptypes.TimestampProto(t)
		if err != nil {
			return nil, errors.Wrap(err, "failed to handle MongoDB's primitive.Timestamp")
		}
		return &QueryActionValue{Kind: &QueryActionValue_Timestamp{Timestamp: ts}}, nil
	case primitive.DateTime:
		// https://docs.mongodb.com/manual/reference/bson-types/#date
		// resolution is up to a millisecond
		var t time.Time
		if v != 0 {
			t = v.Time()
		}
		ts, err := ptypes.TimestampProto(t)
		if err != nil {
			return nil, errors.Wrap(err, "failed to handle MongoDB's primitive.DateTime")
		}
		return &QueryActionValue{Kind: &QueryActionValue_Timestamp{Timestamp: ts}}, nil
	}

	// use reflection for slices (except []byte) and maps
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice:
		size := v.Len()
		s := make([]*QueryActionValue, size)
		for i := 0; i < size; i++ {
			s[i], err = makeValue(v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
		}
		return &QueryActionValue{Kind: &QueryActionValue_Slice{Slice: &QueryActionSlice{Slice: s}}}, nil

	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return nil, errors.Errorf("makeValue: unhandled map key type for %[1]v (%[1]T)", value)
		}

		iter := v.MapRange()
		m := make(map[string]*QueryActionValue, v.Len())
		for iter.Next() {
			m[iter.Key().String()], err = makeValue(iter.Value().Interface())
			if err != nil {
				return nil, err
			}
		}
		return &QueryActionValue{Kind: &QueryActionValue_Map{Map: &QueryActionMap{Map: m}}}, nil

	default:
		return nil, errors.Errorf("unhandled %[1]v (%[1]T)", value)
	}
}

// MarshalActionQuerySQLResult returns serialized form of query Action SQL result.
//
// It supports the following types:
//  * untyped nil;
//  * bool;
//  * int, int8, int16, int32, int64;
//  * uint, uint8, uint16, uint32, uint64;
//  * float32, float64;
//  * string, []byte;
//  * time.Time;
//  * []T for any T from above, including other slices and maps;
//  * map[string]T for any T from above, including other slices and maps.
func MarshalActionQuerySQLResult(columns []string, rows [][]interface{}) ([]byte, error) {
	res := QueryActionResult{
		Columns: columns,
		Rows:    make([]*QueryActionSlice, len(rows)),
	}

	var err error
	for i, row := range rows {
		if len(columns) != len(row) {
			return nil, errors.Errorf("invalid result: expected %d columns in row %d, got %d", len(columns), i, len(row))
		}

		s := QueryActionSlice{
			Slice: make([]*QueryActionValue, len(row)),
		}

		for column, value := range row {
			s.Slice[column], err = makeValue(value)
			if err != nil {
				return nil, err
			}
		}

		res.Rows[i] = &s
	}

	return proto.Marshal(&res)
}

// MarshalActionQueryDocsResult returns serialized form of query Action documents result.
//
// It supports the same types as MarshalActionQuerySQLResult plus:
// * MongoDB's primitive.DateTime and primitive.Timestamp are converted to time.Time.
func MarshalActionQueryDocsResult(docs []map[string]interface{}) ([]byte, error) {
	res := QueryActionResult{
		Docs: make([]*QueryActionMap, len(docs)),
	}

	var err error
	for i, row := range docs {
		m := QueryActionMap{
			Map: make(map[string]*QueryActionValue, len(row)),
		}

		for column, value := range row {
			m.Map[column], err = makeValue(value)
			if err != nil {
				return nil, err
			}
		}

		res.Docs[i] = &m
	}

	return proto.Marshal(&res)
}

func makeInterface(value *QueryActionValue) (interface{}, error) {
	var err error
	switch v := value.Kind.(type) {
	case *QueryActionValue_Nil:
		return nil, nil
	case *QueryActionValue_Bool:
		return v.Bool, nil
	case *QueryActionValue_Int64:
		return v.Int64, nil
	case *QueryActionValue_Uint64:
		return v.Uint64, nil
	case *QueryActionValue_Double:
		return v.Double, nil
	case *QueryActionValue_Bytes:
		// Convert to Go string just for better developer experience;
		// it can contain any byte sequence and not limited to UTF-8.
		// See https://jira.percona.com/browse/SAAS-107.
		return string(v.Bytes), nil
	case *QueryActionValue_Timestamp:
		t, err := ptypes.Timestamp(v.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to handle timestamp")
		}
		return t, nil

	case *QueryActionValue_Slice:
		s := make([]interface{}, len(v.Slice.Slice))
		for i, v := range v.Slice.Slice {
			s[i], err = makeInterface(v)
			if err != nil {
				return nil, err
			}
		}
		return s, nil

	case *QueryActionValue_Map:
		m := make(map[string]interface{}, len(v.Map.Map))
		for k, v := range v.Map.Map {
			m[k], err = makeInterface(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil

	default:
		return nil, errors.Errorf("unhandled %[1]v (%[1]T)", value)
	}
}

func unmarshalActionQuerySQLResult(columns []string, rows []*QueryActionSlice) ([]map[string]interface{}, error) {
	data := make([]map[string]interface{}, len(rows))

	var err error
	for i, s := range rows {
		if len(columns) != len(s.Slice) {
			return nil, errors.Errorf("invalid result: expected %d columns in row %d, got %d", len(columns), i, len(s.Slice))
		}

		row := make(map[string]interface{}, len(s.Slice))

		for si, sv := range s.Slice {
			row[columns[si]], err = makeInterface(sv)
			if err != nil {
				return nil, err
			}
		}

		data[i] = row
	}

	return data, nil
}

func unmarshalActionQueryDocsResult(docs []*QueryActionMap) ([]map[string]interface{}, error) {
	data := make([]map[string]interface{}, len(docs))

	var err error
	for i, m := range docs {
		row := make(map[string]interface{}, len(m.Map))

		for mk, mv := range m.Map {
			row[mk], err = makeInterface(mv)
			if err != nil {
				return nil, err
			}
		}

		data[i] = row
	}

	return data, nil
}

// UnmarshalActionQueryResult returns deserialized form of query Action result, both SQL and documents.
func UnmarshalActionQueryResult(b []byte) ([]map[string]interface{}, error) {
	var res QueryActionResult
	if err := proto.Unmarshal(b, &res); err != nil {
		return nil, err
	}

	lenColumns := len(res.Columns)
	lenRows := len(res.Rows)
	lenDocs := len(res.Docs)
	if (lenColumns != 0 || lenRows != 0) && lenDocs != 0 {
		return nil, errors.Errorf("invalid result: %d columns, %d rows, %d docs", lenColumns, lenRows, lenDocs)
	}

	if lenColumns > 0 {
		return unmarshalActionQuerySQLResult(res.Columns, res.Rows)
	}
	if lenDocs > 0 {
		return unmarshalActionQueryDocsResult(res.Docs)
	}
	return nil, nil
}
