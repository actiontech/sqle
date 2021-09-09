package agentpb

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	mdAgentID          = "pmm-agent-id"
	mdAgentVersion     = "pmm-agent-version"
	mdAgentMetricsPort = "pmm-agent-metrics-port"
	mdAgentNodeID      = "pmm-agent-node-id"
	mdServerVersion    = "pmm-server-version"
)

// AgentConnectMetadata represents metadata sent by pmm-agent with Connect RPC method call.
type AgentConnectMetadata struct {
	ID          string
	Version     string
	MetricsPort uint16
}

// ServerConnectMetadata represents metadata sent by pmm-managed in response to Connect RPC method call.
type ServerConnectMetadata struct {
	AgentRunsOnNodeID string
	ServerVersion     string
}

func getValue(md metadata.MD, key string) string {
	vs := md.Get(key)
	if len(vs) == 1 {
		return vs[0]
	}
	return ""
}

// AddAgentConnectMetadata adds pmm-agent's metadata to outgoing context. Used by pmm-agent.
func AddAgentConnectMetadata(ctx context.Context, md *AgentConnectMetadata) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		mdAgentID, md.ID,
		mdAgentVersion, md.Version,
		mdAgentMetricsPort, strconv.FormatUint(uint64(md.MetricsPort), 10),
	)
}

// ReceiveAgentConnectMetadata receives pmm-agent's metadata. Used by pmm-managed.
func ReceiveAgentConnectMetadata(stream grpc.ServerStream) (*AgentConnectMetadata, error) {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "ReceiveAgentConnectMetadata: failed to get metadata")
	}
	if md == nil || md.Len() == 0 {
		return nil, status.Errorf(codes.DataLoss, "ReceiveAgentConnectMetadata: empty metadata")
	}

	// metrics port is optional
	var mp uint64
	if mpS := getValue(md, mdAgentMetricsPort); mpS != "" {
		var err error
		if mp, err = strconv.ParseUint(mpS, 10, 16); err != nil {
			return nil, status.Errorf(codes.DataLoss, "ReceiveAgentConnectMetadata: %s: %s", mdAgentMetricsPort, err)
		}
	}

	return &AgentConnectMetadata{
		ID:          getValue(md, mdAgentID),
		Version:     getValue(md, mdAgentVersion),
		MetricsPort: uint16(mp),
	}, nil
}

// SendServerConnectMetadata sends pmm-managed's metadata. Used by pmm-managed.
func SendServerConnectMetadata(stream grpc.ServerStream, md *ServerConnectMetadata) error {
	header := metadata.Pairs(
		mdAgentNodeID, md.AgentRunsOnNodeID,
		mdServerVersion, md.ServerVersion,
	)

	// always return gRPC error or nil
	err := stream.SendHeader(header)
	if _, ok := status.FromError(err); err != nil && !ok {
		err = status.Errorf(codes.DataLoss, "SendServerConnectMetadata: SendHeader: %s", err)
	}
	return err
}

// ReceiveServerConnectMetadata receives pmm-managed's metadata. Used by pmm-agent.
func ReceiveServerConnectMetadata(stream grpc.ClientStream) (*ServerConnectMetadata, error) {
	// always return gRPC error or nil
	md, err := stream.Header()
	if _, ok := status.FromError(err); err != nil && !ok {
		err = status.Errorf(codes.DataLoss, "ReceiveServerConnectMetadata: Header: %s", err)
	}
	if err != nil {
		return nil, err
	}

	if md == nil || md.Len() == 0 {
		return nil, status.Errorf(codes.DataLoss, "ReceiveServerConnectMetadata: empty metadata")
	}

	return &ServerConnectMetadata{
		AgentRunsOnNodeID: getValue(md, mdAgentNodeID),
		ServerVersion:     getValue(md, mdServerVersion),
	}, nil
}
