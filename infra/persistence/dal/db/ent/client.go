// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"purchase/infra/persistence/dal/db/ent/migrate"

	"purchase/infra/persistence/dal/db/ent/asynctask"
	"purchase/infra/persistence/dal/db/ent/branch"
	"purchase/infra/persistence/dal/db/ent/pahead"
	"purchase/infra/persistence/dal/db/ent/parow"
	"purchase/infra/persistence/dal/db/ent/trans"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"

	stdsql "database/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// AsyncTask is the client for interacting with the AsyncTask builders.
	AsyncTask *AsyncTaskClient
	// Branch is the client for interacting with the Branch builders.
	Branch *BranchClient
	// PAHead is the client for interacting with the PAHead builders.
	PAHead *PAHeadClient
	// PARow is the client for interacting with the PARow builders.
	PARow *PARowClient
	// Trans is the client for interacting with the Trans builders.
	Trans *TransClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.AsyncTask = NewAsyncTaskClient(c.config)
	c.Branch = NewBranchClient(c.config)
	c.PAHead = NewPAHeadClient(c.config)
	c.PARow = NewPARowClient(c.config)
	c.Trans = NewTransClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:       ctx,
		config:    cfg,
		AsyncTask: NewAsyncTaskClient(cfg),
		Branch:    NewBranchClient(cfg),
		PAHead:    NewPAHeadClient(cfg),
		PARow:     NewPARowClient(cfg),
		Trans:     NewTransClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:       ctx,
		config:    cfg,
		AsyncTask: NewAsyncTaskClient(cfg),
		Branch:    NewBranchClient(cfg),
		PAHead:    NewPAHeadClient(cfg),
		PARow:     NewPARowClient(cfg),
		Trans:     NewTransClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		AsyncTask.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.AsyncTask.Use(hooks...)
	c.Branch.Use(hooks...)
	c.PAHead.Use(hooks...)
	c.PARow.Use(hooks...)
	c.Trans.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.AsyncTask.Intercept(interceptors...)
	c.Branch.Intercept(interceptors...)
	c.PAHead.Intercept(interceptors...)
	c.PARow.Intercept(interceptors...)
	c.Trans.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *AsyncTaskMutation:
		return c.AsyncTask.mutate(ctx, m)
	case *BranchMutation:
		return c.Branch.mutate(ctx, m)
	case *PAHeadMutation:
		return c.PAHead.mutate(ctx, m)
	case *PARowMutation:
		return c.PARow.mutate(ctx, m)
	case *TransMutation:
		return c.Trans.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// AsyncTaskClient is a client for the AsyncTask schema.
type AsyncTaskClient struct {
	config
}

