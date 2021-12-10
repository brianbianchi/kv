package config

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
