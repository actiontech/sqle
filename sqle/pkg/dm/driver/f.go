/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"math"
)

type dm_build_1297 struct{}

var Dm_build_1298 = &dm_build_1297{}

func (Dm_build_1300 *dm_build_1297) Dm_build_1299(dm_build_1301 []byte, dm_build_1302 int, dm_build_1303 byte) int {
	dm_build_1301[dm_build_1302] = dm_build_1303
	return 1
}

func (Dm_build_1305 *dm_build_1297) Dm_build_1304(dm_build_1306 []byte, dm_build_1307 int, dm_build_1308 int8) int {
	dm_build_1306[dm_build_1307] = byte(dm_build_1308)
	return 1
}

func (Dm_build_1310 *dm_build_1297) Dm_build_1309(dm_build_1311 []byte, dm_build_1312 int, dm_build_1313 int16) int {
	dm_build_1311[dm_build_1312] = byte(dm_build_1313)
	dm_build_1312++
	dm_build_1311[dm_build_1312] = byte(dm_build_1313 >> 8)
	return 2
}

func (Dm_build_1315 *dm_build_1297) Dm_build_1314(dm_build_1316 []byte, dm_build_1317 int, dm_build_1318 int32) int {
	dm_build_1316[dm_build_1317] = byte(dm_build_1318)
	dm_build_1317++
	dm_build_1316[dm_build_1317] = byte(dm_build_1318 >> 8)
	dm_build_1317++
	dm_build_1316[dm_build_1317] = byte(dm_build_1318 >> 16)
	dm_build_1317++
	dm_build_1316[dm_build_1317] = byte(dm_build_1318 >> 24)
	dm_build_1317++
	return 4
}

func (Dm_build_1320 *dm_build_1297) Dm_build_1319(dm_build_1321 []byte, dm_build_1322 int, dm_build_1323 int64) int {
	dm_build_1321[dm_build_1322] = byte(dm_build_1323)
	dm_build_1322++
	dm_build_1321[dm_build_1322] = byte(dm_build_1323 >> 8)
	dm_build_1322++
	dm_build_1321[dm_build_1322] = byte(dm_build_1323 >> 16)
	dm_build_1322++
	dm_build_1321[dm_build_1322] = byte(dm_build_1323 >> 24)
	dm_build_1322++
	dm_build_1321[dm_build_1322] = byte(dm_build_1323 >> 32)
	dm_build_1322++
	dm_build_1321[dm_build_1322] = byte(dm_build_1323 >> 40)
	dm_build_1322++
	dm_build_1321[dm_build_1322] = byte(dm_build_1323 >> 48)
	dm_build_1322++
	dm_build_1321[dm_build_1322] = byte(dm_build_1323 >> 56)
	return 8
}

func (Dm_build_1325 *dm_build_1297) Dm_build_1324(dm_build_1326 []byte, dm_build_1327 int, dm_build_1328 float32) int {
	return Dm_build_1325.Dm_build_1344(dm_build_1326, dm_build_1327, math.Float32bits(dm_build_1328))
}

func (Dm_build_1330 *dm_build_1297) Dm_build_1329(dm_build_1331 []byte, dm_build_1332 int, dm_build_1333 float64) int {
	return Dm_build_1330.Dm_build_1349(dm_build_1331, dm_build_1332, math.Float64bits(dm_build_1333))
}

func (Dm_build_1335 *dm_build_1297) Dm_build_1334(dm_build_1336 []byte, dm_build_1337 int, dm_build_1338 uint8) int {
	dm_build_1336[dm_build_1337] = byte(dm_build_1338)
	return 1
}

func (Dm_build_1340 *dm_build_1297) Dm_build_1339(dm_build_1341 []byte, dm_build_1342 int, dm_build_1343 uint16) int {
	dm_build_1341[dm_build_1342] = byte(dm_build_1343)
	dm_build_1342++
	dm_build_1341[dm_build_1342] = byte(dm_build_1343 >> 8)
	return 2
}

