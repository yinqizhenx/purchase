// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"purchase/infra/persistence/dal/db/ent/pahead"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// PAHeadCreate is the builder for creating a PAHead entity.
type PAHeadCreate struct {
	config
	mutation *PAHeadMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCode sets the "code" field.
func (phc *PAHeadCreate) SetCode(s string) *PAHeadCreate {
	phc.mutation.SetCode(s)
	return phc
}

// SetState sets the "state" field.
func (phc *PAHeadCreate) SetState(s string) *PAHeadCreate {
	phc.mutation.SetState(s)
	return phc
}

// SetPayAmount sets the "pay_amount" field.
func (phc *PAHeadCreate) SetPayAmount(s string) *PAHeadCreate {
	phc.mutation.SetPayAmount(s)
	return phc
}

// SetApplicant sets the "applicant" field.
func (phc *PAHeadCreate) SetApplicant(s string) *PAHeadCreate {
	phc.mutation.SetApplicant(s)
	return phc
}

// SetDepartmentCode sets the "department_code" field.
func (phc *PAHeadCreate) SetDepartmentCode(s string) *PAHeadCreate {
	phc.mutation.SetDepartmentCode(s)
	return phc
}

// SetNillableDepartmentCode sets the "department_code" field if the given value is not nil.
func (phc *PAHeadCreate) SetNillableDepartmentCode(s *string) *PAHeadCreate {
	if s != nil {
		phc.SetDepartmentCode(*s)
	}
	return phc
}

// SetSupplierCode sets the "supplier_code" field.
func (phc *PAHeadCreate) SetSupplierCode(s string) *PAHeadCreate {
	phc.mutation.SetSupplierCode(s)
	return phc
}

// SetIsAdv sets the "is_adv" field.
func (phc *PAHeadCreate) SetIsAdv(b bool) *PAHeadCreate {
	phc.mutation.SetIsAdv(b)
	return phc
}

// SetHasInvoice sets the "has_invoice" field.
func (phc *PAHeadCreate) SetHasInvoice(b bool) *PAHeadCreate {
	phc.mutation.SetHasInvoice(b)
	return phc
}

// SetRemark sets the "remark" field.
func (phc *PAHeadCreate) SetRemark(s string) *PAHeadCreate {
	phc.mutation.SetRemark(s)
	return phc
}

// SetCreatedAt sets the "created_at" field.
func (phc *PAHeadCreate) SetCreatedAt(t time.Time) *PAHeadCreate {
	phc.mutation.SetCreatedAt(t)
	return phc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (phc *PAHeadCreate) SetNillableCreatedAt(t *time.Time) *PAHeadCreate {
	if t != nil {
		phc.SetCreatedAt(*t)
	}
	return phc
}

// SetUpdatedAt sets the "updated_at" field.
func (phc *PAHeadCreate) SetUpdatedAt(t time.Time) *PAHeadCreate {
	phc.mutation.SetUpdatedAt(t)
	return phc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (phc *PAHeadCreate) SetNillableUpdatedAt(t *time.Time) *PAHeadCreate {
	if t != nil {
		phc.SetUpdatedAt(*t)
	}
	return phc
}

// SetID sets the "id" field.
func (phc *PAHeadCreate) SetID(i int64) *PAHeadCreate {
	phc.mutation.SetID(i)
	return phc
}

// Mutation returns the PAHeadMutation object of the builder.
func (phc *PAHeadCreate) Mutation() *PAHeadMutation {
	return phc.mutation
}

// Save creates the PAHead in the database.
func (phc *PAHeadCreate) Save(ctx context.Context) (*PAHead, error) {
	phc.defaults()
	return withHooks(ctx, phc.sqlSave, phc.mutation, phc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (phc *PAHeadCreate) SaveX(ctx context.Context) *PAHead {
	v, err := phc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (phc *PAHeadCreate) Exec(ctx context.Context) error {
	_, err := phc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (phc *PAHeadCreate) ExecX(ctx context.Context) {
	if err := phc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (phc *PAHeadCreate) defaults() {
	if _, ok := phc.mutation.CreatedAt(); !ok {
		v := pahead.DefaultCreatedAt()
		phc.mutation.SetCreatedAt(v)
	}
	if _, ok := phc.mutation.UpdatedAt(); !ok {
		v := pahead.DefaultUpdatedAt()
		phc.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (phc *PAHeadCreate) check() error {
	if _, ok := phc.mutation.Code(); !ok {
		return &ValidationError{Name: "code", err: errors.New(`ent: missing required field "PAHead.code"`)}
	}
	if _, ok := phc.mutation.State(); !ok {
		return &ValidationError{Name: "state", err: errors.New(`ent: missing required field "PAHead.state"`)}
	}
	if _, ok := phc.mutation.PayAmount(); !ok {
		return &ValidationError{Name: "pay_amount", err: errors.New(`ent: missing required field "PAHead.pay_amount"`)}
	}
	if _, ok := phc.mutation.Applicant(); !ok {
		return &ValidationError{Name: "applicant", err: errors.New(`ent: missing required field "PAHead.applicant"`)}
	}
	if _, ok := phc.mutation.SupplierCode(); !ok {
		return &ValidationError{Name: "supplier_code", err: errors.New(`ent: missing required field "PAHead.supplier_code"`)}
	}
	if _, ok := phc.mutation.IsAdv(); !ok {
		return &ValidationError{Name: "is_adv", err: errors.New(`ent: missing required field "PAHead.is_adv"`)}
	}
	if _, ok := phc.mutation.HasInvoice(); !ok {
		return &ValidationError{Name: "has_invoice", err: errors.New(`ent: missing required field "PAHead.has_invoice"`)}
	}
	if _, ok := phc.mutation.Remark(); !ok {
		return &ValidationError{Name: "remark", err: errors.New(`ent: missing required field "PAHead.remark"`)}
	}
	if _, ok := phc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "PAHead.created_at"`)}
	}
	if _, ok := phc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "PAHead.updated_at"`)}
	}
	return nil
}

func (phc *PAHeadCreate) sqlSave(ctx context.Context) (*PAHead, error) {
	if err := phc.check(); err != nil {
		return nil, err
	}
	_node, _spec := phc.createSpec()
	if err := sqlgraph.CreateNode(ctx, phc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int64(id)
	}
	phc.mutation.id = &_node.ID
	phc.mutation.done = true
	return _node, nil
}

func (phc *PAHeadCreate) createSpec() (*PAHead, *sqlgraph.CreateSpec) {
	var (
		_node = &PAHead{config: phc.config}
		_spec = sqlgraph.NewCreateSpec(pahead.Table, sqlgraph.NewFieldSpec(pahead.FieldID, field.TypeInt64))
	)
	_spec.OnConflict = phc.conflict
	if id, ok := phc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := phc.mutation.Code(); ok {
		_spec.SetField(pahead.FieldCode, field.TypeString, value)
		_node.Code = value
	}
	if value, ok := phc.mutation.State(); ok {
		_spec.SetField(pahead.FieldState, field.TypeString, value)
		_node.State = value
	}
	if value, ok := phc.mutation.PayAmount(); ok {
		_spec.SetField(pahead.FieldPayAmount, field.TypeString, value)
		_node.PayAmount = value
	}
	if value, ok := phc.mutation.Applicant(); ok {
		_spec.SetField(pahead.FieldApplicant, field.TypeString, value)
		_node.Applicant = value
	}
	if value, ok := phc.mutation.DepartmentCode(); ok {
		_spec.SetField(pahead.FieldDepartmentCode, field.TypeString, value)
		_node.DepartmentCode = value
	}
	if value, ok := phc.mutation.SupplierCode(); ok {
		_spec.SetField(pahead.FieldSupplierCode, field.TypeString, value)
		_node.SupplierCode = value
	}
	if value, ok := phc.mutation.IsAdv(); ok {
		_spec.SetField(pahead.FieldIsAdv, field.TypeBool, value)
		_node.IsAdv = value
	}
	if value, ok := phc.mutation.HasInvoice(); ok {
		_spec.SetField(pahead.FieldHasInvoice, field.TypeBool, value)
		_node.HasInvoice = value
	}
	if value, ok := phc.mutation.Remark(); ok {
		_spec.SetField(pahead.FieldRemark, field.TypeString, value)
		_node.Remark = value
	}
	if value, ok := phc.mutation.CreatedAt(); ok {
		_spec.SetField(pahead.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := phc.mutation.UpdatedAt(); ok {
		_spec.SetField(pahead.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.PAHead.Create().
//		SetCode(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PAHeadUpsert) {
//			SetCode(v+v).
//		}).
//		Exec(ctx)
func (phc *PAHeadCreate) OnConflict(opts ...sql.ConflictOption) *PAHeadUpsertOne {
	phc.conflict = opts
	return &PAHeadUpsertOne{
		create: phc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.PAHead.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (phc *PAHeadCreate) OnConflictColumns(columns ...string) *PAHeadUpsertOne {
	phc.conflict = append(phc.conflict, sql.ConflictColumns(columns...))
	return &PAHeadUpsertOne{
		create: phc,
	}
}

type (
	// PAHeadUpsertOne is the builder for "upsert"-ing
	//  one PAHead node.
	PAHeadUpsertOne struct {
		create *PAHeadCreate
	}

	// PAHeadUpsert is the "OnConflict" setter.
	PAHeadUpsert struct {
		*sql.UpdateSet
	}
)

// SetCode sets the "code" field.
func (u *PAHeadUpsert) SetCode(v string) *PAHeadUpsert {
	u.Set(pahead.FieldCode, v)
	return u
}

// UpdateCode sets the "code" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdateCode() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldCode)
	return u
}

// SetState sets the "state" field.
func (u *PAHeadUpsert) SetState(v string) *PAHeadUpsert {
	u.Set(pahead.FieldState, v)
	return u
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdateState() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldState)
	return u
}

// SetPayAmount sets the "pay_amount" field.
func (u *PAHeadUpsert) SetPayAmount(v string) *PAHeadUpsert {
	u.Set(pahead.FieldPayAmount, v)
	return u
}

// UpdatePayAmount sets the "pay_amount" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdatePayAmount() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldPayAmount)
	return u
}

// SetApplicant sets the "applicant" field.
func (u *PAHeadUpsert) SetApplicant(v string) *PAHeadUpsert {
	u.Set(pahead.FieldApplicant, v)
	return u
}

// UpdateApplicant sets the "applicant" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdateApplicant() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldApplicant)
	return u
}

// SetDepartmentCode sets the "department_code" field.
func (u *PAHeadUpsert) SetDepartmentCode(v string) *PAHeadUpsert {
	u.Set(pahead.FieldDepartmentCode, v)
	return u
}

// UpdateDepartmentCode sets the "department_code" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdateDepartmentCode() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldDepartmentCode)
	return u
}

// ClearDepartmentCode clears the value of the "department_code" field.
func (u *PAHeadUpsert) ClearDepartmentCode() *PAHeadUpsert {
	u.SetNull(pahead.FieldDepartmentCode)
	return u
}

// SetSupplierCode sets the "supplier_code" field.
func (u *PAHeadUpsert) SetSupplierCode(v string) *PAHeadUpsert {
	u.Set(pahead.FieldSupplierCode, v)
	return u
}

// UpdateSupplierCode sets the "supplier_code" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdateSupplierCode() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldSupplierCode)
	return u
}

// SetIsAdv sets the "is_adv" field.
func (u *PAHeadUpsert) SetIsAdv(v bool) *PAHeadUpsert {
	u.Set(pahead.FieldIsAdv, v)
	return u
}

// UpdateIsAdv sets the "is_adv" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdateIsAdv() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldIsAdv)
	return u
}

// SetHasInvoice sets the "has_invoice" field.
func (u *PAHeadUpsert) SetHasInvoice(v bool) *PAHeadUpsert {
	u.Set(pahead.FieldHasInvoice, v)
	return u
}

// UpdateHasInvoice sets the "has_invoice" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdateHasInvoice() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldHasInvoice)
	return u
}

// SetRemark sets the "remark" field.
func (u *PAHeadUpsert) SetRemark(v string) *PAHeadUpsert {
	u.Set(pahead.FieldRemark, v)
	return u
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdateRemark() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldRemark)
	return u
}

// SetCreatedAt sets the "created_at" field.
func (u *PAHeadUpsert) SetCreatedAt(v time.Time) *PAHeadUpsert {
	u.Set(pahead.FieldCreatedAt, v)
	return u
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdateCreatedAt() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldCreatedAt)
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *PAHeadUpsert) SetUpdatedAt(v time.Time) *PAHeadUpsert {
	u.Set(pahead.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *PAHeadUpsert) UpdateUpdatedAt() *PAHeadUpsert {
	u.SetExcluded(pahead.FieldUpdatedAt)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.PAHead.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(pahead.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *PAHeadUpsertOne) UpdateNewValues() *PAHeadUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(pahead.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.PAHead.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *PAHeadUpsertOne) Ignore() *PAHeadUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PAHeadUpsertOne) DoNothing() *PAHeadUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PAHeadCreate.OnConflict
// documentation for more info.
func (u *PAHeadUpsertOne) Update(set func(*PAHeadUpsert)) *PAHeadUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PAHeadUpsert{UpdateSet: update})
	}))
	return u
}

