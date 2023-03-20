package cluster

var DefaultNode Node = &NoClusterNode{}

type Node interface {
	Join(serverId string)
	Leave()
	IsLeader() bool
}

type NoClusterNode struct{}

func (c *NoClusterNode) Join(serverId string) {
	return
}

func (c *NoClusterNode) Leave() {
	return
}

func (c *NoClusterNode) IsLeader() bool {
	return true
}
