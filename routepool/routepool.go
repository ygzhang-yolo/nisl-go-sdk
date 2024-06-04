package routepool

import "github.com/panjf2000/ants/v2"

var Rpool RoutePool

type RoutePool struct {
	pool *ants.Pool
}

func (r *RoutePool) Init(size int) {
	pool, _ := ants.NewPool(size)
	r.pool = pool
}

func (r *RoutePool) Close() {
	r.pool.Release()
}

func (r *RoutePool) Submit(task func()) error {
	return r.pool.Submit(task)
}
