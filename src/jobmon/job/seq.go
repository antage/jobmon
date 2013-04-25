package job

import (
	"sync"
)

type seqInt64 struct {
	sync.Mutex

	currentId int64
}

func (g *seqInt64) next() int64 {
	g.Lock()
	defer g.Unlock()

	g.currentId++
	return g.currentId
}
