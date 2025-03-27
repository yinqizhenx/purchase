// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"purchase/infra/persistence/dal/db/ent/branch"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// BranchCreate is the builder for creating a Branch entity.
type BranchCreate struct {
	config
	mutation *BranchMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCode sets the "code" field.
func (bc *BranchCreate) SetCode(s string) *BranchCreate {
	bc.mutation.SetCode(s)
	return bc
}

// SetTransID sets the "trans_id" field.
func (bc *BranchCreate) SetTransID(i int) *BranchCreate {
	bc.mutation.SetTransID(i)
	return bc
}

// SetType sets the "type" field.
func (bc *BranchCreate) SetType(s string) *BranchCreate {
	bc.mutation.SetType(s)
	return bc
}

// SetState sets the "state" field.
func (bc *BranchCreate) SetState(s string) *BranchCreate {
	bc.mutation.SetState(s)
	return bc
}

// SetName sets the "name" field.
func (bc *BranchCreate) SetName(s string) *BranchCreate {
	bc.mutation.SetName(s)
	return bc
}

// SetAction sets the "action" field.
func (bc *BranchCreate) SetAction(s string) *BranchCreate {
	bc.mutation.SetAction(s)
	return bc
}

// SetCompensate sets the "compensate" field.
func (bc *BranchCreate) SetCompensate(s string) *BranchCreate {
	bc.mutation.SetCompensate(s)
	return bc
}

// SetPayload sets the "payload" field.
func (bc *BranchCreate) SetPayload(s string) *BranchCreate {
	bc.mutation.SetPayload(s)
	return bc
}

// SetActionDepend sets the "action_depend" field.
func (bc *BranchCreate) SetActionDepend(s string) *BranchCreate {
	bc.mutation.SetActionDepend(s)
	return bc
}

// SetCompensateDepend sets the "compensate_depend" field.
func (bc *BranchCreate) SetCompensateDepend(s string) *BranchCreate {
	bc.mutation.SetCompensateDepend(s)
	return bc
}

// SetFinishedAt sets the "finished_at" field.
func (bc *BranchCreate) SetFinishedAt(t time.Time) *BranchCreate {
	bc.mutation.SetFinishedAt(t)
	return bc
}

// SetNillableFinishedAt sets the "finished_at" field if the given value is not nil.
func (bc *BranchCreate) SetNillableFinishedAt(t *time.Time) *BranchCreate {
	if t != nil {
		bc.SetFinishedAt(*t)
	}
	return bc
}

// SetIsDead sets the "is_dead" field.
func (bc *BranchCreate) SetIsDead(b bool) *BranchCreate {
	bc.mutation.SetIsDead(b)
	return bc
}

// SetCreatedAt sets the "created_at" field.
func (bc *BranchCreate) SetCreatedAt(t time.Time) *BranchCreate {
	bc.mutation.SetCreatedAt(t)
	return bc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (bc *BranchCreate) SetNillableCreatedAt(t *time.Time) *BranchCreate {
	if t != nil {
		bc.SetCreatedAt(*t)
	}
	return bc
}

// SetUpdatedAt sets the "updated_at" field.
func (bc *BranchCreate) SetUpdatedAt(t time.Time) *BranchCreate {
	bc.mutation.SetUpdatedAt(t)
	return bc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (bc *BranchCreate) SetNillableUpdatedAt(t *time.Time) *BranchCreate {
	if t != nil {
		bc.SetUpdatedAt(*t)
	}
	return bc
}

// SetUpdatedBy sets the "updated_by" field.
func (bc *BranchCreate) SetUpdatedBy(s string) *BranchCreate {
	bc.mutation.SetUpdatedBy(s)
	return bc
}

// SetCreatedBy sets the "created_by" field.
func (bc *BranchCreate) SetCreatedBy(s string) *BranchCreate {
	bc.mutation.SetCreatedBy(s)
	return bc
}

// Mutation returns the BranchMutation object of the builder.
func (bc *BranchCreate) Mutation() *BranchMutation {
	return bc.mutation
}

// Save creates the Branch in the database.
func (bc *BranchCreate) Save(ctx context.Context) (*Branch, error) {
	bc.defaults()
	return withHooks(ctx, bc.sqlSave, bc.mutation, bc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (bc *BranchCreate) SaveX(ctx context.Context) *Branch {
	v, err := bc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bc *BranchCreate) Exec(ctx context.Context) error {
	_, err := bc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bc *BranchCreate) ExecX(ctx context.Context) {
	if err := bc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (bc *BranchCreate) defaults() {
	if _, ok := bc.mutation.FinishedAt(); !ok {
		v := branch.DefaultFinishedAt()
		bc.mutation.SetFinishedAt(v)
	}
	if _, ok := bc.mutation.CreatedAt(); !ok {
		v := branch.DefaultCreatedAt()
		bc.mutation.SetCreatedAt(v)
	}
	if _, ok := bc.mutation.UpdatedAt(); !ok {
		v := branch.DefaultUpdatedAt()
		bc.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bc *BranchCreate) check() error {
	if _, ok := bc.mutation.Code(); !ok {
		return &ValidationError{Name: "code", err: errors.New(`ent: missing required field "Branch.code"`)}
	}
	if _, ok := bc.mutation.TransID(); !ok {
		return &ValidationError{Name: "trans_id", err: errors.New(`ent: missing required field "Branch.trans_id"`)}
	}
	if _, ok := bc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "Branch.type"`)}
	}
	if _, ok := bc.mutation.State(); !ok {
		return &ValidationError{Name: "state", err: errors.New(`ent: missing required field "Branch.state"`)}
	}
	if _, ok := bc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Branch.name"`)}
	}
	if _, ok := bc.mutation.Action(); !ok {
		return &ValidationError{Name: "action", err: errors.New(`ent: missing required field "Branch.action"`)}
	}
	if _, ok := bc.mutation.Compensate(); !ok {
		return &ValidationError{Name: "compensate", err: errors.New(`ent: missing required field "Branch.compensate"`)}
	}
	if _, ok := bc.mutation.Payload(); !ok {
		return &ValidationError{Name: "payload", err: errors.New(`ent: missing required field "Branch.payload"`)}
	}
	if _, ok := bc.mutation.ActionDepend(); !ok {
		return &ValidationError{Name: "action_depend", err: errors.New(`ent: missing required field "Branch.action_depend"`)}
	}
	if _, ok := bc.mutation.CompensateDepend(); !ok {
		return &ValidationError{Name: "compensate_depend", err: errors.New(`ent: missing required field "Branch.compensate_depend"`)}
	}
	if _, ok := bc.mutation.FinishedAt(); !ok {
		return &ValidationError{Name: "finished_at", err: errors.New(`ent: missing required field "Branch.finished_at"`)}
	}
	if _, ok := bc.mutation.IsDead(); !ok {
		return &ValidationError{Name: "is_dead", err: errors.New(`ent: missing required field "Branch.is_dead"`)}
	}
	if _, ok := bc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Branch.created_at"`)}
	}
	if _, ok := bc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Branch.updated_at"`)}
	}
	if _, ok := bc.mutation.UpdatedBy(); !ok {
		return &ValidationError{Name: "updated_by", err: errors.New(`ent: missing required field "Branch.updated_by"`)}
	}
	if _, ok := bc.mutation.CreatedBy(); !ok {
		return &ValidationError{Name: "created_by", err: errors.New(`ent: missing required field "Branch.created_by"`)}
	}
	return nil
}

func (bc *BranchCreate) sqlSave(ctx context.Context) (*Branch, error) {
	if err := bc.check(); err != nil {
		return nil, err
	}
	_node, _spec := bc.createSpec()
	if err := sqlgraph.CreateNode(ctx, bc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	bc.mutation.id = &_node.ID
	bc.mutation.done = true
	return _node, nil
}

func (bc *BranchCreate) createSpec() (*Branch, *sqlgraph.CreateSpec) {
	var (
		_node = &Branch{config: bc.config}
		_spec = sqlgraph.NewCreateSpec(branch.Table, sqlgraph.NewFieldSpec(branch.FieldID, field.TypeInt))
	)
	_spec.OnConflict = bc.conflict
	if value, ok := bc.mutation.Code(); ok {
		_spec.SetField(branch.FieldCode, field.TypeString, value)
		_node.Code = value
	}
	if value, ok := bc.mutation.TransID(); ok {
		_spec.SetField(branch.FieldTransID, field.TypeInt, value)
		_node.TransID = value
	}
	if value, ok := bc.mutation.GetType(); ok {
		_spec.SetField(branch.FieldType, field.TypeString, value)
		_node.Type = value
	}
	if value, ok := bc.mutation.State(); ok {
		_spec.SetField(branch.FieldState, field.TypeString, value)
		_node.State = value
	}
	if value, ok := bc.mutation.Name(); ok {
		_spec.SetField(branch.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := bc.mutation.Action(); ok {
		_spec.SetField(branch.FieldAction, field.TypeString, value)
		_node.Action = value
	}
	if value, ok := bc.mutation.Compensate(); ok {
		_spec.SetField(branch.FieldCompensate, field.TypeString, value)
		_node.Compensate = value
	}
	if value, ok := bc.mutation.Payload(); ok {
		_spec.SetField(branch.FieldPayload, field.TypeString, value)
		_node.Payload = value
	}
	if value, ok := bc.mutation.ActionDepend(); ok {
		_spec.SetField(branch.FieldActionDepend, field.TypeString, value)
		_node.ActionDepend = value
	}
	if value, ok := bc.mutation.CompensateDepend(); ok {
		_spec.SetField(branch.FieldCompensateDepend, field.TypeString, value)
		_node.CompensateDepend = value
	}
	if value, ok := bc.mutation.FinishedAt(); ok {
		_spec.SetField(branch.FieldFinishedAt, field.TypeTime, value)
		_node.FinishedAt = value
	}
	if value, ok := bc.mutation.IsDead(); ok {
		_spec.SetField(branch.FieldIsDead, field.TypeBool, value)
		_node.IsDead = value
	}
	if value, ok := bc.mutation.CreatedAt(); ok {
		_spec.SetField(branch.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := bc.mutation.UpdatedAt(); ok {
		_spec.SetField(branch.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := bc.mutation.UpdatedBy(); ok {
		_spec.SetField(branch.FieldUpdatedBy, field.TypeString, value)
		_node.UpdatedBy = value
	}
	if value, ok := bc.mutation.CreatedBy(); ok {
		_spec.SetField(branch.FieldCreatedBy, field.TypeString, value)
		_node.CreatedBy = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Branch.Create().
//		SetCode(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.BranchUpsert) {
//			SetCode(v+v).
//		}).
//		Exec(ctx)
func (bc *BranchCreate) OnConflict(opts ...sql.ConflictOption) *BranchUpsertOne {
	bc.conflict = opts
	return &BranchUpsertOne{
		create: bc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Branch.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (bc *BranchCreate) OnConflictColumns(columns ...string) *BranchUpsertOne {
	bc.conflict = append(bc.conflict, sql.ConflictColumns(columns...))
	return &BranchUpsertOne{
		create: bc,
	}
}

type (
	// BranchUpsertOne is the builder for "upsert"-ing
	//  one Branch node.
	BranchUpsertOne struct {
		create *BranchCreate
	}

	// BranchUpsert is the "OnConflict" setter.
	BranchUpsert struct {
		*sql.UpdateSet
	}
)

// SetCode sets the "code" field.
func (u *BranchUpsert) SetCode(v string) *BranchUpsert {
	u.Set(branch.FieldCode, v)
	return u
}

// UpdateCode sets the "code" field to the value that was provided on create.
func (u *BranchUpsert) UpdateCode() *BranchUpsert {
	u.SetExcluded(branch.FieldCode)
	return u
}

// SetTransID sets the "trans_id" field.
func (u *BranchUpsert) SetTransID(v int) *BranchUpsert {
	u.Set(branch.FieldTransID, v)
	return u
}

// UpdateTransID sets the "trans_id" field to the value that was provided on create.
func (u *BranchUpsert) UpdateTransID() *BranchUpsert {
	u.SetExcluded(branch.FieldTransID)
	return u
}

// AddTransID adds v to the "trans_id" field.
func (u *BranchUpsert) AddTransID(v int) *BranchUpsert {
	u.Add(branch.FieldTransID, v)
	return u
}

// SetType sets the "type" field.
func (u *BranchUpsert) SetType(v string) *BranchUpsert {
	u.Set(branch.FieldType, v)
	return u
}

// UpdateType sets the "type" field to the value that was provided on create.
func (u *BranchUpsert) UpdateType() *BranchUpsert {
	u.SetExcluded(branch.FieldType)
	return u
}

// SetState sets the "state" field.
func (u *BranchUpsert) SetState(v string) *BranchUpsert {
	u.Set(branch.FieldState, v)
	return u
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *BranchUpsert) UpdateState() *BranchUpsert {
	u.SetExcluded(branch.FieldState)
	return u
}

// SetName sets the "name" field.
func (u *BranchUpsert) SetName(v string) *BranchUpsert {
	u.Set(branch.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *BranchUpsert) UpdateName() *BranchUpsert {
	u.SetExcluded(branch.FieldName)
	return u
}

// SetAction sets the "action" field.
func (u *BranchUpsert) SetAction(v string) *BranchUpsert {
	u.Set(branch.FieldAction, v)
	return u
}

// UpdateAction sets the "action" field to the value that was provided on create.
func (u *BranchUpsert) UpdateAction() *BranchUpsert {
	u.SetExcluded(branch.FieldAction)
	return u
}

// SetCompensate sets the "compensate" field.
func (u *BranchUpsert) SetCompensate(v string) *BranchUpsert {
	u.Set(branch.FieldCompensate, v)
	return u
}

// UpdateCompensate sets the "compensate" field to the value that was provided on create.
func (u *BranchUpsert) UpdateCompensate() *BranchUpsert {
	u.SetExcluded(branch.FieldCompensate)
	return u
}

// SetPayload sets the "payload" field.
func (u *BranchUpsert) SetPayload(v string) *BranchUpsert {
	u.Set(branch.FieldPayload, v)
	return u
}

// UpdatePayload sets the "payload" field to the value that was provided on create.
func (u *BranchUpsert) UpdatePayload() *BranchUpsert {
	u.SetExcluded(branch.FieldPayload)
	return u
}

// SetActionDepend sets the "action_depend" field.
func (u *BranchUpsert) SetActionDepend(v string) *BranchUpsert {
	u.Set(branch.FieldActionDepend, v)
	return u
}

// UpdateActionDepend sets the "action_depend" field to the value that was provided on create.
func (u *BranchUpsert) UpdateActionDepend() *BranchUpsert {
	u.SetExcluded(branch.FieldActionDepend)
	return u
}

// SetCompensateDepend sets the "compensate_depend" field.
func (u *BranchUpsert) SetCompensateDepend(v string) *BranchUpsert {
	u.Set(branch.FieldCompensateDepend, v)
	return u
}

// UpdateCompensateDepend sets the "compensate_depend" field to the value that was provided on create.
func (u *BranchUpsert) UpdateCompensateDepend() *BranchUpsert {
	u.SetExcluded(branch.FieldCompensateDepend)
	return u
}

// SetFinishedAt sets the "finished_at" field.
func (u *BranchUpsert) SetFinishedAt(v time.Time) *BranchUpsert {
	u.Set(branch.FieldFinishedAt, v)
	return u
}

// UpdateFinishedAt sets the "finished_at" field to the value that was provided on create.
func (u *BranchUpsert) UpdateFinishedAt() *BranchUpsert {
	u.SetExcluded(branch.FieldFinishedAt)
	return u
}

// SetIsDead sets the "is_dead" field.
func (u *BranchUpsert) SetIsDead(v bool) *BranchUpsert {
	u.Set(branch.FieldIsDead, v)
	return u
}

// UpdateIsDead sets the "is_dead" field to the value that was provided on create.
func (u *BranchUpsert) UpdateIsDead() *BranchUpsert {
	u.SetExcluded(branch.FieldIsDead)
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *BranchUpsert) SetCreatedAt(v time.Time) *BranchUpsert {
	u.Set(branch.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *BranchUpsert) UpdateCreatedAt() *BranchUpsert {
	u.SetExcluded(branch.FieldCreatedAt)
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *BranchUpsert) SetUpdatedAt(v time.Time) *BranchUpsert {
	u.Set(branch.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *BranchUpsert) UpdateUpdatedAt() *BranchUpsert {
	u.SetExcluded(branch.FieldUpdatedAt)
	return u
}

// SetUpdatedBy sets the "updated_by" field.
func (u *BranchUpsert) SetUpdatedBy(v string) *BranchUpsert {
	u.Set(branch.FieldUpdatedBy, v)
	return u
}

// UpdateUpdatedBy sets the "updated_by" field to the value that was provided on create.
func (u *BranchUpsert) UpdateUpdatedBy() *BranchUpsert {
	u.SetExcluded(branch.FieldUpdatedBy)
	return u
}

// SetCreatedBy sets the "created_by" field.
func (u *BranchUpsert) SetCreatedBy(v string) *BranchUpsert {
	u.Set(branch.FieldCreatedBy, v)
	return u
}

// UpdateCreatedBy sets the "created_by" field to the value that was provided on create.
func (u *BranchUpsert) UpdateCreatedBy() *BranchUpsert {
	u.SetExcluded(branch.FieldCreatedBy)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Branch.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *BranchUpsertOne) UpdateNewValues() *BranchUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Branch.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *BranchUpsertOne) Ignore() *BranchUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *BranchUpsertOne) DoNothing() *BranchUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the BranchCreate.OnConflict
// documentation for more info.
func (u *BranchUpsertOne) Update(set func(*BranchUpsert)) *BranchUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&BranchUpsert{UpdateSet: update})
	}))
	return u
}

// SetCode sets the "code" field.
func (u *BranchUpsertOne) SetCode(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetCode(v)
	})
}

// UpdateCode sets the "code" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateCode() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCode()
	})
}

// SetTransID sets the "trans_id" field.
func (u *BranchUpsertOne) SetTransID(v int) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetTransID(v)
	})
}

// AddTransID adds v to the "trans_id" field.
func (u *BranchUpsertOne) AddTransID(v int) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.AddTransID(v)
	})
}

