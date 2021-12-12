package config

import (
	"fmt"
	"hash/fnv"

	"github.com/BurntSushi/toml"
)

// Sharding is a method of splitting and storing a single logical dataset in multiple databases.
// By distributing the data among multiple machines,
// a cluster of database systems can store larger dataset and handle additional requests.
// Each address has a unique shard index.
type Shard struct {
	Name    string
	Idx     int
	Address string
}

type Config struct {
	Shards []Shard
}

func ParseFile(filename string) (Config, error) {
	var c Config
	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return Config{}, err
	}
	return c, nil
}

// Better representation of the sharding config.
type Shards struct {
	Count  int
	CurIdx int
	Addrs  map[int]string
}

// Converts and verifies the list of shards.
func ParseShards(shards []Shard, curShardName string) (*Shards, error) {
	shardCount := len(shards)
	shardIdx := -1
	addrs := make(map[int]string)

	for _, s := range shards {
		if _, ok := addrs[s.Idx]; ok {
			return nil, fmt.Errorf("duplicate shard index: %d", s.Idx)
		}

		addrs[s.Idx] = s.Address
		if s.Name == curShardName {
			shardIdx = s.Idx
		}
	}

	for i := 0; i < shardCount; i++ {
		if _, ok := addrs[i]; !ok {
			return nil, fmt.Errorf("shard %d is not found", i)
		}
	}

	if shardIdx < 0 {
		return nil, fmt.Errorf("shard %q was not found", curShardName)
	}

	return &Shards{
		Addrs:  addrs,
		Count:  shardCount,
		CurIdx: shardIdx,
	}, nil
}

// Returns the shard number for the corresponding key.
func (s *Shards) Index(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.Count))
}
