// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"purchase/infra/persistence/dal/db/ent/parow"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// PARowCreate is the builder for creating a PARow entity.
type PARowCreate struct {
	config
	mutation *PARowMutation
	hooks    []Hook
}

// SetDocCode sets the "doc_code" field.
func (prc *PARowCreate) SetDocCode(s string) *PARowCreate {
	prc.mutation.SetDocCode(s)
	return prc
}

// SetOrderCode sets the "order_code" field.
func (prc *PARowCreate) SetOrderCode(s string) *PARowCreate {
	prc.mutation.SetOrderCode(s)
	return prc
}

// SetRowCurrency sets the "row_currency" field.
func (prc *PARowCreate) SetRowCurrency(s string) *PARowCreate {
	prc.mutation.SetRowCurrency(s)
	return prc
}

// SetTaxRatio sets the "tax_ratio" field.
func (prc *PARowCreate) SetTaxRatio(s string) *PARowCreate {
	prc.mutation.SetTaxRatio(s)
	return prc
}

// SetInitialAmount sets the "initial_amount" field.
func (prc *PARowCreate) SetInitialAmount(s string) *PARowCreate {
	prc.mutation.SetInitialAmount(s)
	return prc
}

// SetCreatedAt sets the "created_at" field.
func (prc *PARowCreate) SetCreatedAt(t time.Time) *PARowCreate {
	prc.mutation.SetCreatedAt(t)
	return prc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (prc *PARowCreate) SetNillableCreatedAt(t *time.Time) *PARowCreate {
	if t != nil {
		prc.SetCreatedAt(*t)
	}
	return prc
}

// SetUpdatedAt sets the "updated_at" field.
func (prc *PARowCreate) SetUpdatedAt(t time.Time) *PARowCreate {
	prc.mutation.SetUpdatedAt(t)
	return prc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (prc *PARowCreate) SetNillableUpdatedAt(t *time.Time) *PARowCreate {
	if t != nil {
		prc.SetUpdatedAt(*t)
	}
	return prc
}

// SetID sets the "id" field.
func (prc *PARowCreate) SetID(i int64) *PARowCreate {
	prc.mutation.SetID(i)
	return prc
}

// Mutation returns the PARowMutation object of the builder.
func (prc *PARowCreate) Mutation() *PARowMutation {
	return prc.mutation
}

// Save creates the PARow in the database.
func (prc *PARowCreate) Save(ctx context.Context) (*PARow, error) {
	prc.defaults()
	return withHooks(ctx, prc.sqlSave, prc.mutation, prc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (prc *PARowCreate) SaveX(ctx context.Context) *PARow {
	v, err := prc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (prc *PARowCreate) Exec(ctx context.Context) error {
	_, err := prc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (prc *PARowCreate) ExecX(ctx context.Context) {
	if err := prc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (prc *PARowCreate) defaults() {
	if _, ok := prc.mutation.CreatedAt(); !ok {
		v := parow.DefaultCreatedAt()
		prc.mutation.SetCreatedAt(v)
	}
	if _, ok := prc.mutation.UpdatedAt(); !ok {
		v := parow.DefaultUpdatedAt()
		prc.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (prc *PARowCreate) check() error {
	if _, ok := prc.mutation.DocCode(); !ok {
		return &ValidationError{Name: "doc_code", err: errors.New(`ent: missing required field "PARow.doc_code"`)}
	}
	if _, ok := prc.mutation.OrderCode(); !ok {
		return &ValidationError{Name: "order_code", err: errors.New(`ent: missing required field "PARow.order_code"`)}
	}
	if _, ok := prc.mutation.RowCurrency(); !ok {
		return &ValidationError{Name: "row_currency", err: errors.New(`ent: missing required field "PARow.row_currency"`)}
	}
	if _, ok := prc.mutation.TaxRatio(); !ok {
		return &ValidationError{Name: "tax_ratio", err: errors.New(`ent: missing required field "PARow.tax_ratio"`)}
	}
	if _, ok := prc.mutation.InitialAmount(); !ok {
		return &ValidationError{Name: "initial_amount", err: errors.New(`ent: missing required field "PARow.initial_amount"`)}
	}
	if _, ok := prc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "PARow.created_at"`)}
	}
	if _, ok := prc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "PARow.updated_at"`)}
	}
	return nil
}

func (prc *PARowCreate) sqlSave(ctx context.Context) (*PARow, error) {
	if err := prc.check(); err != nil {
		return nil, err
	}
	_node, _spec := prc.createSpec()
	if err := sqlgraph.CreateNode(ctx, prc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int64(id)
	}
	prc.mutation.id = &_node.ID
	prc.mutation.done = true
	return _node, nil
}

func (prc *PARowCreate) createSpec() (*PARow, *sqlgraph.CreateSpec) {
	var (
		_node = &PARow{config: prc.config}
		_spec = sqlgraph.NewCreateSpec(parow.Table, sqlgraph.NewFieldSpec(parow.FieldID, field.TypeInt64))
	)
	if id, ok := prc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := prc.mutation.DocCode(); ok {
		_spec.SetField(parow.FieldDocCode, field.TypeString, value)
		_node.DocCode = value
	}
	if value, ok := prc.mutation.OrderCode(); ok {
		_spec.SetField(parow.FieldOrderCode, field.TypeString, value)
		_node.OrderCode = value
	}
	if value, ok := prc.mutation.RowCurrency(); ok {
		_spec.SetField(parow.FieldRowCurrency, field.TypeString, value)
		_node.RowCurrency = value
	}
	if value, ok := prc.mutation.TaxRatio(); ok {
		_spec.SetField(parow.FieldTaxRatio, field.TypeString, value)
		_node.TaxRatio = value
	}
	if value, ok := prc.mutation.InitialAmount(); ok {
		_spec.SetField(parow.FieldInitialAmount, field.TypeString, value)
		_node.InitialAmount = value
	}
	if value, ok := prc.mutation.CreatedAt(); ok {
		_spec.SetField(parow.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := prc.mutation.UpdatedAt(); ok {
		_spec.SetField(parow.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	return _node, _spec
}

// PARowCreateBulk is the builder for creating many PARow entities in bulk.
type PARowCreateBulk struct {
	config
	err      error
	builders []*PARowCreate
}

// Save creates the PARow entities in the database.
func (prcb *PARowCreateBulk) Save(ctx context.Context) ([]*PARow, error) {
	if prcb.err != nil {
		return nil, prcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(prcb.builders))
	nodes := make([]*PARow, len(prcb.builders))
	mutators := make([]Mutator, len(prcb.builders))
	for i := range prcb.builders {
		func(i int, root context.Context) {
			builder := prcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PARowMutation)
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
					_, err = mutators[i+1].Mutate(root, prcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, prcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, prcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (prcb *PARowCreateBulk) SaveX(ctx context.Context) []*PARow {
	v, err := prcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (prcb *PARowCreateBulk) Exec(ctx context.Context) error {
	_, err := prcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (prcb *PARowCreateBulk) ExecX(ctx context.Context) {
	if err := prcb.Exec(ctx); err != nil {
		panic(err)
	}
}