// UpdateTransID sets the "trans_id" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateTransID() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateTransID()
	})
}

// SetType sets the "type" field.
func (u *BranchUpsertOne) SetType(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetType(v)
	})
}

// UpdateType sets the "type" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateType() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateType()
	})
}

// SetState sets the "state" field.
func (u *BranchUpsertOne) SetState(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetState(v)
	})
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateState() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateState()
	})
}

// SetName sets the "name" field.
func (u *BranchUpsertOne) SetName(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateName() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateName()
	})
}

// SetAction sets the "action" field.
func (u *BranchUpsertOne) SetAction(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetAction(v)
	})
}

// UpdateAction sets the "action" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateAction() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateAction()
	})
}

// SetCompensate sets the "compensate" field.
func (u *BranchUpsertOne) SetCompensate(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetCompensate(v)
	})
}

// UpdateCompensate sets the "compensate" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateCompensate() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCompensate()
	})
}

// SetPayload sets the "payload" field.
func (u *BranchUpsertOne) SetPayload(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetPayload(v)
	})
}

// UpdatePayload sets the "payload" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdatePayload() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdatePayload()
	})
}

// SetActionDepend sets the "action_depend" field.
func (u *BranchUpsertOne) SetActionDepend(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetActionDepend(v)
	})
}

