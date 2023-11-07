/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"bytes"
	"crypto/tls"
	"dm/security"
	"fmt"
	"net"
	"strconv"
	"time"
	"unicode/utf8"
)

const (
	Dm_build_408 = 8192
	Dm_build_409 = 2 * time.Second
)

type dm_build_410 struct {
	dm_build_411 *net.TCPConn
	dm_build_412 *tls.Conn
	dm_build_413 *Dm_build_78
	dm_build_414 *DmConnection
	dm_build_415 security.Cipher
	dm_build_416 bool
	dm_build_417 bool
	dm_build_418 *security.DhKey

	dm_build_419 bool
	dm_build_420 string
	dm_build_421 bool
}

func dm_build_422(dm_build_423 *DmConnection) (*dm_build_410, error) {
	dm_build_424, dm_build_425 := dm_build_427(dm_build_423.dmConnector.host+":"+strconv.Itoa(int(dm_build_423.dmConnector.port)), time.Duration(dm_build_423.dmConnector.socketTimeout)*time.Second)
	if dm_build_425 != nil {
		return nil, dm_build_425
	}

	dm_build_426 := dm_build_410{}
	dm_build_426.dm_build_411 = dm_build_424
	dm_build_426.dm_build_413 = Dm_build_81(Dm_build_682)
	dm_build_426.dm_build_414 = dm_build_423
	dm_build_426.dm_build_416 = false
	dm_build_426.dm_build_417 = false
	dm_build_426.dm_build_419 = false
	dm_build_426.dm_build_420 = ""
	dm_build_426.dm_build_421 = false
	dm_build_423.Access = &dm_build_426

	return &dm_build_426, nil
}

func dm_build_427(dm_build_428 string, dm_build_429 time.Duration) (*net.TCPConn, error) {
	dm_build_430, dm_build_431 := net.DialTimeout("tcp", dm_build_428, dm_build_429)
	if dm_build_431 != nil {
		return nil, ECGO_COMMUNITION_ERROR.addDetail("\tdial address: " + dm_build_428).throw()
	}

	if tcpConn, ok := dm_build_430.(*net.TCPConn); ok {

		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(Dm_build_409)
		tcpConn.SetNoDelay(true)

		return tcpConn, nil
	}

	return nil, nil
}

func (dm_build_433 *dm_build_410) dm_build_432(dm_build_434 dm_build_803) bool {
	var dm_build_435 = dm_build_433.dm_build_414.dmConnector.compress
	if dm_build_434.dm_build_818() == Dm_build_710 || dm_build_435 == Dm_build_759 {
		return false
	}

	if dm_build_435 == Dm_build_757 {
		return true
	} else if dm_build_435 == Dm_build_758 {
		return !dm_build_433.dm_build_414.Local && dm_build_434.dm_build_816() > Dm_build_756
	}

	return false
}

func (dm_build_437 *dm_build_410) dm_build_436(dm_build_438 dm_build_803) bool {
	var dm_build_439 = dm_build_437.dm_build_414.dmConnector.compress
	if dm_build_438.dm_build_818() == Dm_build_710 || dm_build_439 == Dm_build_759 {
		return false
	}

	if dm_build_439 == Dm_build_757 {
		return true
	} else if dm_build_439 == Dm_build_758 {
		return dm_build_437.dm_build_413.Dm_build_341(Dm_build_718) == 1
	}

	return false
}