// SetCode sets the "code" field.
func (u *PAHeadUpsertOne) SetCode(v string) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetCode(v)
	})
}

// UpdateCode sets the "code" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdateCode() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateCode()
	})
}

// SetState sets the "state" field.
func (u *PAHeadUpsertOne) SetState(v string) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetState(v)
	})
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdateState() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateState()
	})
}

// SetPayAmount sets the "pay_amount" field.
func (u *PAHeadUpsertOne) SetPayAmount(v string) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetPayAmount(v)
	})
}

// UpdatePayAmount sets the "pay_amount" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdatePayAmount() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdatePayAmount()
	})
}

// SetApplicant sets the "applicant" field.
func (u *PAHeadUpsertOne) SetApplicant(v string) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetApplicant(v)
	})
}

// UpdateApplicant sets the "applicant" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdateApplicant() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateApplicant()
	})
}

// SetDepartmentCode sets the "department_code" field.
func (u *PAHeadUpsertOne) SetDepartmentCode(v string) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetDepartmentCode(v)
	})
}

// UpdateDepartmentCode sets the "department_code" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdateDepartmentCode() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateDepartmentCode()
	})
}

// ClearDepartmentCode clears the value of the "department_code" field.
func (u *PAHeadUpsertOne) ClearDepartmentCode() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.ClearDepartmentCode()
	})
}