// UpdateActionDepend sets the "action_depend" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateActionDepend() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateActionDepend()
	})
}

// SetCompensateDepend sets the "compensate_depend" field.
func (u *BranchUpsertOne) SetCompensateDepend(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetCompensateDepend(v)
	})
}

// UpdateCompensateDepend sets the "compensate_depend" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateCompensateDepend() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCompensateDepend()
	})
}

// SetFinishedAt sets the "finished_at" field.
func (u *BranchUpsertOne) SetFinishedAt(v time.Time) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetFinishedAt(v)
	})
}

// UpdateFinishedAt sets the "finished_at" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateFinishedAt() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateFinishedAt()
	})
}

// SetIsDead sets the "is_dead" field.
func (u *BranchUpsertOne) SetIsDead(v bool) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetIsDead(v)
	})
}

// UpdateIsDead sets the "is_dead" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateIsDead() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateIsDead()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *BranchUpsertOne) SetCreatedAt(v time.Time) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateCreatedAt() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *BranchUpsertOne) SetUpdatedAt(v time.Time) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateUpdatedAt() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetUpdatedBy sets the "updated_by" field.
func (u *BranchUpsertOne) SetUpdatedBy(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetUpdatedBy(v)
	})
}

// UpdateUpdatedBy sets the "updated_by" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateUpdatedBy() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateUpdatedBy()
	})
}