func (dm_build_441 *dm_build_410) dm_build_440(dm_build_442 dm_build_803) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(string); ok {
				err = ECGO_COMMUNITION_ERROR.addDetail("\t" + p.(string)).throw()
			} else {
				err = fmt.Errorf("internal error: %v", p)
			}
		}
	}()

	dm_build_444 := dm_build_442.dm_build_816()

	if dm_build_444 > 0 {

		if dm_build_441.dm_build_432(dm_build_442) {
			var retBytes, err = Compress(dm_build_441.dm_build_413, Dm_build_711, int(dm_build_444), int(dm_build_441.dm_build_414.dmConnector.compressID))
			if err != nil {
				return err
			}

			dm_build_441.dm_build_413.Dm_build_92(Dm_build_711)

			dm_build_441.dm_build_413.Dm_build_129(dm_build_444)

			dm_build_441.dm_build_413.Dm_build_157(retBytes)

			dm_build_442.dm_build_817(int32(len(retBytes)) + ULINT_SIZE)

			dm_build_441.dm_build_413.Dm_build_261(Dm_build_718, 1)
		}

		if dm_build_441.dm_build_417 {
			dm_build_444 = dm_build_442.dm_build_816()
			var retBytes = dm_build_441.dm_build_415.Encrypt(dm_build_441.dm_build_413.Dm_build_368(Dm_build_711, int(dm_build_444)), true)

			dm_build_441.dm_build_413.Dm_build_92(Dm_build_711)

			dm_build_441.dm_build_413.Dm_build_157(retBytes)

			dm_build_442.dm_build_817(int32(len(retBytes)))
		}
	}

	if dm_build_441.dm_build_413.Dm_build_90() > Dm_build_683 {
		return ECGO_MSG_TOO_LONG.throw()
	}

	dm_build_442.dm_build_812()
	if dm_build_441.dm_build_672(dm_build_442) {
		if dm_build_441.dm_build_412 != nil {
			dm_build_441.dm_build_413.Dm_build_95(0)
			if _, err := dm_build_441.dm_build_413.Dm_build_114(dm_build_441.dm_build_412); err != nil {
				return err
			}
		}
	} else {
		dm_build_441.dm_build_413.Dm_build_95(0)
		if _, err := dm_build_441.dm_build_413.Dm_build_114(dm_build_441.dm_build_411); err != nil {
			return err
		}
	}
	return nil
}

func (dm_build_446 *dm_build_410) dm_build_445(dm_build_447 dm_build_803) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(string); ok {
				err = ECGO_COMMUNITION_ERROR.addDetail("\t" + p.(string)).throw()
			} else {
				err = fmt.Errorf("internal error: %v", p)
			}
		}
	}()

	dm_build_449 := int32(0)
	if dm_build_446.dm_build_672(dm_build_447) {
		if dm_build_446.dm_build_412 != nil {
			dm_build_446.dm_build_413.Dm_build_92(0)
			if _, err := dm_build_446.dm_build_413.Dm_build_108(dm_build_446.dm_build_412, Dm_build_711); err != nil {
				return err
			}

			dm_build_449 = dm_build_447.dm_build_816()
			if dm_build_449 > 0 {
				if _, err := dm_build_446.dm_build_413.Dm_build_108(dm_build_446.dm_build_412, int(dm_build_449)); err != nil {
					return err
				}
			}
		}
	} else {

		dm_build_446.dm_build_413.Dm_build_92(0)
		if _, err := dm_build_446.dm_build_413.Dm_build_108(dm_build_446.dm_build_411, Dm_build_711); err != nil {
			return err
		}
		dm_build_449 = dm_build_447.dm_build_816()

		if dm_build_449 > 0 {
			if _, err := dm_build_446.dm_build_413.Dm_build_108(dm_build_446.dm_build_411, int(dm_build_449)); err != nil {
				return err
			}
		}
	}

	dm_build_447.dm_build_813()

	dm_build_449 = dm_build_447.dm_build_816()
	if dm_build_449 <= 0 {
		return nil
	}

	if dm_build_446.dm_build_417 {
		ebytes := dm_build_446.dm_build_413.Dm_build_368(Dm_build_711, int(dm_build_449))
		bytes, err := dm_build_446.dm_build_415.Decrypt(ebytes, true)
		if err != nil {
			return err
		}
		dm_build_446.dm_build_413.Dm_build_92(Dm_build_711)
		dm_build_446.dm_build_413.Dm_build_157(bytes)
		dm_build_447.dm_build_817(int32(len(bytes)))
	}

	if dm_build_446.dm_build_436(dm_build_447) {

		dm_build_449 = dm_build_447.dm_build_816()
		cbytes := dm_build_446.dm_build_413.Dm_build_368(Dm_build_711+ULINT_SIZE, int(dm_build_449-ULINT_SIZE))
		bytes, err := UnCompress(cbytes, int(dm_build_446.dm_build_414.dmConnector.compressID))
		if err != nil {
			return err
		}
		dm_build_446.dm_build_413.Dm_build_92(Dm_build_711)
		dm_build_446.dm_build_413.Dm_build_157(bytes)
		dm_build_447.dm_build_817(int32(len(bytes)))
	}
	return nil
}