func (Dm_build_1345 *dm_build_1297) Dm_build_1344(dm_build_1346 []byte, dm_build_1347 int, dm_build_1348 uint32) int {
	dm_build_1346[dm_build_1347] = byte(dm_build_1348)
	dm_build_1347++
	dm_build_1346[dm_build_1347] = byte(dm_build_1348 >> 8)
	dm_build_1347++
	dm_build_1346[dm_build_1347] = byte(dm_build_1348 >> 16)
	dm_build_1347++
	dm_build_1346[dm_build_1347] = byte(dm_build_1348 >> 24)
	return 3
}

func (Dm_build_1350 *dm_build_1297) Dm_build_1349(dm_build_1351 []byte, dm_build_1352 int, dm_build_1353 uint64) int {
	dm_build_1351[dm_build_1352] = byte(dm_build_1353)
	dm_build_1352++
	dm_build_1351[dm_build_1352] = byte(dm_build_1353 >> 8)
	dm_build_1352++
	dm_build_1351[dm_build_1352] = byte(dm_build_1353 >> 16)
	dm_build_1352++
	dm_build_1351[dm_build_1352] = byte(dm_build_1353 >> 24)
	dm_build_1352++
	dm_build_1351[dm_build_1352] = byte(dm_build_1353 >> 32)
	dm_build_1352++
	dm_build_1351[dm_build_1352] = byte(dm_build_1353 >> 40)
	dm_build_1352++
	dm_build_1351[dm_build_1352] = byte(dm_build_1353 >> 48)
	dm_build_1352++
	dm_build_1351[dm_build_1352] = byte(dm_build_1353 >> 56)
	return 3
}

func (Dm_build_1355 *dm_build_1297) Dm_build_1354(dm_build_1356 []byte, dm_build_1357 int, dm_build_1358 []byte, dm_build_1359 int, dm_build_1360 int) int {
	copy(dm_build_1356[dm_build_1357:dm_build_1357+dm_build_1360], dm_build_1358[dm_build_1359:dm_build_1359+dm_build_1360])
	return dm_build_1360
}

func (Dm_build_1362 *dm_build_1297) Dm_build_1361(dm_build_1363 []byte, dm_build_1364 int, dm_build_1365 []byte, dm_build_1366 int, dm_build_1367 int) int {
	dm_build_1364 += Dm_build_1362.Dm_build_1344(dm_build_1363, dm_build_1364, uint32(dm_build_1367))
	return 4 + Dm_build_1362.Dm_build_1354(dm_build_1363, dm_build_1364, dm_build_1365, dm_build_1366, dm_build_1367)
}

func (Dm_build_1369 *dm_build_1297) Dm_build_1368(dm_build_1370 []byte, dm_build_1371 int, dm_build_1372 []byte, dm_build_1373 int, dm_build_1374 int) int {
	dm_build_1371 += Dm_build_1369.Dm_build_1339(dm_build_1370, dm_build_1371, uint16(dm_build_1374))
	return 2 + Dm_build_1369.Dm_build_1354(dm_build_1370, dm_build_1371, dm_build_1372, dm_build_1373, dm_build_1374)
}

func (Dm_build_1376 *dm_build_1297) Dm_build_1375(dm_build_1377 []byte, dm_build_1378 int, dm_build_1379 string, dm_build_1380 string, dm_build_1381 *DmConnection) int {
	dm_build_1382 := Dm_build_1376.Dm_build_1511(dm_build_1379, dm_build_1380, dm_build_1381)
	dm_build_1378 += Dm_build_1376.Dm_build_1344(dm_build_1377, dm_build_1378, uint32(len(dm_build_1382)))
	return 4 + Dm_build_1376.Dm_build_1354(dm_build_1377, dm_build_1378, dm_build_1382, 0, len(dm_build_1382))
}

func (Dm_build_1384 *dm_build_1297) Dm_build_1383(dm_build_1385 []byte, dm_build_1386 int, dm_build_1387 string, dm_build_1388 string, dm_build_1389 *DmConnection) int {
	dm_build_1390 := Dm_build_1384.Dm_build_1511(dm_build_1387, dm_build_1388, dm_build_1389)

	dm_build_1386 += Dm_build_1384.Dm_build_1339(dm_build_1385, dm_build_1386, uint16(len(dm_build_1390)))
	return 2 + Dm_build_1384.Dm_build_1354(dm_build_1385, dm_build_1386, dm_build_1390, 0, len(dm_build_1390))
}

