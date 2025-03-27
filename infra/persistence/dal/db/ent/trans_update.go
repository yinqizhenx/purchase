// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"purchase/infra/persistence/dal/db/ent/predicate"
	"purchase/infra/persistence/dal/db/ent/trans"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// TransUpdate is the builder for updating Trans entities.
type TransUpdate struct {
	config
	hooks    []Hook
	mutation *TransMutation
}

// Where appends a list predicates to the TransUpdate builder.
func (tu *TransUpdate) Where(ps ...predicate.Trans) *TransUpdate {
	tu.mutation.Where(ps...)
	return tu
}

// SetState sets the "state" field.
func (tu *TransUpdate) SetState(s string) *TransUpdate {
	tu.mutation.SetState(s)
	return tu
}

// SetNillableState sets the "state" field if the given value is not nil.
func (tu *TransUpdate) SetNillableState(s *string) *TransUpdate {
	if s != nil {
		tu.SetState(*s)
	}
	return tu
}

// SetName sets the "name" field.
func (tu *TransUpdate) SetName(s string) *TransUpdate {
	tu.mutation.SetName(s)
	return tu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tu *TransUpdate) SetNillableName(s *string) *TransUpdate {
	if s != nil {
		tu.SetName(*s)
	}
	return tu
}

// SetFinishedAt sets the "finished_at" field.
func (tu *TransUpdate) SetFinishedAt(t time.Time) *TransUpdate {
	tu.mutation.SetFinishedAt(t)
	return tu
}

// SetNillableFinishedAt sets the "finished_at" field if the given value is not nil.
func (tu *TransUpdate) SetNillableFinishedAt(t *time.Time) *TransUpdate {
	if t != nil {
		tu.SetFinishedAt(*t)
	}
	return tu
}

// SetCreatedAt sets the "created_at" field.
func (tu *TransUpdate) SetCreatedAt(t time.Time) *TransUpdate {
	tu.mutation.SetCreatedAt(t)
	return tu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tu *TransUpdate) SetNillableCreatedAt(t *time.Time) *TransUpdate {
	if t != nil {
		tu.SetCreatedAt(*t)
	}
	return tu
}

