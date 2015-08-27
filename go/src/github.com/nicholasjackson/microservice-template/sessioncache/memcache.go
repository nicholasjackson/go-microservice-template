package sessioncache

type MemCache interface {
	Set(key string, value interface{}, cacheDuration int)
	Get(key string) (interface{}, error)
	New()
}
