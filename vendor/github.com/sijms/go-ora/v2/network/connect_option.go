package network

import (
	"strconv"
	"strings"

	"github.com/sijms/go-ora/v2/trace"
)

type ClientData struct {
	ProgramPath string
	ProgramName string
	UserName    string
	HostName    string
	DriverName  string
	PID         int
}
type ConnectionOption struct {
	//Port                  int
	TransportConnectTo    int
	SSLVersion            string
	WalletDict            string
	TransportDataUnitSize uint32
	SessionDataUnitSize   uint32
	Protocol              string
	//Host                  string
	UserID      string
	Servers     []string
	Ports       []int
	serverIndex int
	//IP string
	SID string
	//Addr string
	//Server string
	ServiceName  string
	InstanceName string
	DomainName   string
	DBName       string
	ClientData   ClientData
	//InAddrAny bool
	Tracer       trace.Tracer
	connData     string
	SNOConfig    map[string]string
	PrefetchRows int
	SSL          bool
	SSLVerify    bool
}

func (op *ConnectionOption) AddServer(host string, port int) {
	for i := 0; i < len(op.Servers); i++ {
		if strings.ToUpper(host) == strings.ToUpper(op.Servers[i]) &&
			port == op.Ports[i] {
			return
		}
	}
	op.Servers = append(op.Servers, host)
	op.Ports = append(op.Ports, port)
}

func (op *ConnectionOption) GetActiveServer(jump bool) (string, int) {
	if jump {
		op.serverIndex++
	}
	if op.serverIndex >= len(op.Servers) {
		return "", 0
	}
	return op.Servers[op.serverIndex], op.Ports[op.serverIndex]
}
func (op *ConnectionOption) ConnectionData() string {
	//if len(op.connData) > 0 {
	//	return op.connData
	//}
	host, port := op.GetActiveServer(false)
	FulCid := "(CID=(PROGRAM=" + op.ClientData.ProgramPath + ")(HOST=" + op.ClientData.HostName + ")(USER=" + op.ClientData.UserName + "))"
	address := "(ADDRESS=(PROTOCOL=" + op.Protocol + ")(HOST=" + host + ")(PORT=" + strconv.Itoa(port) + "))"
	result := "(CONNECT_DATA="
	if op.SID != "" {
		result += "(SID=" + op.SID + ")"
	} else {
		result += "(SERVICE_NAME=" + op.ServiceName + ")"
	}
	if op.InstanceName != "" {
		result += "(INSTANCE_NAME=" + op.InstanceName + ")"
	}
	result += FulCid
	return "(DESCRIPTION=" + address + result + "))"
}