// SetUpdatedAt sets the "updated_at" field.
func (tu *TransUpdate) SetUpdatedAt(t time.Time) *TransUpdate {
	tu.mutation.SetUpdatedAt(t)
	return tu
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (tu *TransUpdate) SetNillableUpdatedAt(t *time.Time) *TransUpdate {
	if t != nil {
		tu.SetUpdatedAt(*t)
	}
	return tu
}

// SetUpdatedBy sets the "updated_by" field.
func (tu *TransUpdate) SetUpdatedBy(s string) *TransUpdate {
	tu.mutation.SetUpdatedBy(s)
	return tu
}

// SetNillableUpdatedBy sets the "updated_by" field if the given value is not nil.
func (tu *TransUpdate) SetNillableUpdatedBy(s *string) *TransUpdate {
	if s != nil {
		tu.SetUpdatedBy(*s)
	}
	return tu
}

// SetCreatedBy sets the "created_by" field.
func (tu *TransUpdate) SetCreatedBy(s string) *TransUpdate {
	tu.mutation.SetCreatedBy(s)
	return tu
}

// SetNillableCreatedBy sets the "created_by" field if the given value is not nil.
func (tu *TransUpdate) SetNillableCreatedBy(s *string) *TransUpdate {
	if s != nil {
		tu.SetCreatedBy(*s)
	}
	return tu
}

// Mutation returns the TransMutation object of the builder.
func (tu *TransUpdate) Mutation() *TransMutation {
	return tu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tu *TransUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, tu.sqlSave, tu.mutation, tu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TransUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TransUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TransUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tu *TransUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(trans.Table, trans.Columns, sqlgraph.NewFieldSpec(trans.FieldID, field.TypeInt))
	if ps := tu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.State(); ok {
		_spec.SetField(trans.FieldState, field.TypeString, value)
	}
	if value, ok := tu.mutation.Name(); ok {
		_spec.SetField(trans.FieldName, field.TypeString, value)
	}
	if value, ok := tu.mutation.FinishedAt(); ok {
		_spec.SetField(trans.FieldFinishedAt, field.TypeTime, value)
	}
	if value, ok := tu.mutation.CreatedAt(); ok {
		_spec.SetField(trans.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := tu.mutation.UpdatedAt(); ok {
		_spec.SetField(trans.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := tu.mutation.UpdatedBy(); ok {
		_spec.SetField(trans.FieldUpdatedBy, field.TypeString, value)
	}
	if value, ok := tu.mutation.CreatedBy(); ok {
		_spec.SetField(trans.FieldCreatedBy, field.TypeString, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{trans.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tu.mutation.done = true
	return n, nil
}

// TransUpdateOne is the builder for updating a single Trans entity.
type TransUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TransMutation
}

// SetState sets the "state" field.
func (tuo *TransUpdateOne) SetState(s string) *TransUpdateOne {
	tuo.mutation.SetState(s)
	return tuo
}

// SetNillableState sets the "state" field if the given value is not nil.
func (tuo *TransUpdateOne) SetNillableState(s *string) *TransUpdateOne {
	if s != nil {
		tuo.SetState(*s)
	}
	return tuo
}

// SetName sets the "name" field.
func (tuo *TransUpdateOne) SetName(s string) *TransUpdateOne {
	tuo.mutation.SetName(s)
	return tuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tuo *TransUpdateOne) SetNillableName(s *string) *TransUpdateOne {
	if s != nil {
		tuo.SetName(*s)
	}
	return tuo
}

// SetFinishedAt sets the "finished_at" field.
func (tuo *TransUpdateOne) SetFinishedAt(t time.Time) *TransUpdateOne {
	tuo.mutation.SetFinishedAt(t)
	return tuo
}

// SetNillableFinishedAt sets the "finished_at" field if the given value is not nil.
func (tuo *TransUpdateOne) SetNillableFinishedAt(t *time.Time) *TransUpdateOne {
	if t != nil {
		tuo.SetFinishedAt(*t)
	}
	return tuo
}

// SetCreatedAt sets the "created_at" field.
func (tuo *TransUpdateOne) SetCreatedAt(t time.Time) *TransUpdateOne {
	tuo.mutation.SetCreatedAt(t)
	return tuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tuo *TransUpdateOne) SetNillableCreatedAt(t *time.Time) *TransUpdateOne {
	if t != nil {
		tuo.SetCreatedAt(*t)
	}
	return tuo
}

// SetUpdatedAt sets the "updated_at" field.
func (tuo *TransUpdateOne) SetUpdatedAt(t time.Time) *TransUpdateOne {
	tuo.mutation.SetUpdatedAt(t)
	return tuo
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (tuo *TransUpdateOne) SetNillableUpdatedAt(t *time.Time) *TransUpdateOne {
	if t != nil {
		tuo.SetUpdatedAt(*t)
	}
	return tuo
}

// SetUpdatedBy sets the "updated_by" field.
func (tuo *TransUpdateOne) SetUpdatedBy(s string) *TransUpdateOne {
	tuo.mutation.SetUpdatedBy(s)
	return tuo
}

// SetNillableUpdatedBy sets the "updated_by" field if the given value is not nil.
func (tuo *TransUpdateOne) SetNillableUpdatedBy(s *string) *TransUpdateOne {
	if s != nil {
		tuo.SetUpdatedBy(*s)
	}
	return tuo
}

// SetCreatedBy sets the "created_by" field.
func (tuo *TransUpdateOne) SetCreatedBy(s string) *TransUpdateOne {
	tuo.mutation.SetCreatedBy(s)
	return tuo
}

// SetNillableCreatedBy sets the "created_by" field if the given value is not nil.
func (tuo *TransUpdateOne) SetNillableCreatedBy(s *string) *TransUpdateOne {
	if s != nil {
		tuo.SetCreatedBy(*s)
	}
	return tuo
}

// Mutation returns the TransMutation object of the builder.
func (tuo *TransUpdateOne) Mutation() *TransMutation {
	return tuo.mutation
}

// Where appends a list predicates to the TransUpdate builder.
func (tuo *TransUpdateOne) Where(ps ...predicate.Trans) *TransUpdateOne {
	tuo.mutation.Where(ps...)
	return tuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tuo *TransUpdateOne) Select(field string, fields ...string) *TransUpdateOne {
	tuo.fields = append([]string{field}, fields...)
	return tuo
}

// Save executes the query and returns the updated Trans entity.
func (tuo *TransUpdateOne) Save(ctx context.Context) (*Trans, error) {
	return withHooks(ctx, tuo.sqlSave, tuo.mutation, tuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TransUpdateOne) SaveX(ctx context.Context) *Trans {
	node, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tuo *TransUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TransUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tuo *TransUpdateOne) sqlSave(ctx context.Context) (_node *Trans, err error) {
	_spec := sqlgraph.NewUpdateSpec(trans.Table, trans.Columns, sqlgraph.NewFieldSpec(trans.FieldID, field.TypeInt))
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Trans.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, trans.FieldID)
		for _, f := range fields {
			if !trans.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != trans.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tuo.mutation.State(); ok {
		_spec.SetField(trans.FieldState, field.TypeString, value)
	}
	if value, ok := tuo.mutation.Name(); ok {
		_spec.SetField(trans.FieldName, field.TypeString, value)
	}
	if value, ok := tuo.mutation.FinishedAt(); ok {
		_spec.SetField(trans.FieldFinishedAt, field.TypeTime, value)
	}
	if value, ok := tuo.mutation.CreatedAt(); ok {
		_spec.SetField(trans.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := tuo.mutation.UpdatedAt(); ok {
		_spec.SetField(trans.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := tuo.mutation.UpdatedBy(); ok {
		_spec.SetField(trans.FieldUpdatedBy, field.TypeString, value)
	}
	if value, ok := tuo.mutation.CreatedBy(); ok {
		_spec.SetField(trans.FieldCreatedBy, field.TypeString, value)
	}
	_node = &Trans{config: tuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{trans.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tuo.mutation.done = true
	return _node, nil
}