func (Dm_build_1392 *dm_build_1297) Dm_build_1391(dm_build_1393 []byte, dm_build_1394 int) byte {
	return dm_build_1393[dm_build_1394]
}

func (Dm_build_1396 *dm_build_1297) Dm_build_1395(dm_build_1397 []byte, dm_build_1398 int) int16 {
	var dm_build_1399 int16
	dm_build_1399 = int16(dm_build_1397[dm_build_1398] & 0xff)
	dm_build_1398++
	dm_build_1399 |= int16(dm_build_1397[dm_build_1398]&0xff) << 8
	return dm_build_1399
}

func (Dm_build_1401 *dm_build_1297) Dm_build_1400(dm_build_1402 []byte, dm_build_1403 int) int32 {
	var dm_build_1404 int32
	dm_build_1404 = int32(dm_build_1402[dm_build_1403] & 0xff)
	dm_build_1403++
	dm_build_1404 |= int32(dm_build_1402[dm_build_1403]&0xff) << 8
	dm_build_1403++
	dm_build_1404 |= int32(dm_build_1402[dm_build_1403]&0xff) << 16
	dm_build_1403++
	dm_build_1404 |= int32(dm_build_1402[dm_build_1403]&0xff) << 24
	return dm_build_1404
}

func (Dm_build_1406 *dm_build_1297) Dm_build_1405(dm_build_1407 []byte, dm_build_1408 int) int64 {
	var dm_build_1409 int64
	dm_build_1409 = int64(dm_build_1407[dm_build_1408] & 0xff)
	dm_build_1408++
	dm_build_1409 |= int64(dm_build_1407[dm_build_1408]&0xff) << 8
	dm_build_1408++
	dm_build_1409 |= int64(dm_build_1407[dm_build_1408]&0xff) << 16
	dm_build_1408++
	dm_build_1409 |= int64(dm_build_1407[dm_build_1408]&0xff) << 24
	dm_build_1408++
	dm_build_1409 |= int64(dm_build_1407[dm_build_1408]&0xff) << 32
	dm_build_1408++
	dm_build_1409 |= int64(dm_build_1407[dm_build_1408]&0xff) << 40
	dm_build_1408++
	dm_build_1409 |= int64(dm_build_1407[dm_build_1408]&0xff) << 48
	dm_build_1408++
	dm_build_1409 |= int64(dm_build_1407[dm_build_1408]&0xff) << 56
	return dm_build_1409
}

func (Dm_build_1411 *dm_build_1297) Dm_build_1410(dm_build_1412 []byte, dm_build_1413 int) float32 {
	return math.Float32frombits(Dm_build_1411.Dm_build_1427(dm_build_1412, dm_build_1413))
}

func (Dm_build_1415 *dm_build_1297) Dm_build_1414(dm_build_1416 []byte, dm_build_1417 int) float64 {
	return math.Float64frombits(Dm_build_1415.Dm_build_1432(dm_build_1416, dm_build_1417))
}

func (Dm_build_1419 *dm_build_1297) Dm_build_1418(dm_build_1420 []byte, dm_build_1421 int) uint8 {
	return uint8(dm_build_1420[dm_build_1421] & 0xff)
}

func (Dm_build_1423 *dm_build_1297) Dm_build_1422(dm_build_1424 []byte, dm_build_1425 int) uint16 {
	var dm_build_1426 uint16
	dm_build_1426 = uint16(dm_build_1424[dm_build_1425] & 0xff)
	dm_build_1425++
	dm_build_1426 |= uint16(dm_build_1424[dm_build_1425]&0xff) << 8
	return dm_build_1426
}

func (Dm_build_1428 *dm_build_1297) Dm_build_1427(dm_build_1429 []byte, dm_build_1430 int) uint32 {
	var dm_build_1431 uint32
	dm_build_1431 = uint32(dm_build_1429[dm_build_1430] & 0xff)
	dm_build_1430++
	dm_build_1431 |= uint32(dm_build_1429[dm_build_1430]&0xff) << 8
	dm_build_1430++
	dm_build_1431 |= uint32(dm_build_1429[dm_build_1430]&0xff) << 16
	dm_build_1430++
	dm_build_1431 |= uint32(dm_build_1429[dm_build_1430]&0xff) << 24
	return dm_build_1431
}

