package logging

// StatsD interface is to allow the class .. to be injected as a DuckType
type StatsD interface {
	Increment(string)
}
