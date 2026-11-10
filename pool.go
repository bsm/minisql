package minisql

import "sync"

var queryPool sync.Pool

// Pooled attempts to retrieve a query from the global
// query pool.
func Pooled() *Query {
	if v := queryPool.Get(); v != nil {
		return v.(*Query)
	}
	return new(Query)
}

// Release resets a query and puts it into the global pool. Please make
// sure not to use the Query after it has been released.
func Release(q *Query) {
	// Reset here, not in Pooled, so an idle query holds no old arguments.
	q.Reset()
	queryPool.Put(q)
}
