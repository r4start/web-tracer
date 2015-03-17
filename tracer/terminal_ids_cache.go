package tracer

import (
  "sync"
  "sort"

  "github.com/cznic/sortutil"
)

type TerminalIdsCache struct {
  ids []uint64
  guard sync.RWMutex
}

func NewTerminalIdsCache() TerminalIdsCache {
  newCache := TerminalIdsCache{}
  newCache.ids = make([]uint64, 0)

  return newCache
}

func (cache *TerminalIdsCache) AppendIds(ids []uint64) {
  cache.guard.Lock()

  inserted := false
  outOfRange := len(cache.ids)
  for _, v := range ids {
    pos := sortutil.SearchUint64s(cache.ids, v)
    if pos != outOfRange {
      continue
    }

    cache.ids = append(cache.ids, v)
    inserted = true
  }

  if inserted {
    sort.Sort(Uint64Slice(cache.ids))
  }

  cache.guard.Unlock()
}

func (cache *TerminalIdsCache) AppendId(id uint64) {
  cache.guard.Lock()

  outOfRange := len(cache.ids)
  pos := sortutil.SearchUint64s(cache.ids, id)
  if pos != outOfRange {
    cache.guard.Unlock()
    return    
  }

  cache.ids = append(cache.ids, id)

  sort.Sort(Uint64Slice(cache.ids))

  cache.guard.Unlock()
}

func (cache *TerminalIdsCache) GetIds() []uint64 {
  cache.guard.RLock()
  ids := cache.ids
  cache.guard.RUnlock()

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