func (dm_build_451 *dm_build_410) dm_build_450(dm_build_452 dm_build_803) (dm_build_453 interface{}, dm_build_454 error) {
	dm_build_454 = dm_build_452.dm_build_807(dm_build_452)
	if dm_build_454 != nil {
		return nil, dm_build_454
	}

	dm_build_454 = dm_build_451.dm_build_440(dm_build_452)
	if dm_build_454 != nil {
		return nil, dm_build_454
	}

	dm_build_454 = dm_build_451.dm_build_445(dm_build_452)
	if dm_build_454 != nil {
		return nil, dm_build_454
	}

	return dm_build_452.dm_build_811(dm_build_452)
}

func (dm_build_456 *dm_build_410) dm_build_455() (*dm_build_1240, error) {

	Dm_build_457 := dm_build_1246(dm_build_456)
	_, dm_build_458 := dm_build_456.dm_build_450(Dm_build_457)
	if dm_build_458 != nil {
		return nil, dm_build_458
	}

	return Dm_build_457, nil
}

func (dm_build_460 *dm_build_410) dm_build_459() error {

	dm_build_461 := dm_build_1108(dm_build_460)
	_, dm_build_462 := dm_build_460.dm_build_450(dm_build_461)
	if dm_build_462 != nil {
		return dm_build_462
	}

	return nil
}

func (dm_build_464 *dm_build_410) dm_build_463() error {

	var dm_build_465 *dm_build_1240
	var err error
	if dm_build_465, err = dm_build_464.dm_build_455(); err != nil {
		return err
	}

	if dm_build_464.dm_build_414.sslEncrypt == 2 {
		if err = dm_build_464.dm_build_668(false); err != nil {
			return ECGO_INIT_SSL_FAILED.addDetail("\n" + err.Error()).throw()
		}
	} else if dm_build_464.dm_build_414.sslEncrypt == 1 {
		if err = dm_build_464.dm_build_668(true); err != nil {
			return ECGO_INIT_SSL_FAILED.addDetail("\n" + err.Error()).throw()
		}
	}

	if dm_build_464.dm_build_417 || dm_build_464.dm_build_416 {
		k, err := dm_build_464.dm_build_658()
		if err != nil {
			return err
		}
		sessionKey := security.ComputeSessionKey(k, dm_build_465.Dm_build_1244)
		encryptType := dm_build_465.dm_build_1242
		hashType := int(dm_build_465.Dm_build_1243)
		if encryptType == -1 {
			encryptType = security.DES_CFB
		}
		if hashType == -1 {
			hashType = security.MD5
		}
		err = dm_build_464.dm_build_661(encryptType, sessionKey, dm_build_464.dm_build_414.dmConnector.cipherPath, hashType)
		if err != nil {
			return err
		}
	}

	if err := dm_build_464.dm_build_459(); err != nil {
		return err
	}
	return nil
}

func (dm_build_468 *dm_build_410) Dm_build_467(dm_build_469 *DmStatement) error {
	dm_build_470 := dm_build_1269(dm_build_468, dm_build_469)
	_, dm_build_471 := dm_build_468.dm_build_450(dm_build_470)
	if dm_build_471 != nil {
		return dm_build_471
	}

	return nil
}

func (dm_build_473 *dm_build_410) Dm_build_472(dm_build_474 int32) error {
	dm_build_475 := dm_build_1279(dm_build_473, dm_build_474)
	_, dm_build_476 := dm_build_473.dm_build_450(dm_build_475)
	if dm_build_476 != nil {
		return dm_build_476
	}

	return nil
}