func (Dm_build_1433 *dm_build_1297) Dm_build_1432(dm_build_1434 []byte, dm_build_1435 int) uint64 {
	var dm_build_1436 uint64
	dm_build_1436 = uint64(dm_build_1434[dm_build_1435] & 0xff)
	dm_build_1435++
	dm_build_1436 |= uint64(dm_build_1434[dm_build_1435]&0xff) << 8
	dm_build_1435++
	dm_build_1436 |= uint64(dm_build_1434[dm_build_1435]&0xff) << 16
	dm_build_1435++
	dm_build_1436 |= uint64(dm_build_1434[dm_build_1435]&0xff) << 24
	dm_build_1435++
	dm_build_1436 |= uint64(dm_build_1434[dm_build_1435]&0xff) << 32
	dm_build_1435++
	dm_build_1436 |= uint64(dm_build_1434[dm_build_1435]&0xff) << 40
	dm_build_1435++
	dm_build_1436 |= uint64(dm_build_1434[dm_build_1435]&0xff) << 48
	dm_build_1435++
	dm_build_1436 |= uint64(dm_build_1434[dm_build_1435]&0xff) << 56
	return dm_build_1436
}

func (Dm_build_1438 *dm_build_1297) Dm_build_1437(dm_build_1439 []byte, dm_build_1440 int) []byte {
	dm_build_1441 := Dm_build_1438.Dm_build_1427(dm_build_1439, dm_build_1440)

	dm_build_1442 := make([]byte, dm_build_1441)
	copy(dm_build_1442[:int(dm_build_1441)], dm_build_1439[dm_build_1440+4:dm_build_1440+4+int(dm_build_1441)])
	return dm_build_1442
}

func (Dm_build_1444 *dm_build_1297) Dm_build_1443(dm_build_1445 []byte, dm_build_1446 int) []byte {
	dm_build_1447 := Dm_build_1444.Dm_build_1422(dm_build_1445, dm_build_1446)

	dm_build_1448 := make([]byte, dm_build_1447)
	copy(dm_build_1448[:int(dm_build_1447)], dm_build_1445[dm_build_1446+2:dm_build_1446+2+int(dm_build_1447)])
	return dm_build_1448
}

func (Dm_build_1450 *dm_build_1297) Dm_build_1449(dm_build_1451 []byte, dm_build_1452 int, dm_build_1453 int) []byte {

	dm_build_1454 := make([]byte, dm_build_1453)
	copy(dm_build_1454[:dm_build_1453], dm_build_1451[dm_build_1452:dm_build_1452+dm_build_1453])
	return dm_build_1454
}

func (Dm_build_1456 *dm_build_1297) Dm_build_1455(dm_build_1457 []byte, dm_build_1458 int, dm_build_1459 int, dm_build_1460 string, dm_build_1461 *DmConnection) string {
	return Dm_build_1456.Dm_build_1548(dm_build_1457[dm_build_1458:dm_build_1458+dm_build_1459], dm_build_1460, dm_build_1461)
}

func (Dm_build_1463 *dm_build_1297) Dm_build_1462(dm_build_1464 []byte, dm_build_1465 int, dm_build_1466 string, dm_build_1467 *DmConnection) string {
	dm_build_1468 := Dm_build_1463.Dm_build_1427(dm_build_1464, dm_build_1465)
	dm_build_1465 += 4
	return Dm_build_1463.Dm_build_1455(dm_build_1464, dm_build_1465, int(dm_build_1468), dm_build_1466, dm_build_1467)
}

func (Dm_build_1470 *dm_build_1297) Dm_build_1469(dm_build_1471 []byte, dm_build_1472 int, dm_build_1473 string, dm_build_1474 *DmConnection) string {
	dm_build_1475 := Dm_build_1470.Dm_build_1422(dm_build_1471, dm_build_1472)
	dm_build_1472 += 2
	return Dm_build_1470.Dm_build_1455(dm_build_1471, dm_build_1472, int(dm_build_1475), dm_build_1473, dm_build_1474)
}

