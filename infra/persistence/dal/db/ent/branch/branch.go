// Code generated by ent, DO NOT EDIT.

package branch

import (
	"time"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the branch type in the database.
	Label = "branch"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldBranchID holds the string denoting the branch_id field in the database.
	FieldBranchID = "branch_id"
	// FieldTransID holds the string denoting the trans_id field in the database.
	FieldTransID = "trans_id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldState holds the string denoting the state field in the database.
	FieldState = "state"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldAction holds the string denoting the action field in the database.
	FieldAction = "action"
	// FieldCompensate holds the string denoting the compensate field in the database.
	FieldCompensate = "compensate"
	// FieldPayload holds the string denoting the payload field in the database.
	FieldPayload = "payload"
	// FieldActionDepend holds the string denoting the action_depend field in the database.
	FieldActionDepend = "action_depend"
	// FieldCompensateDepend holds the string denoting the compensate_depend field in the database.
	FieldCompensateDepend = "compensate_depend"
	// FieldFinishedAt holds the string denoting the finished_at field in the database.
	FieldFinishedAt = "finished_at"
	// FieldIsDead holds the string denoting the is_dead field in the database.
	FieldIsDead = "is_dead"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldUpdatedBy holds the string denoting the updated_by field in the database.
	FieldUpdatedBy = "updated_by"
	// FieldCreatedBy holds the string denoting the created_by field in the database.
	FieldCreatedBy = "created_by"
	// Table holds the table name of the branch in the database.
	Table = "branches"
)

// Columns holds all SQL columns for branch fields.
var Columns = []string{
	FieldID,
	FieldBranchID,
	FieldTransID,
	FieldType,
	FieldState,
	FieldName,
	FieldAction,
	FieldCompensate,
	FieldPayload,
	FieldActionDepend,
	FieldCompensateDepend,
	FieldFinishedAt,
	FieldIsDead,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldUpdatedBy,
	FieldCreatedBy,
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
	// DefaultFinishedAt holds the default value on creation for the "finished_at" field.
	DefaultFinishedAt func() time.Time
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
)

// OrderOption defines the ordering options for the Branch queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByBranchID orders the results by the branch_id field.
func ByBranchID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBranchID, opts...).ToFunc()
}

// ByTransID orders the results by the trans_id field.
func ByTransID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTransID, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByState orders the results by the state field.
func ByState(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldState, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByAction orders the results by the action field.
func ByAction(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAction, opts...).ToFunc()
}

// ByCompensate orders the results by the compensate field.
func ByCompensate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCompensate, opts...).ToFunc()
}

// ByPayload orders the results by the payload field.
func ByPayload(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPayload, opts...).ToFunc()
}

// ByActionDepend orders the results by the action_depend field.
func ByActionDepend(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldActionDepend, opts...).ToFunc()
}

// ByCompensateDepend orders the results by the compensate_depend field.
func ByCompensateDepend(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCompensateDepend, opts...).ToFunc()
}

// ByFinishedAt orders the results by the finished_at field.
func ByFinishedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFinishedAt, opts...).ToFunc()
}

// ByIsDead orders the results by the is_dead field.
func ByIsDead(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsDead, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByUpdatedBy orders the results by the updated_by field.
func ByUpdatedBy(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedBy, opts...).ToFunc()
}

// ByCreatedBy orders the results by the created_by field.
func ByCreatedBy(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedBy, opts...).ToFunc()
}
