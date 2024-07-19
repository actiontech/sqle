package state

import "time"

type OidToIdxMap map[Oid](map[Oid]int32)

func MakeOidToIdxMap() OidToIdxMap {
	return make(map[Oid](map[Oid]int32))
}

func (m OidToIdxMap) Put(dbOid, objOid Oid, idx int32) {
	if _, ok := m[dbOid]; !ok {
		m[dbOid] = make(map[Oid]int32)
	}
	m[dbOid][objOid] = idx
}

func (m OidToIdxMap) Get(dbOid, objOid Oid) int32 {
	if _, ok := m[dbOid]; !ok {
		return -1
	}
	idx, ok := m[dbOid][objOid]
	if !ok {
		return -1
	}
	return idx
}

// XidToXid8 - Converts Xid (32-bit transaction ID) to Xid8 (64-bit FullTransactionId)
// by calculating and adding an epoch from the current transaction ID
func XidToXid8(xid Xid, currentXactId Xid8) Xid8 {
	// Do not proceed the conversion if either of inputs is 0
	// The currentXactID can be 0 on replicas
	if xid == 0 || currentXactId == 0 {
		return 0
	}
	// If we simply shift the currentXactId, it'll give the epoch of the current transaction ID, which may be different
	// from the epoch of the given xid (the one we want to add).
	// By subtracting the xid from the current one, we can get the epoch of the given xid.
	xidEpoch := int32((currentXactId - Xid8(xid)) >> 32)
	return Xid8(xidEpoch)<<32 | Xid8(xid)
}

func getTimeZoneFromSettings(settings []PostgresSetting) *time.Location {
	for _, setting := range settings {
		if setting.Name != "log_timezone" {
			continue
		}
		if !setting.ResetValue.Valid {
			return nil
		}

		zoneStr := setting.ResetValue.String
		zone, err := time.LoadLocation(zoneStr)
		if err != nil {
			return nil
		}
		return zone
	}
	return nil
}