// SetCreatedBy sets the "created_by" field.
func (u *BranchUpsertOne) SetCreatedBy(v string) *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.SetCreatedBy(v)
	})
}

// UpdateCreatedBy sets the "created_by" field to the value that was provided on create.
func (u *BranchUpsertOne) UpdateCreatedBy() *BranchUpsertOne {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCreatedBy()
	})
}

// Exec executes the query.
func (u *BranchUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for BranchCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *BranchUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *BranchUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *BranchUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// BranchCreateBulk is the builder for creating many Branch entities in bulk.
type BranchCreateBulk struct {
	config
	err      error
	builders []*BranchCreate
	conflict []sql.ConflictOption
}

// Save creates the Branch entities in the database.
func (bcb *BranchCreateBulk) Save(ctx context.Context) ([]*Branch, error) {
	if bcb.err != nil {
		return nil, bcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(bcb.builders))
	nodes := make([]*Branch, len(bcb.builders))
	mutators := make([]Mutator, len(bcb.builders))
	for i := range bcb.builders {
		func(i int, root context.Context) {
			builder := bcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*BranchMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, bcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = bcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, bcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, bcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (bcb *BranchCreateBulk) SaveX(ctx context.Context) []*Branch {
	v, err := bcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (bcb *BranchCreateBulk) Exec(ctx context.Context) error {
	_, err := bcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bcb *BranchCreateBulk) ExecX(ctx context.Context) {
	if err := bcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Branch.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.BranchUpsert) {
//			SetCode(v+v).
//		}).
//		Exec(ctx)
func (bcb *BranchCreateBulk) OnConflict(opts ...sql.ConflictOption) *BranchUpsertBulk {
	bcb.conflict = opts
	return &BranchUpsertBulk{
		create: bcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Branch.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (bcb *BranchCreateBulk) OnConflictColumns(columns ...string) *BranchUpsertBulk {
	bcb.conflict = append(bcb.conflict, sql.ConflictColumns(columns...))
	return &BranchUpsertBulk{
		create: bcb,
	}
}

// BranchUpsertBulk is the builder for "upsert"-ing
// a bulk of Branch nodes.
type BranchUpsertBulk struct {
	create *BranchCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Branch.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *BranchUpsertBulk) UpdateNewValues() *BranchUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Branch.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *BranchUpsertBulk) Ignore() *BranchUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *BranchUpsertBulk) DoNothing() *BranchUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the BranchCreateBulk.OnConflict
// documentation for more info.
func (u *BranchUpsertBulk) Update(set func(*BranchUpsert)) *BranchUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&BranchUpsert{UpdateSet: update})
	}))
	return u
}

// SetCode sets the "code" field.
func (u *BranchUpsertBulk) SetCode(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetCode(v)
	})
}

