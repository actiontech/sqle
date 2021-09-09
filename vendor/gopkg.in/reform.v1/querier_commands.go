package reform

import (
	"fmt"
	"strings"
)

func filteredColumnsAndValues(str Struct, columnsIn []string, isUpdate bool) (columns []string, values []interface{}, err error) {
	columnsSet := make(map[string]struct{}, len(columnsIn))
	for _, c := range columnsIn {
		columnsSet[c] = struct{}{}
	}

	// select columns from set and collect values
	view := str.View()
	allColumns := view.Columns()
	allValues := str.Values()
	columns = make([]string, 0, len(columnsSet))
	values = make([]interface{}, 0, len(columns))

	record, _ := str.(Record)
	var pk uint
	if record != nil {
		pk = view.(Table).PKColumnIndex()
	}

	for i, c := range allColumns {
		if _, ok := columnsSet[c]; ok {
			if isUpdate && record != nil && i == int(pk) {
				err = fmt.Errorf("reform: will not update PK column: %s", c)
				return
			}
			delete(columnsSet, c)
			columns = append(columns, c)
			values = append(values, allValues[i])
		}
	}

	// make error for extra columns
	if len(columnsSet) > 0 {
		columns = make([]string, 0, len(columnsSet))
		for c := range columnsSet {
			columns = append(columns, c)
		}
		// TODO make exported type for that error
		err = fmt.Errorf("reform: unexpected columns: %v", columns)
		return
	}

	return
}

func (q *Querier) insert(str Struct, columns []string, values []interface{}) error {
	for i, c := range columns {
		columns[i] = q.QuoteIdentifier(c)
	}
	placeholders := q.Placeholders(1, len(columns))

	view := str.View()
	record, _ := str.(Record)
	lastInsertIdMethod := q.LastInsertIdMethod()
	defaultValuesMethod := q.DefaultValuesMethod()

	var pk uint
	if record != nil {
		pk = view.(Table).PKColumnIndex()
	}

	// make query
	query := q.startQuery("INSERT") + " INTO " + q.QualifiedView(view)
	if len(columns) > 0 || defaultValuesMethod == EmptyLists {
		query += " (" + strings.Join(columns, ", ") + ")"
	}
	if record != nil && lastInsertIdMethod == OutputInserted {
		query += fmt.Sprintf(" OUTPUT INSERTED.%s", q.QuoteIdentifier(view.Columns()[pk]))
	}
	if len(placeholders) > 0 || defaultValuesMethod == EmptyLists {
		query += fmt.Sprintf(" VALUES (%s)", strings.Join(placeholders, ", "))
	} else {
		query += " DEFAULT VALUES"
	}
	if record != nil && lastInsertIdMethod == Returning {
		query += fmt.Sprintf(" RETURNING %s", q.QuoteIdentifier(view.Columns()[pk]))
	}

	switch lastInsertIdMethod {
	case LastInsertId:
		res, err := q.Exec(query, values...)
		if err != nil {
			return err
		}
		if record != nil && !record.HasPK() {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}

			// TODO optimize to avoid using reflection
			// https://github.com/go-reform/reform/issues/269
			// record.SetPK(id)
			SetPK(record, id) //nolint:staticcheck
		}
		return nil

	case Returning, OutputInserted:
		var err error
		if record != nil {
			err = q.QueryRow(query, values...).Scan(record.PKPointer())
		} else {
			_, err = q.Exec(query, values...)
		}
		return err

	default:
		panic("reform: Unhandled LastInsertIdMethod. Please report this bug.")
	}
}

func (q *Querier) beforeInsert(str Struct) error {
	if bi, ok := str.(BeforeInserter); ok {
		if err := bi.BeforeInsert(); err != nil {
			return err
		}
	}

	return nil
}

// Insert inserts a struct into SQL database table.
// If str implements BeforeInserter, it calls BeforeInsert() before doing so.
//
// It fills record's primary key field.
func (q *Querier) Insert(str Struct) error {
	if err := q.beforeInsert(str); err != nil {
		return err
	}

	view := str.View()
	values := str.Values()
	columns := view.Columns()
	record, _ := str.(Record)

	if record != nil {
		pk := view.(Table).PKColumnIndex()

		// cut primary key
		if !record.HasPK() {
			values = append(values[:pk], values[pk+1:]...)
			columns = append(columns[:pk], columns[pk+1:]...)
		}
	}

	return q.insert(str, columns, values)
}

