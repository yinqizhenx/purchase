package db

type Filter struct {
	where   map[string]interface{}
	limiter *Limiter
	order   *Order
}

type Limiter struct {
	limit  int
	offset int
}

type Ordering string

const (
	Asc  Ordering = "asc"
	Desc Ordering = "desc"
)

type Order struct {
	col      string
	ordering Ordering
}

type Option func(filter *Filter)

func WithWhere(w map[string]interface{}) Option {
	return func(filter *Filter) {
		filter.where = w
	}
}

func WithLimiter(offset, limit int) Option {
	return func(filter *Filter) {
		filter.limiter = &Limiter{
			limit:  limit,
			offset: offset,
		}
	}
}

func WithOrder(col string, ordering Ordering) Option {
	return func(filter *Filter) {
		filter.order = &Order{
			col:      col,
			ordering: ordering,
		}
	}
}
