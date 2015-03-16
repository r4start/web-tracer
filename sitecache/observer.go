package sitecache

import (
  _ "fmt"
  "net/http"
  )

type SiteCacheObserver struct {
  cache *SiteCache
}

func (observer SiteCacheObserver) observeSiteCache() string {
  return "<tr/>"
}

func (observer SiteCacheObserver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  if observer.cache == nil {
    res.WriteHeader(503)
    return
  }

  // cacheRows := observer.observeSiteCache()

  _, err := observer.cache.GetItem("observer.html")
  if err != nil {
    res.WriteHeader(503)
    return
  }
}
