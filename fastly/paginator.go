package fastly

// TODO: In go 1.18 (Feb 2022) use generics to reduce the duplicated code.

// PaginatorACLEntries represents a paginator.
type PaginatorACLEntries interface {
	HasNext() bool
	Remaining() int
	GetNext() ([]*ACLEntry, error)
}

// PaginatorDictItems represents a paginator.
type PaginatorDictItems interface {
	HasNext() bool
	Remaining() int
	GetNext() ([]*DictionaryItem, error)
}

// PaginatorServices represents a paginator.
type PaginatorServices interface {
	HasNext() bool
	Remaining() int
	GetNext() ([]*Service, error)
}
