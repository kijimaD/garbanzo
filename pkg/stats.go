package garbanzo

type Stats struct {
	ReadCount  int
	EventCount int
	CacheCount int
}

func newStats() *Stats {
	return &Stats{
		ReadCount:  0,
		EventCount: 0,
		CacheCount: 0,
	}
}
