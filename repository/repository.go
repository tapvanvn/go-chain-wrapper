package repository

import (
	"fmt"
	"strconv"

	engines "github.com/tapvanvn/godbengine"
)

func GetLastBlock(chain string) uint64 {

	eng := engines.GetEngine()
	pool := eng.GetMemPool()
	key := fmt.Sprintf("mine.lastest_%s", chain)
	if pool != nil {
		if value, err := pool.Get(key); err == nil {
			if lastestBlock, err := strconv.ParseUint(value, 10, 64); err != nil {
				return lastestBlock
			}
		}
	}
	return 0
}

func PutLastBlock(chain string, value uint64) {
	eng := engines.GetEngine()
	pool := eng.GetMemPool()
	key := fmt.Sprintf("mine.lastest_%s", chain)
	if pool != nil {
		pool.Set(key, fmt.Sprintf("%d", value))
	}
}
