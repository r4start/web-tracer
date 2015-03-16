package sitecache

import (
  _ "io"
  "os"
  "log"
  "errors"
  "strings"
  "io/ioutil"
  "hash/fnv"
  "strconv"
)

type CacheItem struct {
  Data []byte
  CheckSum uint64
}

type SiteCache struct {
  cachedItems map[string]CacheItem

  SiteRoot string
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
    cache.cachedItems[rootName + file.Name()] = item
    log.Println("Added file ", rootName + file.Name(), "to cache")
  }

  return nil
}

func NewSiteCache(siteRoot string) (*SiteCache, error) {
  cache := new(SiteCache)
  cache.cachedItems = make(map[string]CacheItem)

  cwd, err := os.Getwd()
  if err != nil {
    return nil, err
  }

  err = os.Chdir(siteRoot)
  if err != nil {
    return nil, err
  }

  cache.SiteRoot = siteRoot
  err = cache.loadCache("./")
  if err != nil {
    cache = nil
  }

  err = os.Chdir(cwd)
  return cache, err
}

func (cache *SiteCache) GetItem(key string) (CacheItem, error) {
  val, exists := cache.cachedItems[key]
  if !exists {
    return CacheItem{}, errors.New("No item with specified key.")
  } else {
    return val, nil
  }
}

func (cache *SiteCache) GetView() map[string]CacheItem {
  return cache.cachedItems
}

func (cache *SiteCache) MarshalJSON() ([]byte, error) {
  items := cache.GetView()

  array := `{ "cache" : [`
  for k, v := range items {
    array := `{"name" : "` + k + `",`
    array += `"hash" : "` + strconv.FormatUint(v.CheckSum, 16) + `",`
    array += `"size" : "` + strconv.FormatInt(int64(len(v.Data)), 10) + `"},`
  }

  if len(items) != 0 {
    array = array[: len(array) - 1]  
  }
  
  array += `]}`
  return []byte(array), nil
}

func (cache *SiteCache) UnmarshalJSON([]byte) error {
  return nil
}
