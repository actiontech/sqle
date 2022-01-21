package go_ora

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sijms/go-ora/v2/network"
)

type Lob struct {
	sourceLocator []byte
	destLocator   []byte
	scn           []byte
	sourceOffset  int
	destOffset    int
	charsetID     int
	size          int64
	data          bytes.Buffer
}

// variableWidthChar if lob has variable width char or not
func (lob *Lob) variableWidthChar() bool {
	if len(lob.sourceLocator) > 6 && lob.sourceLocator[6]&128 == 128 {
		return true
	}
	return false
}

// littleEndianClob if CLOB is littleEndian or not
func (lob *Lob) littleEndianClob() bool {
	return len(lob.sourceLocator) > 7 && lob.sourceLocator[7]&64 > 0
}

// getSize return lob size
func (lob *Lob) getSize(connection *Connection) (size int64, err error) {
	session := connection.session
	connection.connOption.Tracer.Print("Read Lob Size")
	err = lob.write(session, 1)
	if err != nil {
		return
	}
	err = lob.read(connection)
	if err != nil {
		return
	}
	size = lob.size
	connection.connOption.Tracer.Print("Lob Size: ", size)
	return
}

// getData return lob data
func (lob *Lob) getData(connection *Connection) (data []byte, err error) {
	connection.connOption.Tracer.Print("Read Lob Data")
	session := connection.session
	lob.sourceOffset = 1
	err = lob.write(session, 2)
	if err != nil {
		return
	}
	err = lob.read(connection)
	if err != nil {
		return
	}
	data = lob.data.Bytes()
	return
}

// write lob command to network session
func (lob *Lob) write(session *network.Session, operationID int) error {
	session.ResetBuffer()
	session.PutBytes(3, 0x60, 0)
	if len(lob.sourceLocator) == 0 {
		session.PutBytes(0)
	} else {
		session.PutBytes(1)
	}
	session.PutUint(len(lob.sourceLocator), 4, true, true)

	if len(lob.destLocator) == 0 {
		session.PutBytes(0)
	} else {
		session.PutBytes(1)
	}
	session.PutUint(len(lob.destLocator), 4, true, true)

	// put offsets
	if session.TTCVersion < 3 {
		session.PutUint(lob.sourceOffset, 4, true, true)
		session.PutUint(lob.destOffset, 4, true, true)
	} else {
		session.PutBytes(0, 0)
	}

	if lob.charsetID != 0 {
		session.PutBytes(1)
	} else {
		session.PutBytes(0)
	}

	if session.TTCVersion < 3 {
		session.PutBytes(1)
	} else {
		session.PutBytes(0)
	}

	// if bNullO2U (false) {
	// session.PutBytes(1)
	//} else {
	session.PutBytes(0)

	session.PutInt(operationID, 4, true, true)
	if len(lob.scn) == 0 {
		session.PutBytes(0)
	} else {
		session.PutBytes(1)
	}
	session.PutUint(len(lob.scn), 4, true, true)

	if session.TTCVersion >= 3 {
		session.PutUint(lob.sourceOffset, 8, true, true)
		session.PutInt(lob.destOffset, 8, true, true)
		// sendAmount
		session.PutBytes(1)
	}
	if session.TTCVersion >= 4 {
		session.PutBytes(0, 0, 0, 0, 0, 0)
	}

	if len(lob.sourceLocator) > 0 {
		session.PutBytes(lob.sourceLocator...)
	}

	if len(lob.destLocator) > 0 {
		session.PutBytes(lob.destLocator...)
	}

	if lob.charsetID != 0 {
		session.PutUint(lob.charsetID, 2, true, true)
	}
	if session.TTCVersion < 3 {
		session.PutUint(lob.size, 4, true, true)
	}
	for x := 0; x < len(lob.scn); x++ {
		session.PutUint(lob.scn[x], 4, true, true)
	}
	if session.TTCVersion >= 3 {
		session.PutUint(lob.size, 8, true, true)
	}
	return session.Write()
}

// read lob response from network session
func (lob *Lob) read(connection *Connection) error {
	loop := true
	session := connection.session
	for loop {
		msg, err := session.GetByte()
		if err != nil {
			return err
		}
		switch msg {
		case 4:
			session.Summary, err = network.NewSummary(session)
			if err != nil {
				return err
			}
			if session.HasError() {
				if session.Summary.RetCode == 1403 {
					session.Summary = nil
				} else {
					return session.GetError()
				}
			}
			loop = false
		case 8:
			// read rpa message
			if len(lob.sourceLocator) != 0 {
				_, err = session.GetBytes(len(lob.sourceLocator))
				if err != nil {
					return err
				}
			}
			if len(lob.destLocator) != 0 {
				_, err = session.GetBytes(len(lob.destLocator))
				if err != nil {
					return err
				}
			}
			if lob.charsetID != 0 {
				lob.charsetID, err = session.GetInt(2, true, true)
				if err != nil {
					return err
				}
			}
			// get datasize
			if session.TTCVersion < 3 {
				lob.size, err = session.GetInt64(4, true, true)
				if err != nil {
					return err
				}
			} else {
				lob.size, err = session.GetInt64(8, true, true)
				if err != nil {
					return err
				}
			}
		case 9:
			if session.HasEOSCapability {
				session.Summary.EndOfCallStatus, err = session.GetInt(4, true, true)
				if err != nil {
					return err
				}
			}
			loop = false
		case 14:
			// get the data
			err = lob.readData(session)
			if err != nil {
				return err
			}
		case 15:
			warning, err := network.NewWarningObject(session)
			if err != nil {
				return err
			}
			if warning != nil {
				fmt.Println(warning)
			}
		case 23:
			opCode, err := session.GetByte()
			if err != nil {
				return err
			}
			err = connection.getServerNetworkInformation(opCode)
			if err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("TTC error: received code %d during LOB reading", msg))
		}
	}
	return nil
}

// read lob data chunks from network session
func (lob *Lob) readData(session *network.Session) error {
	num1 := 0 // data readed in the call of this function
	var chunkSize int = 0
	var err error
	//num3 := offset // the data readed from the start of read operation
	num4 := 0
	for num4 != 4 {
		switch num4 {
		case 0:
			nb, err := session.GetByte()
			if err != nil {
				return err
			}
			chunkSize = int(nb)
			if chunkSize == 0xFE {
				num4 = 2
			} else {
				num4 = 1
			}
		case 1:
			chunk, err := session.GetBytes(chunkSize)
			if err != nil {
				return err
			}
			lob.data.Write(chunk)
			num1 += int(chunkSize)
			num4 = 4
		case 2:
			if session.UseBigClrChunks {
				chunkSize, err = session.GetInt(4, true, true)
			} else {
				var nb byte
				nb, err = session.GetByte()
				chunkSize = int(nb)
			}
			if err != nil {
				return err
			}
			if chunkSize <= 0 {
				num4 = 4
			} else {
				num4 = 3
			}
		case 3:
			chunk, err := session.GetBytes(int(chunkSize))
			if err != nil {
				return err
			}
			lob.data.Write(chunk)
			num1 += int(chunkSize)
			//num3 += chunkSize
			num4 = 2
		}
	}
	return nil
}
