package minisql

import "sync"

var queryPool sync.Pool

// Pooled attempts to retrieve a query from the global
// query pool.
func Pooled() *Query {
	if v := queryPool.Get(); v != nil {
		q := v.(*Query)
		q.Reset()
		return q
	}
	return new(Query)
}

// Release puts a query into the global pool. Please make sure
// not to use the Query after it has been released.
func Release(q *Query) {
	queryPool.Put(q)
}
