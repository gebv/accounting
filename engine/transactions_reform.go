// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package engine

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type transactionTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("acca").
func (v *transactionTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("transactions").
func (v *transactionTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *transactionTableType) Columns() []string {
	return []string{"tx_id", "invoice_id", "key", "amount", "strategy", "provider", "provider_oper_id", "provider_oper_status", "provider_oper_url", "meta", "status", "next_status", "updated_at", "created_at"}
}

// NewStruct makes a new struct for that view or table.
func (v *transactionTableType) NewStruct() reform.Struct {
	return new(Transaction)
}

// NewRecord makes a new record for that table.
func (v *transactionTableType) NewRecord() reform.Record {
	return new(Transaction)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *transactionTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// TransactionTable represents transactions view or table in SQL database.
var TransactionTable = &transactionTableType{
	s: parse.StructInfo{Type: "Transaction", SQLSchema: "acca", SQLName: "transactions", Fields: []parse.FieldInfo{{Name: "TransactionID", Type: "int64", Column: "tx_id"}, {Name: "InvoiceID", Type: "int64", Column: "invoice_id"}, {Name: "Key", Type: "*string", Column: "key"}, {Name: "Amount", Type: "int64", Column: "amount"}, {Name: "Strategy", Type: "string", Column: "strategy"}, {Name: "Provider", Type: "Provider", Column: "provider"}, {Name: "ProviderOperID", Type: "*string", Column: "provider_oper_id"}, {Name: "ProviderOperStatus", Type: "*string", Column: "provider_oper_status"}, {Name: "ProviderOperUrl", Type: "*string", Column: "provider_oper_url"}, {Name: "Meta", Type: "*[]uint8", Column: "meta"}, {Name: "Status", Type: "TransactionStatus", Column: "status"}, {Name: "NextStatus", Type: "*TransactionStatus", Column: "next_status"}, {Name: "UpdatedAt", Type: "time.Time", Column: "updated_at"}, {Name: "CreatedAt", Type: "time.Time", Column: "created_at"}}, PKFieldIndex: 0},
	z: new(Transaction).Values(),
}

// String returns a string representation of this struct or record.
func (s Transaction) String() string {
	res := make([]string, 14)
	res[0] = "TransactionID: " + reform.Inspect(s.TransactionID, true)
	res[1] = "InvoiceID: " + reform.Inspect(s.InvoiceID, true)
	res[2] = "Key: " + reform.Inspect(s.Key, true)
	res[3] = "Amount: " + reform.Inspect(s.Amount, true)
	res[4] = "Strategy: " + reform.Inspect(s.Strategy, true)
	res[5] = "Provider: " + reform.Inspect(s.Provider, true)
	res[6] = "ProviderOperID: " + reform.Inspect(s.ProviderOperID, true)
	res[7] = "ProviderOperStatus: " + reform.Inspect(s.ProviderOperStatus, true)
	res[8] = "ProviderOperUrl: " + reform.Inspect(s.ProviderOperUrl, true)
	res[9] = "Meta: " + reform.Inspect(s.Meta, true)
	res[10] = "Status: " + reform.Inspect(s.Status, true)
	res[11] = "NextStatus: " + reform.Inspect(s.NextStatus, true)
	res[12] = "UpdatedAt: " + reform.Inspect(s.UpdatedAt, true)
	res[13] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Transaction) Values() []interface{} {
	return []interface{}{
		s.TransactionID,
		s.InvoiceID,
		s.Key,
		s.Amount,
		s.Strategy,
		s.Provider,
		s.ProviderOperID,
		s.ProviderOperStatus,
		s.ProviderOperUrl,
		s.Meta,
		s.Status,
		s.NextStatus,
		s.UpdatedAt,
		s.CreatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Transaction) Pointers() []interface{} {
	return []interface{}{
		&s.TransactionID,
		&s.InvoiceID,
		&s.Key,
		&s.Amount,
		&s.Strategy,
		&s.Provider,
		&s.ProviderOperID,
		&s.ProviderOperStatus,
		&s.ProviderOperUrl,
		&s.Meta,
		&s.Status,
		&s.NextStatus,
		&s.UpdatedAt,
		&s.CreatedAt,
	}
}

// View returns View object for that struct.
func (s *Transaction) View() reform.View {
	return TransactionTable
}

// Table returns Table object for that record.
func (s *Transaction) Table() reform.Table {
	return TransactionTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Transaction) PKValue() interface{} {
	return s.TransactionID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Transaction) PKPointer() interface{} {
	return &s.TransactionID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Transaction) HasPK() bool {
	return s.TransactionID != TransactionTable.z[TransactionTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *Transaction) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.TransactionID = int64(i64)
	} else {
		s.TransactionID = pk.(int64)
	}
}

// check interfaces
var (
	_ reform.View   = TransactionTable
	_ reform.Struct = (*Transaction)(nil)
	_ reform.Table  = TransactionTable
	_ reform.Record = (*Transaction)(nil)
	_ fmt.Stringer  = (*Transaction)(nil)
)

func init() {
	parse.AssertUpToDate(&TransactionTable.s, new(Transaction))
}