func (dm_build_478 *dm_build_410) Dm_build_477(dm_build_479 *DmStatement, dm_build_480 bool, dm_build_481 int16) (*execRetInfo, error) {
	dm_build_482 := dm_build_1146(dm_build_478, dm_build_479, dm_build_480, dm_build_481)
	dm_build_483, dm_build_484 := dm_build_478.dm_build_450(dm_build_482)
	if dm_build_484 != nil {
		return nil, dm_build_484
	}
	return dm_build_483.(*execRetInfo), nil
}

func (dm_build_486 *dm_build_410) Dm_build_485(dm_build_487 *DmStatement, dm_build_488 int16) (*execRetInfo, error) {
	return dm_build_486.Dm_build_477(dm_build_487, false, Dm_build_763)
}

func (dm_build_490 *dm_build_410) Dm_build_489(dm_build_491 *DmStatement, dm_build_492 []OptParameter) (*execRetInfo, error) {
	dm_build_493, dm_build_494 := dm_build_490.dm_build_450(dm_build_905(dm_build_490, dm_build_491, dm_build_492))
	if dm_build_494 != nil {
		return nil, dm_build_494
	}

	return dm_build_493.(*execRetInfo), nil
}

func (dm_build_496 *dm_build_410) Dm_build_495(dm_build_497 *DmStatement, dm_build_498 int16) (*execRetInfo, error) {
	return dm_build_496.Dm_build_477(dm_build_497, true, dm_build_498)
}

func (dm_build_500 *dm_build_410) Dm_build_499(dm_build_501 *DmStatement, dm_build_502 [][]interface{}) (*execRetInfo, error) {
	dm_build_503 := dm_build_928(dm_build_500, dm_build_501, dm_build_502)
	dm_build_504, dm_build_505 := dm_build_500.dm_build_450(dm_build_503)
	if dm_build_505 != nil {
		return nil, dm_build_505
	}
	return dm_build_504.(*execRetInfo), nil
}

func (dm_build_507 *dm_build_410) Dm_build_506(dm_build_508 *DmStatement, dm_build_509 [][]interface{}, dm_build_510 bool) (*execRetInfo, error) {
	var dm_build_511, dm_build_512 = 0, 0
	var dm_build_513 = len(dm_build_509)
	var dm_build_514 [][]interface{}
	var dm_build_515 = NewExceInfo()
	dm_build_515.updateCounts = make([]int64, dm_build_513)
	var dm_build_516 = false
	for dm_build_511 < dm_build_513 {
		for dm_build_512 = dm_build_511; dm_build_512 < dm_build_513; dm_build_512++ {
			paramData := dm_build_509[dm_build_512]
			bindData := make([]interface{}, dm_build_508.paramCount)
			dm_build_516 = false
			for icol := 0; icol < int(dm_build_508.paramCount); icol++ {
				if dm_build_508.params[icol].ioType == IO_TYPE_OUT {
					continue
				}
				if dm_build_507.dm_build_641(bindData, paramData, icol) {
					dm_build_516 = true
					break
				}
			}

			if dm_build_516 {
				break
			}
			dm_build_514 = append(dm_build_514, bindData)
		}

		if dm_build_512 != dm_build_511 {
			tmpExecInfo, err := dm_build_507.Dm_build_499(dm_build_508, dm_build_514)
			if err != nil {
				return nil, err
			}
			dm_build_514 = dm_build_514[0:0]
			dm_build_515.union(tmpExecInfo, dm_build_511, dm_build_512-dm_build_511)
		}

		if dm_build_512 < dm_build_513 {
			tmpExecInfo, err := dm_build_507.Dm_build_517(dm_build_508, dm_build_509[dm_build_512], dm_build_510)
			if err != nil {
				return nil, err
			}

			dm_build_510 = true
			dm_build_515.union(tmpExecInfo, dm_build_512, 1)
		}

		dm_build_511 = dm_build_512 + 1
	}
	for _, i := range dm_build_515.updateCounts {
		if i > 0 {
			dm_build_515.updateCount += i
		}
	}
	return dm_build_515, nil
}

