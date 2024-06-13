// Code generated by ent, DO NOT EDIT.

package parow

import (
	"time"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the parow type in the database.
	Label = "pa_row"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldDocCode holds the string denoting the doc_code field in the database.
	FieldDocCode = "doc_code"
	// FieldOrderCode holds the string denoting the order_code field in the database.
	FieldOrderCode = "order_code"
	// FieldRowCurrency holds the string denoting the row_currency field in the database.
	FieldRowCurrency = "row_currency"
	// FieldTaxRatio holds the string denoting the tax_ratio field in the database.
	FieldTaxRatio = "tax_ratio"
	// FieldInitialAmount holds the string denoting the initial_amount field in the database.
	FieldInitialAmount = "initial_amount"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// Table holds the table name of the parow in the database.
	Table = "pa_rows"
)

// Columns holds all SQL columns for parow fields.
var Columns = []string{
	FieldID,
	FieldDocCode,
	FieldOrderCode,
	FieldRowCurrency,
	FieldTaxRatio,
	FieldInitialAmount,
	FieldCreatedAt,
	FieldUpdatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
)

// OrderOption defines the ordering options for the PARow queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByDocCode orders the results by the doc_code field.
func ByDocCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDocCode, opts...).ToFunc()
}

// ByOrderCode orders the results by the order_code field.
func ByOrderCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOrderCode, opts...).ToFunc()
}

// ByRowCurrency orders the results by the row_currency field.
func ByRowCurrency(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRowCurrency, opts...).ToFunc()
}

// ByTaxRatio orders the results by the tax_ratio field.
func ByTaxRatio(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTaxRatio, opts...).ToFunc()
}

// ByInitialAmount orders the results by the initial_amount field.
func ByInitialAmount(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldInitialAmount, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}
