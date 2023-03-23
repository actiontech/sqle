package cluster

var IsClusterMode bool = false

var DefaultNode Node = &NoClusterNode{}

type Node interface {
	Join(serverId string)
	Leave()
	IsLeader() bool
}

type NoClusterNode struct{}

func (c *NoClusterNode) Join(serverId string) {
}

func (c *NoClusterNode) Leave() {
}

func (c *NoClusterNode) IsLeader() bool {
	return true
}