func (dm_build_518 *dm_build_410) Dm_build_517(dm_build_519 *DmStatement, dm_build_520 []interface{}, dm_build_521 bool) (*execRetInfo, error) {

	var dm_build_522 = make([]interface{}, dm_build_519.paramCount)
	for icol := 0; icol < int(dm_build_519.paramCount); icol++ {
		if dm_build_519.params[icol].ioType == IO_TYPE_OUT {
			continue
		}
		if dm_build_518.dm_build_641(dm_build_522, dm_build_520, icol) {

			if !dm_build_521 {
				preExecute := dm_build_1136(dm_build_518, dm_build_519, dm_build_519.params)
				dm_build_518.dm_build_450(preExecute)
				dm_build_521 = true
			}

			dm_build_518.dm_build_647(dm_build_519, dm_build_519.params[icol], icol, dm_build_520[icol].(iOffRowBinder))
			dm_build_522[icol] = ParamDataEnum_OFF_ROW
		}
	}

	var dm_build_523 = make([][]interface{}, 1, 1)
	dm_build_523[0] = dm_build_522

	dm_build_524 := dm_build_928(dm_build_518, dm_build_519, dm_build_523)
	dm_build_525, dm_build_526 := dm_build_518.dm_build_450(dm_build_524)
	if dm_build_526 != nil {
		return nil, dm_build_526
	}
	return dm_build_525.(*execRetInfo), nil
}

func (dm_build_528 *dm_build_410) Dm_build_527(dm_build_529 *DmStatement, dm_build_530 int16) (*execRetInfo, error) {
	dm_build_531 := dm_build_1123(dm_build_528, dm_build_529, dm_build_530)

	dm_build_532, dm_build_533 := dm_build_528.dm_build_450(dm_build_531)
	if dm_build_533 != nil {
		return nil, dm_build_533
	}
	return dm_build_532.(*execRetInfo), nil
}

func (dm_build_535 *dm_build_410) Dm_build_534(dm_build_536 *innerRows, dm_build_537 int64) (*execRetInfo, error) {
	dm_build_538 := dm_build_1028(dm_build_535, dm_build_536, dm_build_537, INT64_MAX)
	dm_build_539, dm_build_540 := dm_build_535.dm_build_450(dm_build_538)
	if dm_build_540 != nil {
		return nil, dm_build_540
	}
	return dm_build_539.(*execRetInfo), nil
}

func (dm_build_542 *dm_build_410) Commit() error {
	dm_build_543 := dm_build_891(dm_build_542)
	_, dm_build_544 := dm_build_542.dm_build_450(dm_build_543)
	if dm_build_544 != nil {
		return dm_build_544
	}

	return nil
}

func (dm_build_546 *dm_build_410) Rollback() error {
	dm_build_547 := dm_build_1184(dm_build_546)
	_, dm_build_548 := dm_build_546.dm_build_450(dm_build_547)
	if dm_build_548 != nil {
		return dm_build_548
	}

	return nil
}

func (dm_build_550 *dm_build_410) Dm_build_549(dm_build_551 *DmConnection) error {
	dm_build_552 := dm_build_1189(dm_build_550, dm_build_551.IsoLevel)
	_, dm_build_553 := dm_build_550.dm_build_450(dm_build_552)
	if dm_build_553 != nil {
		return dm_build_553
	}

	return nil
}

func (dm_build_555 *dm_build_410) Dm_build_554(dm_build_556 *DmStatement, dm_build_557 string) error {
	dm_build_558 := dm_build_896(dm_build_555, dm_build_556, dm_build_557)
	_, dm_build_559 := dm_build_555.dm_build_450(dm_build_558)
	if dm_build_559 != nil {
		return dm_build_559
	}

	return nil
}

func (dm_build_561 *dm_build_410) Dm_build_560(dm_build_562 []uint32) ([]int64, error) {
	dm_build_563 := dm_build_1287(dm_build_561, dm_build_562)
	dm_build_564, dm_build_565 := dm_build_561.dm_build_450(dm_build_563)
	if dm_build_565 != nil {
		return nil, dm_build_565
	}
	return dm_build_564.([]int64), nil
}