// NewAsyncTaskClient returns a client for the AsyncTask from the given config.
func NewAsyncTaskClient(c config) *AsyncTaskClient {
	return &AsyncTaskClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `asynctask.Hooks(f(g(h())))`.
func (c *AsyncTaskClient) Use(hooks ...Hook) {
	c.hooks.AsyncTask = append(c.hooks.AsyncTask, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `asynctask.Intercept(f(g(h())))`.
func (c *AsyncTaskClient) Intercept(interceptors ...Interceptor) {
	c.inters.AsyncTask = append(c.inters.AsyncTask, interceptors...)
}

// Create returns a builder for creating a AsyncTask entity.
func (c *AsyncTaskClient) Create() *AsyncTaskCreate {
	mutation := newAsyncTaskMutation(c.config, OpCreate)
	return &AsyncTaskCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of AsyncTask entities.
func (c *AsyncTaskClient) CreateBulk(builders ...*AsyncTaskCreate) *AsyncTaskCreateBulk {
	return &AsyncTaskCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *AsyncTaskClient) MapCreateBulk(slice any, setFunc func(*AsyncTaskCreate, int)) *AsyncTaskCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &AsyncTaskCreateBulk{err: fmt.Errorf("calling to AsyncTaskClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*AsyncTaskCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &AsyncTaskCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for AsyncTask.
func (c *AsyncTaskClient) Update() *AsyncTaskUpdate {
	mutation := newAsyncTaskMutation(c.config, OpUpdate)
	return &AsyncTaskUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *AsyncTaskClient) UpdateOne(at *AsyncTask) *AsyncTaskUpdateOne {
	mutation := newAsyncTaskMutation(c.config, OpUpdateOne, withAsyncTask(at))
	return &AsyncTaskUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *AsyncTaskClient) UpdateOneID(id int64) *AsyncTaskUpdateOne {
	mutation := newAsyncTaskMutation(c.config, OpUpdateOne, withAsyncTaskID(id))
	return &AsyncTaskUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for AsyncTask.
func (c *AsyncTaskClient) Delete() *AsyncTaskDelete {
	mutation := newAsyncTaskMutation(c.config, OpDelete)
	return &AsyncTaskDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *AsyncTaskClient) DeleteOne(at *AsyncTask) *AsyncTaskDeleteOne {
	return c.DeleteOneID(at.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *AsyncTaskClient) DeleteOneID(id int64) *AsyncTaskDeleteOne {
	builder := c.Delete().Where(asynctask.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &AsyncTaskDeleteOne{builder}
}

// Query returns a query builder for AsyncTask.
func (c *AsyncTaskClient) Query() *AsyncTaskQuery {
	return &AsyncTaskQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeAsyncTask},
		inters: c.Interceptors(),
	}
}

// Get returns a AsyncTask entity by its id.
func (c *AsyncTaskClient) Get(ctx context.Context, id int64) (*AsyncTask, error) {
	return c.Query().Where(asynctask.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *AsyncTaskClient) GetX(ctx context.Context, id int64) *AsyncTask {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *AsyncTaskClient) Hooks() []Hook {
	return c.hooks.AsyncTask
}

// Interceptors returns the client interceptors.
func (c *AsyncTaskClient) Interceptors() []Interceptor {
	return c.inters.AsyncTask
}

func (c *AsyncTaskClient) mutate(ctx context.Context, m *AsyncTaskMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&AsyncTaskCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&AsyncTaskUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&AsyncTaskUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&AsyncTaskDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown AsyncTask mutation op: %q", m.Op())
	}
}

// BranchClient is a client for the Branch schema.
type BranchClient struct {
	config
}

// NewBranchClient returns a client for the Branch from the given config.
func NewBranchClient(c config) *BranchClient {
	return &BranchClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `branch.Hooks(f(g(h())))`.
func (c *BranchClient) Use(hooks ...Hook) {
	c.hooks.Branch = append(c.hooks.Branch, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `branch.Intercept(f(g(h())))`.
func (c *BranchClient) Intercept(interceptors ...Interceptor) {
	c.inters.Branch = append(c.inters.Branch, interceptors...)
}

// Create returns a builder for creating a Branch entity.
func (c *BranchClient) Create() *BranchCreate {
	mutation := newBranchMutation(c.config, OpCreate)
	return &BranchCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Branch entities.
func (c *BranchClient) CreateBulk(builders ...*BranchCreate) *BranchCreateBulk {
	return &BranchCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *BranchClient) MapCreateBulk(slice any, setFunc func(*BranchCreate, int)) *BranchCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &BranchCreateBulk{err: fmt.Errorf("calling to BranchClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*BranchCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &BranchCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Branch.
func (c *BranchClient) Update() *BranchUpdate {
	mutation := newBranchMutation(c.config, OpUpdate)
	return &BranchUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *BranchClient) UpdateOne(b *Branch) *BranchUpdateOne {
	mutation := newBranchMutation(c.config, OpUpdateOne, withBranch(b))
	return &BranchUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *BranchClient) UpdateOneID(id int) *BranchUpdateOne {
	mutation := newBranchMutation(c.config, OpUpdateOne, withBranchID(id))
	return &BranchUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Branch.
func (c *BranchClient) Delete() *BranchDelete {
	mutation := newBranchMutation(c.config, OpDelete)
	return &BranchDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *BranchClient) DeleteOne(b *Branch) *BranchDeleteOne {
	return c.DeleteOneID(b.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *BranchClient) DeleteOneID(id int) *BranchDeleteOne {
	builder := c.Delete().Where(branch.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &BranchDeleteOne{builder}
}

// Query returns a query builder for Branch.
func (c *BranchClient) Query() *BranchQuery {
	return &BranchQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeBranch},
		inters: c.Interceptors(),
	}
}

// Get returns a Branch entity by its id.
func (c *BranchClient) Get(ctx context.Context, id int) (*Branch, error) {
	return c.Query().Where(branch.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *BranchClient) GetX(ctx context.Context, id int) *Branch {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *BranchClient) Hooks() []Hook {
	return c.hooks.Branch
}

// Interceptors returns the client interceptors.
func (c *BranchClient) Interceptors() []Interceptor {
	return c.inters.Branch
}

func (c *BranchClient) mutate(ctx context.Context, m *BranchMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&BranchCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&BranchUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&BranchUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&BranchDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Branch mutation op: %q", m.Op())
	}
}

// PAHeadClient is a client for the PAHead schema.
type PAHeadClient struct {
	config
}

// NewPAHeadClient returns a client for the PAHead from the given config.
func NewPAHeadClient(c config) *PAHeadClient {
	return &PAHeadClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `pahead.Hooks(f(g(h())))`.
func (c *PAHeadClient) Use(hooks ...Hook) {
	c.hooks.PAHead = append(c.hooks.PAHead, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `pahead.Intercept(f(g(h())))`.
func (c *PAHeadClient) Intercept(interceptors ...Interceptor) {
	c.inters.PAHead = append(c.inters.PAHead, interceptors...)
}

// Create returns a builder for creating a PAHead entity.
func (c *PAHeadClient) Create() *PAHeadCreate {
	mutation := newPAHeadMutation(c.config, OpCreate)
	return &PAHeadCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of PAHead entities.
func (c *PAHeadClient) CreateBulk(builders ...*PAHeadCreate) *PAHeadCreateBulk {
	return &PAHeadCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *PAHeadClient) MapCreateBulk(slice any, setFunc func(*PAHeadCreate, int)) *PAHeadCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &PAHeadCreateBulk{err: fmt.Errorf("calling to PAHeadClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*PAHeadCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &PAHeadCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for PAHead.
func (c *PAHeadClient) Update() *PAHeadUpdate {
	mutation := newPAHeadMutation(c.config, OpUpdate)
	return &PAHeadUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PAHeadClient) UpdateOne(ph *PAHead) *PAHeadUpdateOne {
	mutation := newPAHeadMutation(c.config, OpUpdateOne, withPAHead(ph))
	return &PAHeadUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *PAHeadClient) UpdateOneID(id int64) *PAHeadUpdateOne {
	mutation := newPAHeadMutation(c.config, OpUpdateOne, withPAHeadID(id))
	return &PAHeadUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for PAHead.
func (c *PAHeadClient) Delete() *PAHeadDelete {
	mutation := newPAHeadMutation(c.config, OpDelete)
	return &PAHeadDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *PAHeadClient) DeleteOne(ph *PAHead) *PAHeadDeleteOne {
	return c.DeleteOneID(ph.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *PAHeadClient) DeleteOneID(id int64) *PAHeadDeleteOne {
	builder := c.Delete().Where(pahead.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PAHeadDeleteOne{builder}
}

// Query returns a query builder for PAHead.
func (c *PAHeadClient) Query() *PAHeadQuery {
	return &PAHeadQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypePAHead},
		inters: c.Interceptors(),
	}
}

// Get returns a PAHead entity by its id.
func (c *PAHeadClient) Get(ctx context.Context, id int64) (*PAHead, error) {
	return c.Query().Where(pahead.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PAHeadClient) GetX(ctx context.Context, id int64) *PAHead {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *PAHeadClient) Hooks() []Hook {
	return c.hooks.PAHead
}

// Interceptors returns the client interceptors.
func (c *PAHeadClient) Interceptors() []Interceptor {
	return c.inters.PAHead
}

func (c *PAHeadClient) mutate(ctx context.Context, m *PAHeadMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&PAHeadCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&PAHeadUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&PAHeadUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&PAHeadDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown PAHead mutation op: %q", m.Op())
	}
}

// PARowClient is a client for the PARow schema.
type PARowClient struct {
	config
}

// NewPARowClient returns a client for the PARow from the given config.
func NewPARowClient(c config) *PARowClient {
	return &PARowClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `parow.Hooks(f(g(h())))`.
func (c *PARowClient) Use(hooks ...Hook) {
	c.hooks.PARow = append(c.hooks.PARow, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `parow.Intercept(f(g(h())))`.
func (c *PARowClient) Intercept(interceptors ...Interceptor) {
	c.inters.PARow = append(c.inters.PARow, interceptors...)
}

// Create returns a builder for creating a PARow entity.
func (c *PARowClient) Create() *PARowCreate {
	mutation := newPARowMutation(c.config, OpCreate)
	return &PARowCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of PARow entities.
func (c *PARowClient) CreateBulk(builders ...*PARowCreate) *PARowCreateBulk {
	return &PARowCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *PARowClient) MapCreateBulk(slice any, setFunc func(*PARowCreate, int)) *PARowCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &PARowCreateBulk{err: fmt.Errorf("calling to PARowClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*PARowCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &PARowCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for PARow.
func (c *PARowClient) Update() *PARowUpdate {
	mutation := newPARowMutation(c.config, OpUpdate)
	return &PARowUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PARowClient) UpdateOne(pr *PARow) *PARowUpdateOne {
	mutation := newPARowMutation(c.config, OpUpdateOne, withPARow(pr))
	return &PARowUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *PARowClient) UpdateOneID(id int64) *PARowUpdateOne {
	mutation := newPARowMutation(c.config, OpUpdateOne, withPARowID(id))
	return &PARowUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for PARow.
func (c *PARowClient) Delete() *PARowDelete {
	mutation := newPARowMutation(c.config, OpDelete)
	return &PARowDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *PARowClient) DeleteOne(pr *PARow) *PARowDeleteOne {
	return c.DeleteOneID(pr.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *PARowClient) DeleteOneID(id int64) *PARowDeleteOne {
	builder := c.Delete().Where(parow.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PARowDeleteOne{builder}
}

// Query returns a query builder for PARow.
func (c *PARowClient) Query() *PARowQuery {
	return &PARowQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypePARow},
		inters: c.Interceptors(),
	}
}

// Get returns a PARow entity by its id.
func (c *PARowClient) Get(ctx context.Context, id int64) (*PARow, error) {
	return c.Query().Where(parow.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PARowClient) GetX(ctx context.Context, id int64) *PARow {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *PARowClient) Hooks() []Hook {
	return c.hooks.PARow
}

// Interceptors returns the client interceptors.
func (c *PARowClient) Interceptors() []Interceptor {
	return c.inters.PARow
}

func (c *PARowClient) mutate(ctx context.Context, m *PARowMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&PARowCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&PARowUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&PARowUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&PARowDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown PARow mutation op: %q", m.Op())
	}
}

// TransClient is a client for the Trans schema.
type TransClient struct {
	config
}

// NewTransClient returns a client for the Trans from the given config.
func NewTransClient(c config) *TransClient {
	return &TransClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `trans.Hooks(f(g(h())))`.
func (c *TransClient) Use(hooks ...Hook) {
	c.hooks.Trans = append(c.hooks.Trans, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `trans.Intercept(f(g(h())))`.
func (c *TransClient) Intercept(interceptors ...Interceptor) {
	c.inters.Trans = append(c.inters.Trans, interceptors...)
}

// Create returns a builder for creating a Trans entity.
func (c *TransClient) Create() *TransCreate {
	mutation := newTransMutation(c.config, OpCreate)
	return &TransCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Trans entities.
func (c *TransClient) CreateBulk(builders ...*TransCreate) *TransCreateBulk {
	return &TransCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *TransClient) MapCreateBulk(slice any, setFunc func(*TransCreate, int)) *TransCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &TransCreateBulk{err: fmt.Errorf("calling to TransClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*TransCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &TransCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Trans.
func (c *TransClient) Update() *TransUpdate {
	mutation := newTransMutation(c.config, OpUpdate)
	return &TransUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TransClient) UpdateOne(t *Trans) *TransUpdateOne {
	mutation := newTransMutation(c.config, OpUpdateOne, withTrans(t))
	return &TransUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *TransClient) UpdateOneID(id int) *TransUpdateOne {
	mutation := newTransMutation(c.config, OpUpdateOne, withTransID(id))
	return &TransUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Trans.
func (c *TransClient) Delete() *TransDelete {
	mutation := newTransMutation(c.config, OpDelete)
	return &TransDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *TransClient) DeleteOne(t *Trans) *TransDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *TransClient) DeleteOneID(id int) *TransDeleteOne {
	builder := c.Delete().Where(trans.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TransDeleteOne{builder}
}

// Query returns a query builder for Trans.
func (c *TransClient) Query() *TransQuery {
	return &TransQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeTrans},
		inters: c.Interceptors(),
	}
}

// Get returns a Trans entity by its id.
func (c *TransClient) Get(ctx context.Context, id int) (*Trans, error) {
	return c.Query().Where(trans.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TransClient) GetX(ctx context.Context, id int) *Trans {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *TransClient) Hooks() []Hook {
	return c.hooks.Trans
}

// Interceptors returns the client interceptors.
func (c *TransClient) Interceptors() []Interceptor {
	return c.inters.Trans
}

func (c *TransClient) mutate(ctx context.Context, m *TransMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&TransCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&TransUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&TransUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&TransDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Trans mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		AsyncTask, Branch, PAHead, PARow, Trans []ent.Hook
	}
	inters struct {
		AsyncTask, Branch, PAHead, PARow, Trans []ent.Interceptor
	}
)

// ExecContext allows calling the underlying ExecContext method of the driver if it is supported by it.
// See, database/sql#DB.ExecContext for more information.
func (c *config) ExecContext(ctx context.Context, query string, args ...any) (stdsql.Result, error) {
	ex, ok := c.driver.(interface {
		ExecContext(context.Context, string, ...any) (stdsql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.ExecContext is not supported")
	}
	return ex.ExecContext(ctx, query, args...)
}

// QueryContext allows calling the underlying QueryContext method of the driver if it is supported by it.
// See, database/sql#DB.QueryContext for more information.
func (c *config) QueryContext(ctx context.Context, query string, args ...any) (*stdsql.Rows, error) {
	q, ok := c.driver.(interface {
		QueryContext(context.Context, string, ...any) (*stdsql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.QueryContext is not supported")
	}
	return q.QueryContext(ctx, query, args...)
}
