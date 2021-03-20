package postgres

import "sync"

// simple slow way to solve race condition problem for sql tables
var mux = sync.RWMutex{}