func (dm_build_567 *dm_build_410) Close() error {
	if dm_build_567.dm_build_421 {
		return nil
	}

	dm_build_568 := dm_build_567.dm_build_411.Close()
	if dm_build_568 != nil {
		return dm_build_568
	}

	dm_build_567.dm_build_414 = nil
	dm_build_567.dm_build_421 = true
	return nil
}

func (dm_build_570 *dm_build_410) dm_build_569(dm_build_571 *lob) (int64, error) {
	dm_build_572 := dm_build_1059(dm_build_570, dm_build_571)
	dm_build_573, dm_build_574 := dm_build_570.dm_build_450(dm_build_572)
	if dm_build_574 != nil {
		return 0, dm_build_574
	}
	return dm_build_573.(int64), nil
}

func (dm_build_576 *dm_build_410) dm_build_575(dm_build_577 *lob, dm_build_578 int32, dm_build_579 int32) ([]byte, error) {
	dm_build_580 := dm_build_1046(dm_build_576, dm_build_577, int(dm_build_578), int(dm_build_579))
	dm_build_581, dm_build_582 := dm_build_576.dm_build_450(dm_build_580)
	if dm_build_582 != nil {
		return nil, dm_build_582
	}
	return dm_build_581.([]byte), nil
}

func (dm_build_584 *dm_build_410) dm_build_583(dm_build_585 *DmBlob, dm_build_586 int32, dm_build_587 int32) ([]byte, error) {
	var dm_build_588 = make([]byte, dm_build_587)
	var dm_build_589 int32 = 0
	var dm_build_590 int32 = 0
	var dm_build_591 []byte
	var dm_build_592 error
	for dm_build_589 < dm_build_587 {
		dm_build_590 = dm_build_587 - dm_build_589
		if dm_build_590 > Dm_build_796 {
			dm_build_590 = Dm_build_796
		}
		dm_build_591, dm_build_592 = dm_build_584.dm_build_575(&dm_build_585.lob, dm_build_586, dm_build_590)
		if dm_build_592 != nil {
			return nil, dm_build_592
		}
		if dm_build_591 == nil || len(dm_build_591) == 0 {
			break
		}
		Dm_build_1298.Dm_build_1354(dm_build_588, int(dm_build_589), dm_build_591, 0, len(dm_build_591))
		dm_build_589 += int32(len(dm_build_591))
		dm_build_586 += int32(len(dm_build_591))
		if dm_build_585.readOver {
			break
		}
	}
	return dm_build_588, nil
}

func (dm_build_594 *dm_build_410) dm_build_593(dm_build_595 *DmClob, dm_build_596 int32, dm_build_597 int32) (string, error) {
	var dm_build_598 bytes.Buffer
	var dm_build_599 int32 = 0
	var dm_build_600 int32 = 0
	var dm_build_601 []byte
	var dm_build_602 string
	var dm_build_603 error
	for dm_build_599 < dm_build_597 {
		dm_build_600 = dm_build_597 - dm_build_599
		if dm_build_600 > Dm_build_796/2 {
			dm_build_600 = Dm_build_796 / 2
		}
		dm_build_601, dm_build_603 = dm_build_594.dm_build_575(&dm_build_595.lob, dm_build_596, dm_build_600)
		if dm_build_603 != nil {
			return "", dm_build_603
		}
		if dm_build_601 == nil || len(dm_build_601) == 0 {
			break
		}
		dm_build_602 = Dm_build_1298.Dm_build_1455(dm_build_601, 0, len(dm_build_601), dm_build_595.serverEncoding, dm_build_594.dm_build_414)

		dm_build_598.WriteString(dm_build_602)
		strLen := utf8.RuneCountInString(dm_build_602)
		dm_build_599 += int32(strLen)
		dm_build_596 += int32(strLen)
		if dm_build_595.readOver {
			break
		}
	}
	return dm_build_598.String(), nil
}

