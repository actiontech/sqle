package advanced_nego

import (
	"errors"
	"fmt"
	"github.com/sijms/go-ora/v2/network/security"
)

type encryptService struct {
	defaultService
	algoID int
}

func NewEncryptService(comm *AdvancedNegoComm) (*encryptService, error) {
	output := &encryptService{
		defaultService: defaultService{
			comm:        comm,
			serviceType: 2,
			version:     0xB200200,
			availableServiceNames: []string{"", "RC4_40", "RC4_56", "RC4_128", "RC4_256",
				"DES40C", "DES56C", "3DES112", "3DES168", "AES128", "AES192", "AES256"},
			availableServiceIDs: []int{0, 1, 8, 10, 6, 3, 2, 11, 12, 15, 16, 17},
		},
	}
	str := ""
	level := ""
	connOption := comm.session.Context.ConnOption
	if connOption != nil {
		snConfig := connOption.SNOConfig
		if snConfig != nil {
			var exists bool
			str, exists = snConfig["sqlnet.encryption_types_client"]
			if !exists {
				str = ""
			}
			level, exists = snConfig["sqlnet.encryption_client"]
			if !exists {
				level = ""
			}
		}
	}
	output.readAdvNegoLevel(level)
	//level := conops.Encryption != null ? conops.Encryption : snoConfig[];
	err := output.buildServiceList(str, true, true)
	//output.selectedServ, err = output.validate(strings.Split(str,","), true)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (serv *encryptService) readServiceData(subPacketnum int) error {
	var err error
	comm := serv.comm
	serv.version, err = comm.readVersion()
	if err != nil {
		return err
	}
	resp, err := comm.readUB1()
	if err != nil {
		return err
	}
	serv.algoID = int(resp)

	return nil
}
func (serv *encryptService) writeServiceData() error {
	serv.writeHeader(3)
	comm := serv.comm
	comm.writeVersion(serv.getVersion())
	selectedIndices := make([]byte, len(serv.selectedIndices))
	for i := 0; i < len(serv.selectedIndices); i++ {
		index := serv.selectedIndices[i]
		selectedIndices[i] = uint8(serv.availableServiceIDs[index])
	}
	comm.writeBytes(selectedIndices)
	// send selected driver
	comm.writeUB1(1)
	return nil
}

func (serv *encryptService) getServiceDataLength() int {
	return 17 + len(serv.selectedIndices)
}

func (serv *encryptService) activateAlgorithm() error {
	key := serv.comm.session.Context.AdvancedService.SessionKey
	iv := make([]byte, 16)
	switch serv.algoID {
	case 0:
		return nil
	case 15:
		algo, err := security.NewOracleNetworkCBCEncrypter(key[:16], iv)
		if err != nil {
			return err
		}
		serv.comm.session.Context.AdvancedService.CryptAlgo = algo
	case 16:
		algo, err := security.NewOracleNetworkCBCEncrypter(key[:24], iv)
		if err != nil {
			return err
		}
		serv.comm.session.Context.AdvancedService.CryptAlgo = algo
	case 17:
		algo, err := security.NewOracleNetworkCBCEncrypter(key[:32], iv)
		if err != nil {
			return err
		}
		serv.comm.session.Context.AdvancedService.CryptAlgo = algo
	default:
		return errors.New(fmt.Sprintf("advanced negotiation error: encryption service algorithm: %d still not supported", serv.algoID))
	}
	return nil
}