func (Dm_build_1477 *dm_build_1297) Dm_build_1476(dm_build_1478 byte) []byte {
	return []byte{dm_build_1478}
}

func (Dm_build_1480 *dm_build_1297) Dm_build_1479(dm_build_1481 int16) []byte {
	return []byte{byte(dm_build_1481), byte(dm_build_1481 >> 8)}
}

func (Dm_build_1483 *dm_build_1297) Dm_build_1482(dm_build_1484 int32) []byte {
	return []byte{byte(dm_build_1484), byte(dm_build_1484 >> 8), byte(dm_build_1484 >> 16), byte(dm_build_1484 >> 24)}
}

func (Dm_build_1486 *dm_build_1297) Dm_build_1485(dm_build_1487 int64) []byte {
	return []byte{byte(dm_build_1487), byte(dm_build_1487 >> 8), byte(dm_build_1487 >> 16), byte(dm_build_1487 >> 24), byte(dm_build_1487 >> 32),
		byte(dm_build_1487 >> 40), byte(dm_build_1487 >> 48), byte(dm_build_1487 >> 56)}
}

func (Dm_build_1489 *dm_build_1297) Dm_build_1488(dm_build_1490 float32) []byte {
	return Dm_build_1489.Dm_build_1500(math.Float32bits(dm_build_1490))
}

func (Dm_build_1492 *dm_build_1297) Dm_build_1491(dm_build_1493 float64) []byte {
	return Dm_build_1492.Dm_build_1503(math.Float64bits(dm_build_1493))
}

func (Dm_build_1495 *dm_build_1297) Dm_build_1494(dm_build_1496 uint8) []byte {
	return []byte{byte(dm_build_1496)}
}

func (Dm_build_1498 *dm_build_1297) Dm_build_1497(dm_build_1499 uint16) []byte {
	return []byte{byte(dm_build_1499), byte(dm_build_1499 >> 8)}
}

func (Dm_build_1501 *dm_build_1297) Dm_build_1500(dm_build_1502 uint32) []byte {
	return []byte{byte(dm_build_1502), byte(dm_build_1502 >> 8), byte(dm_build_1502 >> 16), byte(dm_build_1502 >> 24)}
}

func (Dm_build_1504 *dm_build_1297) Dm_build_1503(dm_build_1505 uint64) []byte {
	return []byte{byte(dm_build_1505), byte(dm_build_1505 >> 8), byte(dm_build_1505 >> 16), byte(dm_build_1505 >> 24), byte(dm_build_1505 >> 32), byte(dm_build_1505 >> 40), byte(dm_build_1505 >> 48), byte(dm_build_1505 >> 56)}
}

func (Dm_build_1507 *dm_build_1297) Dm_build_1506(dm_build_1508 []byte, dm_build_1509 string, dm_build_1510 *DmConnection) []byte {
	if dm_build_1509 == "UTF-8" {
		return dm_build_1508
	}

	if dm_build_1510 == nil {
		if e := dm_build_1553(dm_build_1509); e != nil {
			tmp, err := ioutil.ReadAll(
				transform.NewReader(bytes.NewReader(dm_build_1508), e.NewEncoder()),
			)
			if err != nil {
				panic("UTF8 To Charset error!")
			}

			return tmp
		}

		panic("Unsupported Charset!")
	}

	if dm_build_1510.encodeBuffer == nil {
		dm_build_1510.encodeBuffer = bytes.NewBuffer(nil)
		dm_build_1510.encode = dm_build_1553(dm_build_1510.getServerEncoding())
		dm_build_1510.transformReaderDst = make([]byte, 4096)
		dm_build_1510.transformReaderSrc = make([]byte, 4096)
	}

	if e := dm_build_1510.encode; e != nil {

		dm_build_1510.encodeBuffer.Reset()

		n, err := dm_build_1510.encodeBuffer.ReadFrom(
			Dm_build_1567(bytes.NewReader(dm_build_1508), e.NewEncoder(), dm_build_1510.transformReaderDst, dm_build_1510.transformReaderSrc),
		)
		if err != nil {
			panic("UTF8 To Charset error!")
		}
		var tmp = make([]byte, n)
		if _, err = dm_build_1510.encodeBuffer.Read(tmp); err != nil {
			panic("UTF8 To Charset error!")
		}
		return tmp
	}

	panic("Unsupported Charset!")
}

