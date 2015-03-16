package sitecache

import (
  _ "io"
  _ "os"
  "log"
  "strings"
  "io/ioutil"
  "hash/fnv"
)

type CacheItem struct {
  Data []byte
  CheckSum uint64
}

type SiteCache struct {
  CachedItems map[string]CacheItem
}

func NewCacheItem(data []byte) CacheItem {
  hasher := fnv.New64()
  hasher.Write(data)

  return CacheItem{data, hasher.Sum64()}
}

func (cache *SiteCache) loadCache(siteRoot string) error {
  log.Println("Loading directory ", siteRoot)

  info, err := ioutil.ReadDir(siteRoot)
  if err != nil {
    return err
  }

  var rootName string

    if strings.HasSuffix(siteRoot, "/") {
      rootName = siteRoot
    } else {
      rootName = siteRoot + "/"
    }

  for _, file := range info {
    var data []byte
    
    if file.Mode().IsDir() {
      err = cache.loadCache(rootName + file.Name())
    } else {
      data, err = ioutil.ReadFile(rootName + file.Name())
    }

    if err != nil {
      return err
    }

    if len(data) == 0 {
      continue
    }

    item := NewCacheItem(data)
    cache.CachedItems[rootName + file.Name()] = item
    log.Println("Added file ", rootName + file.Name(), "to cache")
  }

  return nil
}

func NewSiteCache(siteRoot string) (*SiteCache, error) {
  cache := new(SiteCache)
  cache.CachedItems = make(map[string]CacheItem)

  err := cache.loadCache(siteRoot)
  if err != nil {
    cache = nil
  }

  return cache, err
}
