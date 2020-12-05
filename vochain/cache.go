package vochain

import (
	"time"

	"gitlab.com/vocdoni/go-dvote/types"
)

// CacheAdd adds a new vote proof to the local cache
func (v *State) CacheAdd(id [32]byte, vc *types.CacheTx) {
	if len(id) == 0 {
		return
	}
	v.voteCacheLock.Lock()
	defer v.voteCacheLock.Unlock()
	v.voteCache[id] = vc
}

// CacheDel deletes an existing vote proof from the local cache
func (v *State) CacheDel(id [32]byte) {
	v.voteCacheLock.Lock()
	defer v.voteCacheLock.Unlock()
	delete(v.voteCache, id)
	if v.MemPoolRemoveTxKey != nil {
		v.MemPoolRemoveTxKey(id, false)
	}
}

// CacheGet fetch an existing vote proof from the local cache
func (v *State) CacheGet(id [32]byte) *types.CacheTx {
	if len(id) == 0 {
		return nil
	}
	v.voteCacheLock.RLock()
	defer v.voteCacheLock.RUnlock()
	return v.voteCache[id]
}

// CachePurge removes the old cache saved votes
func (v *State) CachePurge(height int64) {
	if height%6 != 0 {
		return
	}
	v.voteCacheLock.Lock()
	defer v.voteCacheLock.Unlock()
	for id, vp := range v.voteCache {
		if time.Since(vp.Created) > voteCachePurgeThreshold {
			delete(v.voteCache, id)
		}
	}
}

// CacheSize returns the current size of the vote cache
func (v *State) CacheSize() int {
	v.voteCacheLock.RLock()
	defer v.voteCacheLock.RUnlock()
	return len(v.voteCache)
}