// SetSupplierCode sets the "supplier_code" field.
func (u *PAHeadUpsertOne) SetSupplierCode(v string) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetSupplierCode(v)
	})
}

// UpdateSupplierCode sets the "supplier_code" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdateSupplierCode() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateSupplierCode()
	})
}

// SetIsAdv sets the "is_adv" field.
func (u *PAHeadUpsertOne) SetIsAdv(v bool) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetIsAdv(v)
	})
}

// UpdateIsAdv sets the "is_adv" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdateIsAdv() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateIsAdv()
	})
}

// SetHasInvoice sets the "has_invoice" field.
func (u *PAHeadUpsertOne) SetHasInvoice(v bool) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetHasInvoice(v)
	})
}

// UpdateHasInvoice sets the "has_invoice" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdateHasInvoice() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateHasInvoice()
	})
}

// SetRemark sets the "remark" field.
func (u *PAHeadUpsertOne) SetRemark(v string) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdateRemark() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateRemark()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *PAHeadUpsertOne) SetCreatedAt(v time.Time) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdateCreatedAt() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *PAHeadUpsertOne) SetUpdatedAt(v time.Time) *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *PAHeadUpsertOne) UpdateUpdatedAt() *PAHeadUpsertOne {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateUpdatedAt()
	})
}

