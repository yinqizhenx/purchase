// Code generated by ent, DO NOT EDIT.

package pahead

import (
	"time"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the pahead type in the database.
	Label = "pa_head"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCode holds the string denoting the code field in the database.
	FieldCode = "code"
	// FieldState holds the string denoting the state field in the database.
	FieldState = "state"
	// FieldPayAmount holds the string denoting the pay_amount field in the database.
	FieldPayAmount = "pay_amount"
	// FieldApplicant holds the string denoting the applicant field in the database.
	FieldApplicant = "applicant"
	// FieldDepartmentCode holds the string denoting the department_code field in the database.
	FieldDepartmentCode = "department_code"
	// FieldSupplierCode holds the string denoting the supplier_code field in the database.
	FieldSupplierCode = "supplier_code"
	// FieldIsAdv holds the string denoting the is_adv field in the database.
	FieldIsAdv = "is_adv"
	// FieldHasInvoice holds the string denoting the has_invoice field in the database.
	FieldHasInvoice = "has_invoice"
	// FieldRemark holds the string denoting the remark field in the database.
	FieldRemark = "remark"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// Table holds the table name of the pahead in the database.
	Table = "pa_heads"
)

// Columns holds all SQL columns for pahead fields.
var Columns = []string{
	FieldID,
	FieldCode,
	FieldState,
	FieldPayAmount,
	FieldApplicant,
	FieldDepartmentCode,
	FieldSupplierCode,
	FieldIsAdv,
	FieldHasInvoice,
	FieldRemark,
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

// OrderOption defines the ordering options for the PAHead queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCode orders the results by the code field.
func ByCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCode, opts...).ToFunc()
}

// ByState orders the results by the state field.
func ByState(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldState, opts...).ToFunc()
}

// ByPayAmount orders the results by the pay_amount field.
func ByPayAmount(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPayAmount, opts...).ToFunc()
}

// ByApplicant orders the results by the applicant field.
func ByApplicant(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldApplicant, opts...).ToFunc()
}

// ByDepartmentCode orders the results by the department_code field.
func ByDepartmentCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDepartmentCode, opts...).ToFunc()
}

// BySupplierCode orders the results by the supplier_code field.
func BySupplierCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSupplierCode, opts...).ToFunc()
}

// ByIsAdv orders the results by the is_adv field.
func ByIsAdv(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsAdv, opts...).ToFunc()
}

// ByHasInvoice orders the results by the has_invoice field.
func ByHasInvoice(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHasInvoice, opts...).ToFunc()
}

// ByRemark orders the results by the remark field.
func ByRemark(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRemark, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}