// InsertColumns inserts a struct into SQL database table with specified columns.
// Other columns are omitted from generated INSERT statement.
// If str implements BeforeInserter, it calls BeforeInsert() before doing so.
//
// It fills record's primary key field.
func (q *Querier) InsertColumns(str Struct, columns ...string) error {
	if err := q.beforeInsert(str); err != nil {
		return err
	}

	columns, values, err := filteredColumnsAndValues(str, columns, false)
	if err != nil {
		return err
	}

	return q.insert(str, columns, values)
}

// InsertMulti inserts several structs into SQL database table with single query.
// If they implement BeforeInserter, it calls BeforeInsert() before doing so.
//
// All structs should belong to the same view/table.
// All records should either have or not have primary key set.
// It doesn't fill primary key fields.
// Given all these limitations, most users should use Querier.Insert in a loop, not this method.
func (q *Querier) InsertMulti(structs ...Struct) error {
	if len(structs) == 0 {
		return nil
	}

	// check that view is the same
	view := structs[0].View()
	for _, str := range structs {
		if str.View() != view {
			return fmt.Errorf("reform: different tables in InsertMulti: %s and %s", view.Name(), str.View().Name())
		}
	}

	var err error
	for _, str := range structs {
		if bi, ok := str.(BeforeInserter); ok {
			e := bi.BeforeInsert()
			if err == nil {
				err = e
			}
		}
	}
	if err != nil {
		return err
	}

	// check if all PK are present or all are absent
	record, _ := structs[0].(Record)
	if record != nil {
		for _, str := range structs {
			rec, _ := str.(Record)
			if record.HasPK() != rec.HasPK() {
				return fmt.Errorf("reform: PK in present in one struct and absent in other: first: %s, second: %s",
					record, rec)
			}
		}
	}

	columns := view.Columns()
	for i, c := range columns {
		columns[i] = q.QuoteIdentifier(c)
	}

	var pk uint
	if record != nil && !record.HasPK() {
		pk = view.(Table).PKColumnIndex()
		columns = append(columns[:pk], columns[pk+1:]...)
	}

	placeholders := q.Placeholders(1, len(columns)*len(structs))
	query := fmt.Sprintf("%s INTO %s (%s) VALUES ",
		q.startQuery("INSERT"),
		q.QualifiedView(view),
		strings.Join(columns, ", "),
	)
	for i := 0; i < len(structs); i++ {
		query += fmt.Sprintf("(%s), ", strings.Join(placeholders[len(columns)*i:len(columns)*(i+1)], ", "))
	}
	query = query[:len(query)-2] // cut last ", "

	values := make([]interface{}, 0, len(placeholders))
	for _, str := range structs {
		v := str.Values()
		if record != nil && !record.HasPK() {
			v = append(v[:pk], v[pk+1:]...)
		}
		values = append(values, v...)
	}

	_, err = q.Exec(query, values...)
	return err
}

func (q *Querier) update(str Struct, columns []string, values []interface{}, tail string, args ...interface{}) (uint, error) {
	for i, c := range columns {
		columns[i] = q.QuoteIdentifier(c)
	}
	placeholders := q.Placeholders(1, len(columns))

	p := make([]string, len(columns))
	for i, c := range columns {
		p[i] = c + " = " + placeholders[i]
	}
	table := str.View()
	query := fmt.Sprintf("%s %s SET %s %s",
		q.startQuery("UPDATE"),
		q.QualifiedView(table),
		strings.Join(p, ", "),
		tail,
	)

	args = append(values, args...)
	res, err := q.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return uint(ra), nil
}

func (q *Querier) beforeUpdate(str Struct) error {
	if bu, ok := str.(BeforeUpdater); ok {
		if err := bu.BeforeUpdate(); err != nil {
			return err
		}
	}

	return nil
}