// Exec executes the query.
func (u *PAHeadUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for PAHeadCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PAHeadUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *PAHeadUpsertOne) ID(ctx context.Context) (id int64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *PAHeadUpsertOne) IDX(ctx context.Context) int64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// PAHeadCreateBulk is the builder for creating many PAHead entities in bulk.
type PAHeadCreateBulk struct {
	config
	err      error
	builders []*PAHeadCreate
	conflict []sql.ConflictOption
}

// Save creates the PAHead entities in the database.
func (phcb *PAHeadCreateBulk) Save(ctx context.Context) ([]*PAHead, error) {
	if phcb.err != nil {
		return nil, phcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(phcb.builders))
	nodes := make([]*PAHead, len(phcb.builders))
	mutators := make([]Mutator, len(phcb.builders))
	for i := range phcb.builders {
		func(i int, root context.Context) {
			builder := phcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PAHeadMutation)
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
					_, err = mutators[i+1].Mutate(root, phcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = phcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, phcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int64(id)
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
		if _, err := mutators[0].Mutate(ctx, phcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (phcb *PAHeadCreateBulk) SaveX(ctx context.Context) []*PAHead {
	v, err := phcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (phcb *PAHeadCreateBulk) Exec(ctx context.Context) error {
	_, err := phcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (phcb *PAHeadCreateBulk) ExecX(ctx context.Context) {
	if err := phcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.PAHead.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.PAHeadUpsert) {
//			SetCode(v+v).
//		}).
//		Exec(ctx)
func (phcb *PAHeadCreateBulk) OnConflict(opts ...sql.ConflictOption) *PAHeadUpsertBulk {
	phcb.conflict = opts
	return &PAHeadUpsertBulk{
		create: phcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.PAHead.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (phcb *PAHeadCreateBulk) OnConflictColumns(columns ...string) *PAHeadUpsertBulk {
	phcb.conflict = append(phcb.conflict, sql.ConflictColumns(columns...))
	return &PAHeadUpsertBulk{
		create: phcb,
	}
}

// PAHeadUpsertBulk is the builder for "upsert"-ing
// a bulk of PAHead nodes.
type PAHeadUpsertBulk struct {
	create *PAHeadCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.PAHead.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(pahead.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *PAHeadUpsertBulk) UpdateNewValues() *PAHeadUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(pahead.FieldID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.PAHead.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *PAHeadUpsertBulk) Ignore() *PAHeadUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *PAHeadUpsertBulk) DoNothing() *PAHeadUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the PAHeadCreateBulk.OnConflict
// documentation for more info.
func (u *PAHeadUpsertBulk) Update(set func(*PAHeadUpsert)) *PAHeadUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&PAHeadUpsert{UpdateSet: update})
	}))
	return u
}

// SetCode sets the "code" field.
func (u *PAHeadUpsertBulk) SetCode(v string) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetCode(v)
	})
}

// UpdateCode sets the "code" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdateCode() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateCode()
	})
}

