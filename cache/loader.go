package cache

// Fetcher loads the value based on key
type Fetcher func(key any) (any, error)

type Loader interface {
	Get(k any) (any, bool)
}
