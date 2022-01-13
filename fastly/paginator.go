package fastly

// TODO: In go 1.18 (Feb 2022) use generics to reduce the duplicated code.

// PaginatorACLEntries represents a paginator.
type PaginatorACLEntries interface {
	HasNext() bool
	Remaining() int
	GetNext() ([]*ACLEntry, error)
}

// PaginatorDictionaryItems represents a paginator.
type PaginatorDictionaryItems interface {
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