// SetState sets the "state" field.
func (u *PAHeadUpsertBulk) SetState(v string) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetState(v)
	})
}

// UpdateState sets the "state" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdateState() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateState()
	})
}

// SetPayAmount sets the "pay_amount" field.
func (u *PAHeadUpsertBulk) SetPayAmount(v string) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetPayAmount(v)
	})
}

// UpdatePayAmount sets the "pay_amount" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdatePayAmount() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdatePayAmount()
	})
}

// SetApplicant sets the "applicant" field.
func (u *PAHeadUpsertBulk) SetApplicant(v string) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetApplicant(v)
	})
}

// UpdateApplicant sets the "applicant" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdateApplicant() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateApplicant()
	})
}

// SetDepartmentCode sets the "department_code" field.
func (u *PAHeadUpsertBulk) SetDepartmentCode(v string) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetDepartmentCode(v)
	})
}

// UpdateDepartmentCode sets the "department_code" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdateDepartmentCode() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateDepartmentCode()
	})
}

// ClearDepartmentCode clears the value of the "department_code" field.
func (u *PAHeadUpsertBulk) ClearDepartmentCode() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.ClearDepartmentCode()
	})
}

// SetSupplierCode sets the "supplier_code" field.
func (u *PAHeadUpsertBulk) SetSupplierCode(v string) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetSupplierCode(v)
	})
}

// UpdateSupplierCode sets the "supplier_code" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdateSupplierCode() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateSupplierCode()
	})
}

// SetIsAdv sets the "is_adv" field.
func (u *PAHeadUpsertBulk) SetIsAdv(v bool) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetIsAdv(v)
	})
}

// UpdateIsAdv sets the "is_adv" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdateIsAdv() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateIsAdv()
	})
}

// SetHasInvoice sets the "has_invoice" field.
func (u *PAHeadUpsertBulk) SetHasInvoice(v bool) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetHasInvoice(v)
	})
}

// UpdateHasInvoice sets the "has_invoice" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdateHasInvoice() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateHasInvoice()
	})
}

// SetRemark sets the "remark" field.
func (u *PAHeadUpsertBulk) SetRemark(v string) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetRemark(v)
	})
}

// UpdateRemark sets the "remark" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdateRemark() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateRemark()
	})
}

// SetCreatedAt sets the "created_at" field.
func (u *PAHeadUpsertBulk) SetCreatedAt(v time.Time) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetCreatedAt(v)
	})
}

// UpdateCreatedAt sets the "created_at" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdateCreatedAt() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateCreatedAt()
	})
}

// SetUpdatedAt sets the "updated_at" field.
func (u *PAHeadUpsertBulk) SetUpdatedAt(v time.Time) *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *PAHeadUpsertBulk) UpdateUpdatedAt() *PAHeadUpsertBulk {
	return u.Update(func(s *PAHeadUpsert) {
		s.UpdateUpdatedAt()
	})
}

// Exec executes the query.
func (u *PAHeadUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the PAHeadCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for PAHeadCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *PAHeadUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
