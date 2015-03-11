package sitecache

type CacheItem struct {
  Data []byte
  CheckSum uint64
}

type SiteCache struct {
  CachedItems map[string]CacheItem
}

func (cache *SiteCache) loadCache(siteRoot string) error {
  return nil
}

func NewSiteCache(siteRoot string) (*SiteCache, error) {
  cache := new(SiteCache)
  err := cache.loadCache(siteRoot)
  if err != nil {
    cache = nil
  }

  return cache, err
}