func (dm_build_605 *dm_build_410) dm_build_604(dm_build_606 *DmClob, dm_build_607 int, dm_build_608 string, dm_build_609 string) (int, error) {
	var dm_build_610 = Dm_build_1298.Dm_build_1511(dm_build_608, dm_build_609, dm_build_605.dm_build_414)
	var dm_build_611 = 0
	var dm_build_612 = len(dm_build_610)
	var dm_build_613 = 0
	var dm_build_614 = 0
	var dm_build_615 = 0
	var dm_build_616 = dm_build_612/Dm_build_795 + 1
	var dm_build_617 byte = 0
	var dm_build_618 byte = 0x01
	var dm_build_619 byte = 0x02
	for i := 0; i < dm_build_616; i++ {
		dm_build_617 = 0
		if i == 0 {
			dm_build_617 |= dm_build_618
		}
		if i == dm_build_616-1 {
			dm_build_617 |= dm_build_619
		}
		dm_build_615 = dm_build_612 - dm_build_614
		if dm_build_615 > Dm_build_795 {
			dm_build_615 = Dm_build_795
		}

		setLobData := dm_build_1203(dm_build_605, &dm_build_606.lob, dm_build_617, dm_build_607, dm_build_610, dm_build_611, dm_build_615)
		ret, err := dm_build_605.dm_build_450(setLobData)
		if err != nil {
			return 0, err
		}
		tmp := ret.(int32)
		if err != nil {
			return -1, err
		}
		if tmp <= 0 {
			return dm_build_613, nil
		} else {
			dm_build_607 += int(tmp)
			dm_build_613 += int(tmp)
			dm_build_614 += dm_build_615
			dm_build_611 += dm_build_615
		}
	}
	return dm_build_613, nil
}

func (dm_build_621 *dm_build_410) dm_build_620(dm_build_622 *DmBlob, dm_build_623 int, dm_build_624 []byte) (int, error) {
	var dm_build_625 = 0
	var dm_build_626 = len(dm_build_624)
	var dm_build_627 = 0
	var dm_build_628 = 0
	var dm_build_629 = 0
	var dm_build_630 = dm_build_626/Dm_build_795 + 1
	var dm_build_631 byte = 0
	var dm_build_632 byte = 0x01
	var dm_build_633 byte = 0x02
	for i := 0; i < dm_build_630; i++ {
		dm_build_631 = 0
		if i == 0 {
			dm_build_631 |= dm_build_632
		}
		if i == dm_build_630-1 {
			dm_build_631 |= dm_build_633
		}
		dm_build_629 = dm_build_626 - dm_build_628
		if dm_build_629 > Dm_build_795 {
			dm_build_629 = Dm_build_795
		}

		setLobData := dm_build_1203(dm_build_621, &dm_build_622.lob, dm_build_631, dm_build_623, dm_build_624, dm_build_625, dm_build_629)
		ret, err := dm_build_621.dm_build_450(setLobData)
		if err != nil {
			return 0, err
		}
		tmp := ret.(int32)
		if tmp <= 0 {
			return dm_build_627, nil
		} else {
			dm_build_623 += int(tmp)
			dm_build_627 += int(tmp)
			dm_build_628 += dm_build_629
			dm_build_625 += dm_build_629
		}
	}
	return dm_build_627, nil
}

func (dm_build_635 *dm_build_410) dm_build_634(dm_build_636 *lob, dm_build_637 int) (int64, error) {
	dm_build_638 := dm_build_1070(dm_build_635, dm_build_636, dm_build_637)
	dm_build_639, dm_build_640 := dm_build_635.dm_build_450(dm_build_638)
	if dm_build_640 != nil {
		return dm_build_636.length, dm_build_640
	}
	return dm_build_639.(int64), nil
}

func (dm_build_642 *dm_build_410) dm_build_641(dm_build_643 []interface{}, dm_build_644 []interface{}, dm_build_645 int) bool {
	var dm_build_646 = false
	if dm_build_645 >= len(dm_build_644) || dm_build_644[dm_build_645] == nil {
		dm_build_643[dm_build_645] = ParamDataEnum_Null
	} else if binder, ok := dm_build_644[dm_build_645].(iOffRowBinder); ok {
		dm_build_646 = true
		dm_build_643[dm_build_645] = ParamDataEnum_OFF_ROW
		var lob lob
		if l, ok := binder.getObj().(DmBlob); ok {
			lob = l.lob
		} else if l, ok := binder.getObj().(DmClob); ok {
			lob = l.lob
		}
		if &lob != nil && lob.canOptimized(dm_build_642.dm_build_414) {
			dm_build_643[dm_build_645] = &lobCtl{lob.buildCtlData()}
			dm_build_646 = false
		}
	} else {
		dm_build_643[dm_build_645] = dm_build_644[dm_build_645]
	}
	return dm_build_646
}

