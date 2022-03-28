package cache

type Cache interface{
	SetShortUrl(string, string) error
	GetShortUrl(string) (string, error)
}