// UpdateCode sets the "code" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateCode() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCode()
	})
}

// SetTransID sets the "trans_id" field.
func (u *BranchUpsertBulk) SetTransID(v int) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetTransID(v)
	})
}

// AddTransID adds v to the "trans_id" field.
func (u *BranchUpsertBulk) AddTransID(v int) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.AddTransID(v)
	})
}

// UpdateTransID sets the "trans_id" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateTransID() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateTransID()
	})
}

// SetType sets the "type" field.
func (u *BranchUpsertBulk) SetType(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetType(v)
	})
}

// UpdateType sets the "type" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateType() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateType()
	})
}

// SetState sets the "state" field.
func (u *BranchUpsertBulk) SetState(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetState(v)
	})
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateState() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateState()
	})
}

// SetName sets the "name" field.
func (u *BranchUpsertBulk) SetName(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateName() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateName()
	})
}

// SetAction sets the "action" field.
func (u *BranchUpsertBulk) SetAction(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetAction(v)
	})
}

// UpdateAction sets the "action" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateAction() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateAction()
	})
}

// SetCompensate sets the "compensate" field.
func (u *BranchUpsertBulk) SetCompensate(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetCompensate(v)
	})
}

// UpdateCompensate sets the "compensate" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateCompensate() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCompensate()
	})
}

