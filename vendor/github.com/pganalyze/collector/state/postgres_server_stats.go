package state

// PostgresServerStats - Statistics for a Postgres server.
type PostgresServerStats struct {
	CurrentXactId   Xid8
	NextMultiXactId Xid8

	XminHorizonBackend                Xid
	XminHorizonReplicationSlot        Xid
	XminHorizonReplicationSlotCatalog Xid
	XminHorizonPreparedXact           Xid
	XminHorizonStandby                Xid
}

// FullXminHorizonBackend - Returns XminHorizonBackend in 64-bit FullTransactionId
func (ss PostgresServerStats) FullXminHorizonBackend() int64 {
	return int64(XidToXid8(ss.XminHorizonBackend, Xid8(ss.CurrentXactId)))
}

// FullXminHorizonReplicationSlot - Returns XminHorizonReplicationSlot in 64-bit FullTransactionId
func (ss PostgresServerStats) FullXminHorizonReplicationSlot() int64 {
	return int64(XidToXid8(ss.XminHorizonReplicationSlot, Xid8(ss.CurrentXactId)))
}

// FullXminHorizonReplicationSlotCatalog - Returns XminHorizonReplicationSlotCatalog in 64-bit FullTransactionId
func (ss PostgresServerStats) FullXminHorizonReplicationSlotCatalog() int64 {
	return int64(XidToXid8(ss.XminHorizonReplicationSlotCatalog, Xid8(ss.CurrentXactId)))
}

// FullXminHorizonPreparedXact - Returns XminHorizonPreparedXact in 64-bit FullTransactionId
func (ss PostgresServerStats) FullXminHorizonPreparedXact() int64 {
	return int64(XidToXid8(ss.XminHorizonPreparedXact, Xid8(ss.CurrentXactId)))
}

// FullXminHorizonStandby - Returns XminHorizonStandby in 64-bit FullTransactionId
func (ss PostgresServerStats) FullXminHorizonStandby() int64 {
	return int64(XidToXid8(ss.XminHorizonStandby, Xid8(ss.CurrentXactId)))
}
