// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"purchase/infra/persistence/dal/db/ent/trans"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// Trans is the model entity for the Trans schema.
type Trans struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// TransID holds the value of the "trans_id" field.
	TransID string `json:"trans_id,omitempty"`
	// State holds the value of the "state" field.
	State string `json:"state,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// FinishedAt holds the value of the "finished_at" field.
	FinishedAt time.Time `json:"finished_at,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// UpdatedBy holds the value of the "updated_by" field.
	UpdatedBy string `json:"updated_by,omitempty"`
	// CreatedBy holds the value of the "created_by" field.
	CreatedBy    string `json:"created_by,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Trans) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case trans.FieldID:
			values[i] = new(sql.NullInt64)
		case trans.FieldTransID, trans.FieldState, trans.FieldName, trans.FieldUpdatedBy, trans.FieldCreatedBy:
			values[i] = new(sql.NullString)
		case trans.FieldFinishedAt, trans.FieldCreatedAt, trans.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Trans fields.
func (t *Trans) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case trans.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			t.ID = int(value.Int64)
		case trans.FieldTransID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field trans_id", values[i])
			} else if value.Valid {
				t.TransID = value.String
			}
		case trans.FieldState:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field state", values[i])
			} else if value.Valid {
				t.State = value.String
			}
		case trans.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				t.Name = value.String
			}
		case trans.FieldFinishedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field finished_at", values[i])
			} else if value.Valid {
				t.FinishedAt = value.Time
			}
		case trans.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				t.CreatedAt = value.Time
			}
		case trans.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				t.UpdatedAt = value.Time
			}
		case trans.FieldUpdatedBy:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field updated_by", values[i])
			} else if value.Valid {
				t.UpdatedBy = value.String
			}
		case trans.FieldCreatedBy:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field created_by", values[i])
			} else if value.Valid {
				t.CreatedBy = value.String
			}
		default:
			t.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Trans.
// This includes values selected through modifiers, order, etc.
func (t *Trans) Value(name string) (ent.Value, error) {
	return t.selectValues.Get(name)
}

// Update returns a builder for updating this Trans.
// Note that you need to call Trans.Unwrap() before calling this method if this Trans
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Trans) Update() *TransUpdateOne {
	return NewTransClient(t.config).UpdateOne(t)
}

// Unwrap unwraps the Trans entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (t *Trans) Unwrap() *Trans {
	_tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("ent: Trans is not a transactional entity")
	}
	t.config.driver = _tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Trans) String() string {
	var builder strings.Builder
	builder.WriteString("Trans(")
	builder.WriteString(fmt.Sprintf("id=%v, ", t.ID))
	builder.WriteString("trans_id=")
	builder.WriteString(t.TransID)
	builder.WriteString(", ")
	builder.WriteString("state=")
	builder.WriteString(t.State)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(t.Name)
	builder.WriteString(", ")
	builder.WriteString("finished_at=")
	builder.WriteString(t.FinishedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(t.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(t.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_by=")
	builder.WriteString(t.UpdatedBy)
	builder.WriteString(", ")
	builder.WriteString("created_by=")
	builder.WriteString(t.CreatedBy)
	builder.WriteByte(')')
	return builder.String()
}

// TransSlice is a parsable slice of Trans.
type TransSlice []*Trans