// SetPayload sets the "payload" field.
func (u *BranchUpsertBulk) SetPayload(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetPayload(v)
	})
}

// UpdatePayload sets the "payload" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdatePayload() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdatePayload()
	})
}

// SetActionDepend sets the "action_depend" field.
func (u *BranchUpsertBulk) SetActionDepend(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetActionDepend(v)
	})
}

// UpdateActionDepend sets the "action_depend" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateActionDepend() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateActionDepend()
	})
}

// SetCompensateDepend sets the "compensate_depend" field.
func (u *BranchUpsertBulk) SetCompensateDepend(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetCompensateDepend(v)
	})
}

// UpdateCompensateDepend sets the "compensate_depend" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateCompensateDepend() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCompensateDepend()
	})
}

// SetFinishedAt sets the "finished_at" field.
func (u *BranchUpsertBulk) SetFinishedAt(v time.Time) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetFinishedAt(v)
	})
}

// UpdateFinishedAt sets the "finished_at" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateFinishedAt() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateFinishedAt()
	})
}

// SetIsDead sets the "is_dead" field.
func (u *BranchUpsertBulk) SetIsDead(v bool) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetIsDead(v)
	})
}

// UpdateIsDead sets the "is_dead" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateIsDead() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateIsDead()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *BranchUpsertBulk) SetCreatedAt(v time.Time) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateCreatedAt() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *BranchUpsertBulk) SetUpdatedAt(v time.Time) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateUpdatedAt() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetUpdatedBy sets the "updated_by" field.
func (u *BranchUpsertBulk) SetUpdatedBy(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetUpdatedBy(v)
	})
}

// UpdateUpdatedBy sets the "updated_by" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateUpdatedBy() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateUpdatedBy()
	})
}

// SetCreatedBy sets the "created_by" field.
func (u *BranchUpsertBulk) SetCreatedBy(v string) *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.SetCreatedBy(v)
	})
}

// UpdateCreatedBy sets the "created_by" field to the value that was provided on create.
func (u *BranchUpsertBulk) UpdateCreatedBy() *BranchUpsertBulk {
	return u.Update(func(s *BranchUpsert) {
		s.UpdateCreatedBy()
	})
}

// Exec executes the query.
func (u *BranchUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the BranchCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for BranchCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *BranchUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
