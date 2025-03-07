// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"
	"purchase/infra/persistence/dal/db/ent/pahead"
	"purchase/infra/persistence/dal/db/ent/predicate"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// PAHeadQuery is the builder for querying PAHead entities.
type PAHeadQuery struct {
	config
	ctx        *QueryContext
	order      []pahead.OrderOption
	inters     []Interceptor
	predicates []predicate.PAHead
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PAHeadQuery builder.
func (phq *PAHeadQuery) Where(ps ...predicate.PAHead) *PAHeadQuery {
	phq.predicates = append(phq.predicates, ps...)
	return phq
}

// Limit the number of records to be returned by this query.
func (phq *PAHeadQuery) Limit(limit int) *PAHeadQuery {
	phq.ctx.Limit = &limit
	return phq
}

// Offset to start from.
func (phq *PAHeadQuery) Offset(offset int) *PAHeadQuery {
	phq.ctx.Offset = &offset
	return phq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (phq *PAHeadQuery) Unique(unique bool) *PAHeadQuery {
	phq.ctx.Unique = &unique
	return phq
}

// Order specifies how the records should be ordered.
func (phq *PAHeadQuery) Order(o ...pahead.OrderOption) *PAHeadQuery {
	phq.order = append(phq.order, o...)
	return phq
}

// First returns the first PAHead entity from the query.
// Returns a *NotFoundError when no PAHead was found.
func (phq *PAHeadQuery) First(ctx context.Context) (*PAHead, error) {
	nodes, err := phq.Limit(1).All(setContextOp(ctx, phq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{pahead.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (phq *PAHeadQuery) FirstX(ctx context.Context) *PAHead {
	node, err := phq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first PAHead ID from the query.
// Returns a *NotFoundError when no PAHead ID was found.
func (phq *PAHeadQuery) FirstID(ctx context.Context) (id int64, err error) {
	var ids []int64
	if ids, err = phq.Limit(1).IDs(setContextOp(ctx, phq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{pahead.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (phq *PAHeadQuery) FirstIDX(ctx context.Context) int64 {
	id, err := phq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single PAHead entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one PAHead entity is found.
// Returns a *NotFoundError when no PAHead entities are found.
func (phq *PAHeadQuery) Only(ctx context.Context) (*PAHead, error) {
	nodes, err := phq.Limit(2).All(setContextOp(ctx, phq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{pahead.Label}
	default:
		return nil, &NotSingularError{pahead.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (phq *PAHeadQuery) OnlyX(ctx context.Context) *PAHead {
	node, err := phq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only PAHead ID in the query.
// Returns a *NotSingularError when more than one PAHead ID is found.
// Returns a *NotFoundError when no entities are found.
func (phq *PAHeadQuery) OnlyID(ctx context.Context) (id int64, err error) {
	var ids []int64
	if ids, err = phq.Limit(2).IDs(setContextOp(ctx, phq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{pahead.Label}
	default:
		err = &NotSingularError{pahead.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (phq *PAHeadQuery) OnlyIDX(ctx context.Context) int64 {
	id, err := phq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PAHeads.
func (phq *PAHeadQuery) All(ctx context.Context) ([]*PAHead, error) {
	ctx = setContextOp(ctx, phq.ctx, "All")
	if err := phq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*PAHead, *PAHeadQuery]()
	return withInterceptors[[]*PAHead](ctx, phq, qr, phq.inters)
}

// AllX is like All, but panics if an error occurs.
func (phq *PAHeadQuery) AllX(ctx context.Context) []*PAHead {
	nodes, err := phq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of PAHead IDs.
func (phq *PAHeadQuery) IDs(ctx context.Context) (ids []int64, err error) {
	if phq.ctx.Unique == nil && phq.path != nil {
		phq.Unique(true)
	}
	ctx = setContextOp(ctx, phq.ctx, "IDs")
	if err = phq.Select(pahead.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (phq *PAHeadQuery) IDsX(ctx context.Context) []int64 {
	ids, err := phq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (phq *PAHeadQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, phq.ctx, "Count")
	if err := phq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, phq, querierCount[*PAHeadQuery](), phq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (phq *PAHeadQuery) CountX(ctx context.Context) int {
	count, err := phq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (phq *PAHeadQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, phq.ctx, "Exist")
	switch _, err := phq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (phq *PAHeadQuery) ExistX(ctx context.Context) bool {
	exist, err := phq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PAHeadQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (phq *PAHeadQuery) Clone() *PAHeadQuery {
	if phq == nil {
		return nil
	}
	return &PAHeadQuery{
		config:     phq.config,
		ctx:        phq.ctx.Clone(),
		order:      append([]pahead.OrderOption{}, phq.order...),
		inters:     append([]Interceptor{}, phq.inters...),
		predicates: append([]predicate.PAHead{}, phq.predicates...),
		// clone intermediate query.
		sql:  phq.sql.Clone(),
		path: phq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Code string `json:"code,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.PAHead.Query().
//		GroupBy(pahead.FieldCode).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (phq *PAHeadQuery) GroupBy(field string, fields ...string) *PAHeadGroupBy {
	phq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &PAHeadGroupBy{build: phq}
	grbuild.flds = &phq.ctx.Fields
	grbuild.label = pahead.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Code string `json:"code,omitempty"`
//	}
//
//	client.PAHead.Query().
//		Select(pahead.FieldCode).
//		Scan(ctx, &v)
func (phq *PAHeadQuery) Select(fields ...string) *PAHeadSelect {
	phq.ctx.Fields = append(phq.ctx.Fields, fields...)
	sbuild := &PAHeadSelect{PAHeadQuery: phq}
	sbuild.label = pahead.Label
	sbuild.flds, sbuild.scan = &phq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a PAHeadSelect configured with the given aggregations.
func (phq *PAHeadQuery) Aggregate(fns ...AggregateFunc) *PAHeadSelect {
	return phq.Select().Aggregate(fns...)
}

func (phq *PAHeadQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range phq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, phq); err != nil {
				return err
			}
		}
	}
	for _, f := range phq.ctx.Fields {
		if !pahead.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if phq.path != nil {
		prev, err := phq.path(ctx)
		if err != nil {
			return err
		}
		phq.sql = prev
	}
	return nil
}

func (phq *PAHeadQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*PAHead, error) {
	var (
		nodes = []*PAHead{}
		_spec = phq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*PAHead).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &PAHead{config: phq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(phq.modifiers) > 0 {
		_spec.Modifiers = phq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, phq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (phq *PAHeadQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := phq.querySpec()
	if len(phq.modifiers) > 0 {
		_spec.Modifiers = phq.modifiers
	}
	_spec.Node.Columns = phq.ctx.Fields
	if len(phq.ctx.Fields) > 0 {
		_spec.Unique = phq.ctx.Unique != nil && *phq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, phq.driver, _spec)
}

func (phq *PAHeadQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(pahead.Table, pahead.Columns, sqlgraph.NewFieldSpec(pahead.FieldID, field.TypeInt64))
	_spec.From = phq.sql
	if unique := phq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if phq.path != nil {
		_spec.Unique = true
	}
	if fields := phq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, pahead.FieldID)
		for i := range fields {
			if fields[i] != pahead.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := phq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := phq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := phq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := phq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (phq *PAHeadQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(phq.driver.Dialect())
	t1 := builder.Table(pahead.Table)
	columns := phq.ctx.Fields
	if len(columns) == 0 {
		columns = pahead.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if phq.sql != nil {
		selector = phq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if phq.ctx.Unique != nil && *phq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range phq.modifiers {
		m(selector)
	}
	for _, p := range phq.predicates {
		p(selector)
	}
	for _, p := range phq.order {
		p(selector)
	}
	if offset := phq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := phq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (phq *PAHeadQuery) ForUpdate(opts ...sql.LockOption) *PAHeadQuery {
	if phq.driver.Dialect() == dialect.Postgres {
		phq.Unique(false)
	}
	phq.modifiers = append(phq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return phq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (phq *PAHeadQuery) ForShare(opts ...sql.LockOption) *PAHeadQuery {
	if phq.driver.Dialect() == dialect.Postgres {
		phq.Unique(false)
	}
	phq.modifiers = append(phq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return phq
}

// PAHeadGroupBy is the group-by builder for PAHead entities.
type PAHeadGroupBy struct {
	selector
	build *PAHeadQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (phgb *PAHeadGroupBy) Aggregate(fns ...AggregateFunc) *PAHeadGroupBy {
	phgb.fns = append(phgb.fns, fns...)
	return phgb
}

// Scan applies the selector query and scans the result into the given value.
func (phgb *PAHeadGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, phgb.build.ctx, "GroupBy")
	if err := phgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PAHeadQuery, *PAHeadGroupBy](ctx, phgb.build, phgb, phgb.build.inters, v)
}

func (phgb *PAHeadGroupBy) sqlScan(ctx context.Context, root *PAHeadQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(phgb.fns))
	for _, fn := range phgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*phgb.flds)+len(phgb.fns))
		for _, f := range *phgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*phgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := phgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// PAHeadSelect is the builder for selecting fields of PAHead entities.
type PAHeadSelect struct {
	*PAHeadQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (phs *PAHeadSelect) Aggregate(fns ...AggregateFunc) *PAHeadSelect {
	phs.fns = append(phs.fns, fns...)
	return phs
}

// Scan applies the selector query and scans the result into the given value.
func (phs *PAHeadSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, phs.ctx, "Select")
	if err := phs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PAHeadQuery, *PAHeadSelect](ctx, phs.PAHeadQuery, phs, phs.inters, v)
}

func (phs *PAHeadSelect) sqlScan(ctx context.Context, root *PAHeadQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(phs.fns))
	for _, fn := range phs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*phs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := phs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
