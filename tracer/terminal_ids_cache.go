package tracer

import (
  "log"
  "sync"
  "sort"

  "github.com/cznic/sortutil"
  "github.com/jinzhu/gorm"
  _ "github.com/mattn/go-sqlite3"
)

type TerminalIdsCache struct {
  ids map[uint64] interface{}
  guard sync.RWMutex
}

func NewTerminalIdsCache() TerminalIdsCache {
  newCache := TerminalIdsCache{}
  newCache.ids = make(map[uint64] interface{}, 0)
  newCache.guard = sync.RWMutex{}

  return newCache
}

func LoadIdsFromDb(dbName string) []uint64 {
  conn, err := gorm.Open("sqlite3", dbName)
  if err != nil {
    return make([]uint64, 0)
  }

  ids := make([]uint64, 0)
  rows, e := conn.Raw("select distinct terminal_id from log_entries;").
                  Rows()
  
  if e != nil {
    log.Println(e)
    return make([]uint64, 0)
  }

  defer rows.Close()
  for rows.Next() {
    var termId uint64
    
    err = rows.Scan(&termId)
    if err == nil {
      ids = append(ids, termId)
    }
  }

  return ids
}

func (cache *TerminalIdsCache) AppendIds(ids []uint64) {
  sort.Sort(Uint64Slice(ids))
  ids = ids[:sortutil.Dedupe(Uint64Slice(ids))]

  cache.guard.Lock()
  defer cache.guard.Unlock()

  for _, v := range ids {
    cache.ids[v] = nil
  }
}

func (cache *TerminalIdsCache) AppendId(id uint64) {
  cache.guard.Lock()
  defer cache.guard.Unlock()

  cache.ids[id] = nil
}

func (cache *TerminalIdsCache) GetIds() []uint64 {
  cache.guard.RLock()
  defer cache.guard.RUnlock()
  ids := make([]uint64, 0, len(cache.ids))

  for k, _ := range cache.ids {
    ids = append(ids, k)
  }

  return ids
}

type Uint64Slice []uint64

func (s Uint64Slice) Len() int {
  return len(s)
}

func (s Uint64Slice) Less(i, j int) bool {
  return s[i] < s[j]
}

func (s Uint64Slice) Swap(i, j int) {
  s[i] = s[i] ^ s[j]
  s[j] = s[i] ^ s[j]
  s[i] = s[j] ^ s[i]
}
