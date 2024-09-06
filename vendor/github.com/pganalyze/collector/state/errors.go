package state

import "errors"

var ErrReplicaCollectionDisabled error = errors.New("monitored server is replica and replication collection disabled via config")
