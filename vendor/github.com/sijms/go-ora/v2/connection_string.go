package go_ora

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type PromotableTransaction int

const (
	Promotable PromotableTransaction = 1
	Local      PromotableTransaction = 0
)

type DBAPrivilege int

const (
	NONE    DBAPrivilege = 0
	SYSDBA  DBAPrivilege = 0x20
	SYSOPER DBAPrivilege = 0x40
)
const defaultPort int = 1521

func DBAPrivilegeFromString(s string) DBAPrivilege {
	S := strings.ToUpper(s)
	if S == "SYSDBA" {
		return SYSDBA
	} else if S == "SYSOPER" {
		return SYSOPER
	} else {
		return NONE
	}
}

type EnList int

const (
	FALSE   EnList = 0
	TRUE    EnList = 1
	DYNAMIC EnList = 2
)

func EnListFromString(s string) EnList {
	S := strings.ToUpper(s)
	if S == "TRUE" {
		return TRUE
	} else if S == "DYNAMIC" {
		return DYNAMIC
	} else {
		return FALSE
	}
}

type ConnectionString struct {
	DataSource            string
	Host                  string
	Port                  int
	Servers               []string
	Ports                 []int
	SID                   string
	ServiceName           string
	InstanceName          string
	DBAPrivilege          DBAPrivilege
	EnList                EnList
	ConnectionLifeTime    int
	IncrPoolSize          int
	DecrPoolSize          int
	MaxPoolSize           int
	MinPoolSize           int
	Password              string
	PasswordSecurityInfo  bool
	Pooling               bool
	ConnectionTimeOut     int
	UserID                string
	PromotableTransaction PromotableTransaction
	ProxyUserID           string
	ProxyPassword         string
	ValidateConnection    bool
	StmtCacheSize         int
	StmtCachePurge        bool
	HaEvent               bool
	LoadBalance           bool
	MetadataBooling       bool
	ContextConnection     bool
	SelfTuning            bool
	SSL                   bool
	SSLVerify             bool
	ApplicationEdition    string
	PoolRegulator         int
	ConnectionPoolTimeout int
	Trace                 string // Trace file
	PrefetchRows          int
	WalletPath            string
	w                     *wallet
}

// BuildUrl create databaseURL from server, port, service, user, password, urlOptions
// this function help build a will formed databaseURL and accept any character as it
// convert special charters to corresponding values in URL
func BuildUrl(server string, port int, service, user, password string, options map[string]string) string {
	ret := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", url.QueryEscape(user), url.QueryEscape(password),
		url.QueryEscape(server), port, url.QueryEscape(service))
	if options != nil {
		ret += "?"
		for key, val := range options {
			val = strings.TrimSpace(val)
			for _, temp := range strings.Split(val, ",") {
				temp = strings.TrimSpace(temp)
				if strings.ToUpper(key) == "SERVER" {
					ret += fmt.Sprintf("%s=%s&", key, temp)
				} else {
					ret += fmt.Sprintf("%s=%s&", key, url.QueryEscape(temp))
				}
			}
		}
		ret = strings.TrimRight(ret, "&")
	}
	return ret
}