func (Dm_build_1512 *dm_build_1297) Dm_build_1511(dm_build_1513 string, dm_build_1514 string, dm_build_1515 *DmConnection) []byte {
	return Dm_build_1512.Dm_build_1506([]byte(dm_build_1513), dm_build_1514, dm_build_1515)
}

func (Dm_build_1517 *dm_build_1297) Dm_build_1516(dm_build_1518 []byte) byte {
	return Dm_build_1517.Dm_build_1391(dm_build_1518, 0)
}

func (Dm_build_1520 *dm_build_1297) Dm_build_1519(dm_build_1521 []byte) int16 {
	return Dm_build_1520.Dm_build_1395(dm_build_1521, 0)
}

func (Dm_build_1523 *dm_build_1297) Dm_build_1522(dm_build_1524 []byte) int32 {
	return Dm_build_1523.Dm_build_1400(dm_build_1524, 0)
}

func (Dm_build_1526 *dm_build_1297) Dm_build_1525(dm_build_1527 []byte) int64 {
	return Dm_build_1526.Dm_build_1405(dm_build_1527, 0)
}

func (Dm_build_1529 *dm_build_1297) Dm_build_1528(dm_build_1530 []byte) float32 {
	return Dm_build_1529.Dm_build_1410(dm_build_1530, 0)
}

func (Dm_build_1532 *dm_build_1297) Dm_build_1531(dm_build_1533 []byte) float64 {
	return Dm_build_1532.Dm_build_1414(dm_build_1533, 0)
}

func (Dm_build_1535 *dm_build_1297) Dm_build_1534(dm_build_1536 []byte) uint8 {
	return Dm_build_1535.Dm_build_1418(dm_build_1536, 0)
}

func (Dm_build_1538 *dm_build_1297) Dm_build_1537(dm_build_1539 []byte) uint16 {
	return Dm_build_1538.Dm_build_1422(dm_build_1539, 0)
}

func (Dm_build_1541 *dm_build_1297) Dm_build_1540(dm_build_1542 []byte) uint32 {
	return Dm_build_1541.Dm_build_1427(dm_build_1542, 0)
}

func (Dm_build_1544 *dm_build_1297) Dm_build_1543(dm_build_1545 []byte, dm_build_1546 string, dm_build_1547 *DmConnection) []byte {
	if dm_build_1546 == "UTF-8" {
		return dm_build_1545
	}

	if dm_build_1547 == nil {
		if e := dm_build_1553(dm_build_1546); e != nil {

			tmp, err := ioutil.ReadAll(
				transform.NewReader(bytes.NewReader(dm_build_1545), e.NewDecoder()),
			)
			if err != nil {

				panic("Charset To UTF8 error!")
			}

			return tmp
		}

		panic("Unsupported Charset!")
	}

	if dm_build_1547.encodeBuffer == nil {
		dm_build_1547.encodeBuffer = bytes.NewBuffer(nil)
		dm_build_1547.encode = dm_build_1553(dm_build_1547.getServerEncoding())
		dm_build_1547.transformReaderDst = make([]byte, 4096)
		dm_build_1547.transformReaderSrc = make([]byte, 4096)
	}

	if e := dm_build_1547.encode; e != nil {

		dm_build_1547.encodeBuffer.Reset()

		n, err := dm_build_1547.encodeBuffer.ReadFrom(
			Dm_build_1567(bytes.NewReader(dm_build_1545), e.NewDecoder(), dm_build_1547.transformReaderDst, dm_build_1547.transformReaderSrc),
		)
		if err != nil {

			panic("Charset To UTF8 error!")
		}

		return dm_build_1547.encodeBuffer.Next(int(n))
	}

	panic("Unsupported Charset!")
}

func (Dm_build_1549 *dm_build_1297) Dm_build_1548(dm_build_1550 []byte, dm_build_1551 string, dm_build_1552 *DmConnection) string {
	return string(Dm_build_1549.Dm_build_1543(dm_build_1550, dm_build_1551, dm_build_1552))
}