func (dm_build_648 *dm_build_410) dm_build_647(dm_build_649 *DmStatement, dm_build_650 parameter, dm_build_651 int, dm_build_652 iOffRowBinder) error {
	var dm_build_653 = Dm_build_4()
	dm_build_652.read(dm_build_653)
	var dm_build_654 = 0
	for !dm_build_652.isReadOver() || dm_build_653.Dm_build_5() > 0 {
		if !dm_build_652.isReadOver() && dm_build_653.Dm_build_5() < Dm_build_795 {
			dm_build_652.read(dm_build_653)
		}
		if dm_build_653.Dm_build_5() > Dm_build_795 {
			dm_build_654 = Dm_build_795
		} else {
			dm_build_654 = dm_build_653.Dm_build_5()
		}

		putData := dm_build_1174(dm_build_648, dm_build_649, int16(dm_build_651), dm_build_653, int32(dm_build_654))
		_, err := dm_build_648.dm_build_450(putData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dm_build_656 *dm_build_410) dm_build_655() ([]byte, error) {
	var dm_build_657 error
	if dm_build_656.dm_build_418 == nil {
		if dm_build_656.dm_build_418, dm_build_657 = security.NewClientKeyPair(); dm_build_657 != nil {
			return nil, dm_build_657
		}
	}
	return security.Bn2Bytes(dm_build_656.dm_build_418.GetY(), security.DH_KEY_LENGTH), nil
}

func (dm_build_659 *dm_build_410) dm_build_658() (*security.DhKey, error) {
	var dm_build_660 error
	if dm_build_659.dm_build_418 == nil {
		if dm_build_659.dm_build_418, dm_build_660 = security.NewClientKeyPair(); dm_build_660 != nil {
			return nil, dm_build_660
		}
	}
	return dm_build_659.dm_build_418, nil
}

func (dm_build_662 *dm_build_410) dm_build_661(dm_build_663 int, dm_build_664 []byte, dm_build_665 string, dm_build_666 int) (dm_build_667 error) {
	if dm_build_663 > 0 && dm_build_663 < security.MIN_EXTERNAL_CIPHER_ID && dm_build_664 != nil {
		dm_build_662.dm_build_415, dm_build_667 = security.NewSymmCipher(dm_build_663, dm_build_664)
	} else if dm_build_663 >= security.MIN_EXTERNAL_CIPHER_ID {
		if dm_build_662.dm_build_415, dm_build_667 = security.NewThirdPartCipher(dm_build_663, dm_build_664, dm_build_665, dm_build_666); dm_build_667 != nil {
			dm_build_667 = THIRD_PART_CIPHER_INIT_FAILED.addDetailln(dm_build_667.Error()).throw()
		}
	}
	return
}

func (dm_build_669 *dm_build_410) dm_build_668(dm_build_670 bool) (dm_build_671 error) {
	if dm_build_669.dm_build_412, dm_build_671 = security.NewTLSFromTCP(dm_build_669.dm_build_411, dm_build_669.dm_build_414.dmConnector.sslCertPath, dm_build_669.dm_build_414.dmConnector.sslKeyPath, dm_build_669.dm_build_414.dmConnector.user); dm_build_671 != nil {
		return
	}
	if !dm_build_670 {
		dm_build_669.dm_build_412 = nil
	}
	return
}

func (dm_build_673 *dm_build_410) dm_build_672(dm_build_674 dm_build_803) bool {
	return dm_build_674.dm_build_818() != Dm_build_710 && dm_build_673.dm_build_414.sslEncrypt == 1
}
