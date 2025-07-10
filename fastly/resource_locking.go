package fastly

import (
	"context"
	"sync"
	"weak"
)

// A ResourceLockManager stores [sync.Mutex] objects used to manage
// serialization of requests.
//
// The objects are identified using string keys, and there will be at
// most one [sync.Mutex] per distinct key.
//
// This structure is safe to use from concurrent goroutines, as the
// functions which manipulate it will use its internal [sync.Mutex] to
// protect the data structure.
type ResourceLockManager struct {
	m     sync.Mutex
	locks map[string]weak.Pointer[sync.Mutex]
}

// NewResourceLockManager creates constructs a [ResourceLockManager]
// with an empty map.
func NewResourceLockManager() *ResourceLockManager {
	return &ResourceLockManager{
		locks: make(map[string]weak.Pointer[sync.Mutex]),
	}
}

// Get returns a [sync.Mutex] which should be used for serializing
// requests associated with the resource named in the key parameter.
//
// If no [sync.Mutex] exists for the specified key, one will be
// created, stored, and returned. Subsequent requests for the same key
// will return the same [sync.Mutex], unless all previous requests have
// completed their usage of the [sync.Mutex] and the garbage collector
// has reaped it. In that case, a new [sync.Mutex] will be created,
// stored, and returned.
//
// This function is safe to use from concurrent goroutines.
func (rlm *ResourceLockManager) Get(key string) *sync.Mutex {
	rlm.m.Lock()
	defer rlm.m.Unlock()

	ptr, exists := rlm.locks[key]
	if !exists || ptr.Value() == nil {
		l := sync.Mutex{}
		rlm.locks[key] = weak.Make(&l)
		return &l
	}

	return ptr.Value()
}

type key struct{}

var contextKey key

// NewContextForResourceID returns a [context.Context] which contains
// the resource ID specified in the id parameter.
//
// This [context.Context] should be passed to the various mutating
// request methods of [fastly.Client] so that those request methods
// can serialize concurrent requests against the same resource.
func NewContextForResourceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey, id)
}

func resourceIDFromContext(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(contextKey).(string)
	return s, ok
}