func dm_build_1553(dm_build_1554 string) encoding.Encoding {
	if e, err := ianaindex.MIB.Encoding(dm_build_1554); err == nil && e != nil {
		return e
	}
	return nil
}

type Dm_build_1555 struct {
	dm_build_1556 io.Reader
	dm_build_1557 transform.Transformer
	dm_build_1558 error

	dm_build_1559                []byte
	dm_build_1560, dm_build_1561 int

	dm_build_1562                []byte
	dm_build_1563, dm_build_1564 int

	dm_build_1565 bool
}

const dm_build_1566 = 4096

func Dm_build_1567(dm_build_1568 io.Reader, dm_build_1569 transform.Transformer, dm_build_1570 []byte, dm_build_1571 []byte) *Dm_build_1555 {
	dm_build_1569.Reset()
	return &Dm_build_1555{
		dm_build_1556: dm_build_1568,
		dm_build_1557: dm_build_1569,
		dm_build_1559: dm_build_1570,
		dm_build_1562: dm_build_1571,
	}
}

func (dm_build_1573 *Dm_build_1555) Read(dm_build_1574 []byte) (int, error) {
	dm_build_1575, dm_build_1576 := 0, error(nil)
	for {

		if dm_build_1573.dm_build_1560 != dm_build_1573.dm_build_1561 {
			dm_build_1575 = copy(dm_build_1574, dm_build_1573.dm_build_1559[dm_build_1573.dm_build_1560:dm_build_1573.dm_build_1561])
			dm_build_1573.dm_build_1560 += dm_build_1575
			if dm_build_1573.dm_build_1560 == dm_build_1573.dm_build_1561 && dm_build_1573.dm_build_1565 {
				return dm_build_1575, dm_build_1573.dm_build_1558
			}
			return dm_build_1575, nil
		} else if dm_build_1573.dm_build_1565 {
			return 0, dm_build_1573.dm_build_1558
		}

		if dm_build_1573.dm_build_1563 != dm_build_1573.dm_build_1564 || dm_build_1573.dm_build_1558 != nil {
			dm_build_1573.dm_build_1560 = 0
			dm_build_1573.dm_build_1561, dm_build_1575, dm_build_1576 = dm_build_1573.dm_build_1557.Transform(dm_build_1573.dm_build_1559, dm_build_1573.dm_build_1562[dm_build_1573.dm_build_1563:dm_build_1573.dm_build_1564], dm_build_1573.dm_build_1558 == io.EOF)
			dm_build_1573.dm_build_1563 += dm_build_1575

			switch {
			case dm_build_1576 == nil:
				if dm_build_1573.dm_build_1563 != dm_build_1573.dm_build_1564 {
					dm_build_1573.dm_build_1558 = nil
				}

				dm_build_1573.dm_build_1565 = dm_build_1573.dm_build_1558 != nil
				continue
			case dm_build_1576 == transform.ErrShortDst && (dm_build_1573.dm_build_1561 != 0 || dm_build_1575 != 0):

				continue
			case dm_build_1576 == transform.ErrShortSrc && dm_build_1573.dm_build_1564-dm_build_1573.dm_build_1563 != len(dm_build_1573.dm_build_1562) && dm_build_1573.dm_build_1558 == nil:

			default:
				dm_build_1573.dm_build_1565 = true

				if dm_build_1573.dm_build_1558 == nil || dm_build_1573.dm_build_1558 == io.EOF {
					dm_build_1573.dm_build_1558 = dm_build_1576
				}
				continue
			}
		}

		if dm_build_1573.dm_build_1563 != 0 {
			dm_build_1573.dm_build_1563, dm_build_1573.dm_build_1564 = 0, copy(dm_build_1573.dm_build_1562, dm_build_1573.dm_build_1562[dm_build_1573.dm_build_1563:dm_build_1573.dm_build_1564])
		}
		dm_build_1575, dm_build_1573.dm_build_1558 = dm_build_1573.dm_build_1556.Read(dm_build_1573.dm_build_1562[dm_build_1573.dm_build_1564:])
		dm_build_1573.dm_build_1564 += dm_build_1575
	}
}