// Update updates all columns of row specified by primary key in SQL database table with given record.
// If record implements BeforeUpdater, it calls BeforeUpdate() before doing so.
//
// Method returns ErrNoRows if no rows were updated.
// Method returns ErrNoPK if primary key is not set.
func (q *Querier) Update(record Record) error {
	if err := q.beforeUpdate(record); err != nil {
		return err
	}
	if !record.HasPK() {
		return ErrNoPK
	}

	table := record.Table()
	values := record.Values()
	columns := table.Columns()

	// cut primary key, make tail
	pk := table.PKColumnIndex()
	pkColumn := columns[pk]
	values = append(values[:pk], values[pk+1:]...)
	columns = append(columns[:pk], columns[pk+1:]...)
	tail := fmt.Sprintf("WHERE %s = %s", q.QuoteIdentifier(pkColumn), q.Placeholder(len(columns)+1))

	ra, err := q.update(record, columns, values, tail, record.PKValue())
	if ra > 1 {
		panic(fmt.Sprintf("reform: %d rows by UPDATE by primary key. Please report this bug.", ra))
	}
	if err == nil && ra == 0 {
		err = ErrNoRows
	}
	return err
}

// UpdateColumns updates specified columns of row specified by primary key in SQL database table with given record.
// Other columns are omitted from generated UPDATE statement.
// If record implements BeforeUpdater, it calls BeforeUpdate() before doing so.
//
// Method returns ErrNoRows if no rows were updated.
// Method returns ErrNoPK if primary key is not set.
func (q *Querier) UpdateColumns(record Record, columns ...string) error {
	if err := q.beforeUpdate(record); err != nil {
		return err
	}
	if !record.HasPK() {
		return ErrNoPK
	}

	columns, values, err := filteredColumnsAndValues(record, columns, true)
	if err != nil {
		return err
	}

	if len(values) == 0 {
		// TODO make exported type for that error
		return fmt.Errorf("reform: nothing to update")
	}

	// make tail
	table := record.Table()
	pkColumn := table.Columns()[table.PKColumnIndex()]
	tail := fmt.Sprintf("WHERE %s = %s", q.QuoteIdentifier(pkColumn), q.Placeholder(len(columns)+1))

	ra, err := q.update(record, columns, values, tail, record.PKValue())
	if ra > 1 {
		panic(fmt.Sprintf("reform: %d rows by UPDATE by primary key. Please report this bug.", ra))
	}
	if err == nil && ra == 0 {
		err = ErrNoRows
	}
	return err
}

// UpdateView updates specified columns of rows specified by tail and args in SQL database table with given struct,
// and returns a number of updated rows.
// Other columns are omitted from generated UPDATE statement.
// If struct implements BeforeUpdater, it calls BeforeUpdate() before doing so.
//
// Method never returns ErrNoRows.
func (q *Querier) UpdateView(str Struct, columns []string, tail string, args ...interface{}) (uint, error) {
	if err := q.beforeUpdate(str); err != nil {
		return 0, err
	}

	columns, values, err := filteredColumnsAndValues(str, columns, true)
	if err != nil {
		return 0, err
	}

	if len(values) == 0 {
		// TODO make exported type for that error
		return 0, fmt.Errorf("reform: nothing to update")
	}

	return q.update(str, columns, values, tail, args...)
}

// Save saves record in SQL database table.
// If primary key is set, it first calls Update and checks if row was affected (matched).
// If primary key is absent or no row was affected, it calls Insert. This allows to call Save with Record
// with primary key set.
func (q *Querier) Save(record Record) error {
	if record.HasPK() {
		err := q.Update(record)
		if err != ErrNoRows {
			return err
		}
	}

	return q.Insert(record)
}

// Delete deletes record from SQL database table by primary key.
//
// Method returns ErrNoRows if no rows were deleted.
// Method returns ErrNoPK if primary key is not set.
func (q *Querier) Delete(record Record) error {
	if !record.HasPK() {
		return ErrNoPK
	}

	table := record.Table()
	pk := table.PKColumnIndex()
	query := fmt.Sprintf("%s FROM %s WHERE %s = %s",
		q.startQuery("DELETE"),
		q.QualifiedView(table),
		q.QuoteIdentifier(table.Columns()[pk]),
		q.Placeholder(1),
	)

	res, err := q.Exec(query, record.PKValue())
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return ErrNoRows
	}
	if ra > 1 {
		panic(fmt.Sprintf("reform: %d rows by DELETE by primary key. Please report this bug.", ra))
	}
	return nil
}

// DeleteFrom deletes rows from view with tail and args and returns a number of deleted rows.
//
// Method never returns ErrNoRows.
func (q *Querier) DeleteFrom(view View, tail string, args ...interface{}) (uint, error) {
	query := fmt.Sprintf("%s FROM %s %s",
		q.startQuery("DELETE"),
		q.QualifiedView(view),
		tail,
	)

	res, err := q.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return uint(ra), nil
}
