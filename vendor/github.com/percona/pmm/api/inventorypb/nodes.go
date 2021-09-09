package inventorypb

//go-sumtype:decl Node

// Node is a common interface for all types of Nodes.
type Node interface {
	sealedNode() //nolint:unused
}

// in order of NodeType enum

func (*GenericNode) sealedNode()             {}
func (*ContainerNode) sealedNode()           {}
func (*RemoteNode) sealedNode()              {}
func (*RemoteRDSNode) sealedNode()           {}
func (*RemoteAzureDatabaseNode) sealedNode() {}
