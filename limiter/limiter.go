package limiter

type limiter interface {
	Allow(target string) bool
}