// newConnectionStringFromUrl create new connection string from databaseURL data and options
func newConnectionStringFromUrl(databaseUrl string) (*ConnectionString, error) {
	u, err := url.Parse(databaseUrl)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	p := u.Port()
	ret := &ConnectionString{
		Port:                  defaultPort,
		DBAPrivilege:          NONE,
		EnList:                TRUE,
		IncrPoolSize:          5,
		DecrPoolSize:          5,
		MaxPoolSize:           100,
		MinPoolSize:           1,
		ConnectionTimeOut:     15,
		PromotableTransaction: Promotable,
		StmtCacheSize:         20,
		MetadataBooling:       true,
		SelfTuning:            true,
		PoolRegulator:         100,
		ConnectionPoolTimeout: 15,
		PrefetchRows:          25,
		SSL:                   false,
		SSLVerify:             true,
		Servers:               make([]string, 0, 3),
		Ports:                 make([]int, 0, 3),
	}
	ret.UserID = u.User.Username()
	ret.Password, _ = u.User.Password()
	if p != "" {
		port, err := strconv.Atoi(p)
		if err != nil {
			ret.Ports = append(ret.Ports, defaultPort)
		} else {
			ret.Ports = append(ret.Ports, port)
		}

	} else {
		ret.Ports = append(ret.Ports, defaultPort)
	}
	if len(u.Host) > 0 {
		idx := strings.Index(u.Host, ":")
		if idx > 0 {
			ret.Servers = append(ret.Servers, u.Host[:idx])
		} else {
			ret.Servers = append(ret.Servers, u.Host)
		}
	}
	ret.ServiceName = strings.Trim(u.Path, "/")
	if q != nil {
		for key, val := range q {
			switch strings.ToUpper(key) {
			//case "DATA SOURCE":
			//	conStr.DataSource = val
			case "SERVER":
				for _, srv := range val {
					srv = strings.TrimSpace(srv)
					idx := strings.Index(srv, ":")
					if idx > 0 {
						ret.Servers = append(ret.Servers, srv[:idx])
						port, err := strconv.Atoi(srv[idx+1:])
						if err != nil {
							port = 0
						}
						if port == 0 {
							ret.Ports = append(ret.Ports, defaultPort)
						} else {
							ret.Ports = append(ret.Ports, port)
						}
					} else {
						ret.Servers = append(ret.Servers, srv)
						ret.Ports = append(ret.Ports, defaultPort)
					}
				}
			case "SERVICE NAME":
				ret.ServiceName = val[0]
			case "SID":
				ret.SID = val[0]
			case "INSTANCE NAME":
				ret.InstanceName = val[0]
			case "WALLET":
				ret.WalletPath = val[0]
			case "SSL":
				ret.SSL = strings.ToUpper(val[0]) == "TRUE" || strings.ToUpper(val[0]) == "ENABLE"
			case "SSL VERIFY":
				ret.SSLVerify = strings.ToUpper(val[0]) == "TRUE" || strings.ToUpper(val[0]) == "ENABLE"
			case "DBA PRIVILEGE":
				ret.DBAPrivilege = DBAPrivilegeFromString(val[0])
			case "ENLIST":
				ret.EnList = EnListFromString(val[0])
			case "CONNECT TIMEOUT":
				fallthrough
			case "CONNECTION TIMEOUT":
				ret.ConnectionTimeOut, err = strconv.Atoi(val[0])
				if err != nil {
					return nil, errors.New("CONNECTION TIMEOUT value must be an integer")
				}
			case "INC POOL SIZE":
				ret.IncrPoolSize, err = strconv.Atoi(val[0])
				if err != nil {
					return nil, errors.New("INC POOL SIZE value must be an integer")
				}
			case "DECR POOL SIZE":
				ret.DecrPoolSize, err = strconv.Atoi(val[0])
				if err != nil {
					return nil, errors.New("DECR POOL SIZE value must be an integer")
				}
			case "MAX POOL SIZE":
				ret.MaxPoolSize, err = strconv.Atoi(val[0])
				if err != nil {
					return nil, errors.New("MAX POOL SIZE value must be an integer")
				}
			case "MIN POOL SIZE":
				ret.MinPoolSize, err = strconv.Atoi(val[0])
				if err != nil {
					return nil, errors.New("MIN POOL SIZE value must be an integer")
				}
			case "POOL REGULATOR":
				ret.PoolRegulator, err = strconv.Atoi(val[0])
				if err != nil {
					return nil, errors.New("POOL REGULATOR value must be an integer")
				}
			case "STATEMENT CACHE SIZE":
				ret.StmtCacheSize, err = strconv.Atoi(val[0])
				if err != nil {
					return nil, errors.New("STATEMENT CACHE SIZE value must be an integer")
				}
			case "CONNECTION POOL TIMEOUT":
				ret.ConnectionPoolTimeout, err = strconv.Atoi(val[0])
				if err != nil {
					return nil, errors.New("CONNECTION POOL TIMEOUT value must be an integer")
				}
			case "CONNECTION LIFETIME":
				ret.ConnectionLifeTime, err = strconv.Atoi(val[0])
				if err != nil {
					return nil, errors.New("CONNECTION LIFETIME value must be an integer")
				}
			case "PERSIST SECURITY INFO":
				ret.PasswordSecurityInfo = val[0] == "TRUE"
			case "POOLING":
				ret.Pooling = val[0] == "TRUE"
			case "VALIDATE CONNECTION":
				ret.ValidateConnection = val[0] == "TRUE"
			case "STATEMENT CACHE PURGE":
				ret.StmtCachePurge = val[0] == "TRUE"
			case "HA EVENTS":
				ret.HaEvent = val[0] == "TRUE"
			case "LOAD BALANCING":
				ret.LoadBalance = val[0] == "TRUE"
			case "METADATA POOLING":
				ret.MetadataBooling = val[0] == "TRUE"
			case "SELF TUNING":
				ret.SelfTuning = val[0] == "TRUE"
			case "CONTEXT CONNECTION":
				ret.ContextConnection = val[0] == "TRUE"
			case "PROMOTABLE TRANSACTION":
				if val[0] == "PROMOTABLE" {
					ret.PromotableTransaction = Promotable
				} else {
					ret.PromotableTransaction = Local
				}
			case "APPLICATION EDITION":
				ret.ApplicationEdition = val[0]
			//case "USER ID":
			//	val = strings.Trim(val, "'")
			//	conStr.UserID = strings.Trim(val, "\"")
			//	if conStr.UserID == "\\" {
			//		// get os user and password
			//	}
			case "PROXY USER ID":
				ret.ProxyUserID = val[0]
			//case "PASSWORD":
			//	val = strings.Trim(val, "'")
			//	conStr.Password = strings.Trim(val, "\"")
			case "PROXY PASSWORD":
				ret.ProxyPassword = val[0]
			case "TRACE FILE":
				ret.Trace = val[0]
			case "PREFETCH_ROWS":
				ret.PrefetchRows, err = strconv.Atoi(val[0])
				if err != nil {
					ret.PrefetchRows = 25
				}
			}
		}
	}
	if len(ret.Servers) == 0 {
		return nil, errors.New("empty connection servers")
	}
	if len(ret.WalletPath) > 0 {
		if len(ret.ServiceName) == 0 {
			return nil, errors.New("you should specify server/service if you will use wallet")
		}
		ret.w, err = NewWallet(path.Join(ret.WalletPath, "cwallet.sso"))
		if err != nil {
			return nil, err
		}
		if len(ret.Password) == 0 {
			serv := ret.Servers[0]
			port := ret.Ports[0]
			cred, err := ret.w.getCredential(serv, port, ret.ServiceName, ret.UserID)
			if err != nil {
				return nil, err
			}
			if cred == nil {
				return nil, errors.New(
					fmt.Sprintf("cannot find credentials for server: %s, service: %s,  username: %s",
						ret.Host, ret.ServiceName, ret.UserID))
			}
			ret.UserID = cred.username
			ret.Password = cred.password
		}
	}
	return ret, ret.validate()
}

// validate check is data in connection string is correct and fulfilled
func (connStr *ConnectionString) validate() error {
	if !connStr.Pooling {
		connStr.MaxPoolSize = -1
		connStr.MinPoolSize = 0
		connStr.IncrPoolSize = -1
		connStr.DecrPoolSize = 0
		connStr.PoolRegulator = 0
	}
	//if connStr.SSL && (connStr.w == nil || len(connStr.w.certificates) == 0) {
	//	return errors.New("tcps need a valid wallet contains server and client certificates")
	//}
	if len(connStr.UserID) == 0 {
		return errors.New("empty user name")
	}
	if len(connStr.Password) == 0 {
		return errors.New("empty password")
	}
	if len(connStr.SID) == 0 && len(connStr.ServiceName) == 0 {
		return errors.New("empty SID and service name")
	}
	return nil
}
