/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"io"
	"math"
)

type Dm_build_78 struct {
	dm_build_79 []byte
	dm_build_80 int
}

func Dm_build_81(dm_build_82 int) *Dm_build_78 {
	return &Dm_build_78{make([]byte, 0, dm_build_82), 0}
}

func Dm_build_83(dm_build_84 []byte) *Dm_build_78 {
	return &Dm_build_78{dm_build_84, 0}
}

func (dm_build_86 *Dm_build_78) dm_build_85(dm_build_87 int) *Dm_build_78 {

	dm_build_88 := len(dm_build_86.dm_build_79)
	dm_build_89 := cap(dm_build_86.dm_build_79)

	if dm_build_88+dm_build_87 <= dm_build_89 {
		dm_build_86.dm_build_79 = dm_build_86.dm_build_79[:dm_build_88+dm_build_87]
	} else {

		var calCap = int64(math.Max(float64(2*dm_build_89), float64(dm_build_87+dm_build_88)))

		nbuf := make([]byte, dm_build_87+dm_build_88, calCap)
		copy(nbuf, dm_build_86.dm_build_79)
		dm_build_86.dm_build_79 = nbuf
	}

	return dm_build_86
}

func (dm_build_91 *Dm_build_78) Dm_build_90() int {
	return len(dm_build_91.dm_build_79)
}

func (dm_build_93 *Dm_build_78) Dm_build_92(dm_build_94 int) *Dm_build_78 {
	for i := dm_build_94; i < len(dm_build_93.dm_build_79); i++ {
		dm_build_93.dm_build_79[i] = 0
	}
	dm_build_93.dm_build_79 = dm_build_93.dm_build_79[:dm_build_94]
	return dm_build_93
}

func (dm_build_96 *Dm_build_78) Dm_build_95(dm_build_97 int) *Dm_build_78 {
	dm_build_96.dm_build_80 = dm_build_97
	return dm_build_96
}

func (dm_build_99 *Dm_build_78) Dm_build_98() int {
	return dm_build_99.dm_build_80
}

func (dm_build_101 *Dm_build_78) Dm_build_100(dm_build_102 bool) int {
	return len(dm_build_101.dm_build_79) - dm_build_101.dm_build_80
}

func (dm_build_104 *Dm_build_78) Dm_build_103(dm_build_105 int, dm_build_106 bool, dm_build_107 bool) *Dm_build_78 {

	if dm_build_106 {
		if dm_build_107 {
			dm_build_104.dm_build_85(dm_build_105)
		} else {
			dm_build_104.dm_build_79 = dm_build_104.dm_build_79[:len(dm_build_104.dm_build_79)-dm_build_105]
		}
	} else {
		if dm_build_107 {
			dm_build_104.dm_build_80 += dm_build_105
		} else {
			dm_build_104.dm_build_80 -= dm_build_105
		}
	}

	return dm_build_104
}

func (dm_build_109 *Dm_build_78) Dm_build_108(dm_build_110 io.Reader, dm_build_111 int) (int, error) {
	dm_build_112 := len(dm_build_109.dm_build_79)
	dm_build_109.dm_build_85(dm_build_111)
	dm_build_113 := 0
	for dm_build_111 > 0 {
		n, err := dm_build_110.Read(dm_build_109.dm_build_79[dm_build_112+dm_build_113:])
		if n > 0 && err == io.EOF {
			dm_build_113 += n
			dm_build_109.dm_build_79 = dm_build_109.dm_build_79[:dm_build_112+dm_build_113]
			return dm_build_113, nil
		} else if n > 0 && err == nil {
			dm_build_111 -= n
			dm_build_113 += n
		} else if n == 0 && err != nil {
			return -1, ECGO_COMMUNITION_ERROR.addDetailln(err.Error()).throw()
		}
	}

	return dm_build_113, nil
}

func (dm_build_115 *Dm_build_78) Dm_build_114(dm_build_116 io.Writer) (*Dm_build_78, error) {
	if _, err := dm_build_116.Write(dm_build_115.dm_build_79); err != nil {
		return nil, ECGO_COMMUNITION_ERROR.addDetailln(err.Error()).throw()
	}
	return dm_build_115, nil
}

func (dm_build_118 *Dm_build_78) Dm_build_117(dm_build_119 bool) int {
	dm_build_120 := len(dm_build_118.dm_build_79)
	dm_build_118.dm_build_85(1)

	if dm_build_119 {
		return copy(dm_build_118.dm_build_79[dm_build_120:], []byte{1})
	} else {
		return copy(dm_build_118.dm_build_79[dm_build_120:], []byte{0})
	}
}

func (dm_build_122 *Dm_build_78) Dm_build_121(dm_build_123 byte) int {
	dm_build_124 := len(dm_build_122.dm_build_79)
	dm_build_122.dm_build_85(1)

	return copy(dm_build_122.dm_build_79[dm_build_124:], Dm_build_1298.Dm_build_1476(dm_build_123))
}

func (dm_build_126 *Dm_build_78) Dm_build_125(dm_build_127 int16) int {
	dm_build_128 := len(dm_build_126.dm_build_79)
	dm_build_126.dm_build_85(2)

	return copy(dm_build_126.dm_build_79[dm_build_128:], Dm_build_1298.Dm_build_1479(dm_build_127))
}

func (dm_build_130 *Dm_build_78) Dm_build_129(dm_build_131 int32) int {
	dm_build_132 := len(dm_build_130.dm_build_79)
	dm_build_130.dm_build_85(4)

	return copy(dm_build_130.dm_build_79[dm_build_132:], Dm_build_1298.Dm_build_1482(dm_build_131))
}

func (dm_build_134 *Dm_build_78) Dm_build_133(dm_build_135 uint8) int {
	dm_build_136 := len(dm_build_134.dm_build_79)
	dm_build_134.dm_build_85(1)

	return copy(dm_build_134.dm_build_79[dm_build_136:], Dm_build_1298.Dm_build_1494(dm_build_135))
}

func (dm_build_138 *Dm_build_78) Dm_build_137(dm_build_139 uint16) int {
	dm_build_140 := len(dm_build_138.dm_build_79)
	dm_build_138.dm_build_85(2)

	return copy(dm_build_138.dm_build_79[dm_build_140:], Dm_build_1298.Dm_build_1497(dm_build_139))
}

func (dm_build_142 *Dm_build_78) Dm_build_141(dm_build_143 uint32) int {
	dm_build_144 := len(dm_build_142.dm_build_79)
	dm_build_142.dm_build_85(4)

	return copy(dm_build_142.dm_build_79[dm_build_144:], Dm_build_1298.Dm_build_1500(dm_build_143))
}

func (dm_build_146 *Dm_build_78) Dm_build_145(dm_build_147 uint64) int {
	dm_build_148 := len(dm_build_146.dm_build_79)
	dm_build_146.dm_build_85(8)

	return copy(dm_build_146.dm_build_79[dm_build_148:], Dm_build_1298.Dm_build_1503(dm_build_147))
}

func (dm_build_150 *Dm_build_78) Dm_build_149(dm_build_151 float32) int {
	dm_build_152 := len(dm_build_150.dm_build_79)
	dm_build_150.dm_build_85(4)

	return copy(dm_build_150.dm_build_79[dm_build_152:], Dm_build_1298.Dm_build_1500(math.Float32bits(dm_build_151)))
}

func (dm_build_154 *Dm_build_78) Dm_build_153(dm_build_155 float64) int {
	dm_build_156 := len(dm_build_154.dm_build_79)
	dm_build_154.dm_build_85(8)

	return copy(dm_build_154.dm_build_79[dm_build_156:], Dm_build_1298.Dm_build_1503(math.Float64bits(dm_build_155)))
}

func (dm_build_158 *Dm_build_78) Dm_build_157(dm_build_159 []byte) int {
	dm_build_160 := len(dm_build_158.dm_build_79)
	dm_build_158.dm_build_85(len(dm_build_159))
	return copy(dm_build_158.dm_build_79[dm_build_160:], dm_build_159)
}

func (dm_build_162 *Dm_build_78) Dm_build_161(dm_build_163 []byte) int {
	return dm_build_162.Dm_build_129(int32(len(dm_build_163))) + dm_build_162.Dm_build_157(dm_build_163)
}

func (dm_build_165 *Dm_build_78) Dm_build_164(dm_build_166 []byte) int {
	return dm_build_165.Dm_build_133(uint8(len(dm_build_166))) + dm_build_165.Dm_build_157(dm_build_166)
}

func (dm_build_168 *Dm_build_78) Dm_build_167(dm_build_169 []byte) int {
	return dm_build_168.Dm_build_137(uint16(len(dm_build_169))) + dm_build_168.Dm_build_157(dm_build_169)
}

func (dm_build_171 *Dm_build_78) Dm_build_170(dm_build_172 []byte) int {
	return dm_build_171.Dm_build_157(dm_build_172) + dm_build_171.Dm_build_121(0)
}

func (dm_build_174 *Dm_build_78) Dm_build_173(dm_build_175 string, dm_build_176 string, dm_build_177 *DmConnection) int {
	dm_build_178 := Dm_build_1298.Dm_build_1511(dm_build_175, dm_build_176, dm_build_177)
	return dm_build_174.Dm_build_161(dm_build_178)
}

func (dm_build_180 *Dm_build_78) Dm_build_179(dm_build_181 string, dm_build_182 string, dm_build_183 *DmConnection) int {
	dm_build_184 := Dm_build_1298.Dm_build_1511(dm_build_181, dm_build_182, dm_build_183)
	return dm_build_180.Dm_build_164(dm_build_184)
}

func (dm_build_186 *Dm_build_78) Dm_build_185(dm_build_187 string, dm_build_188 string, dm_build_189 *DmConnection) int {
	dm_build_190 := Dm_build_1298.Dm_build_1511(dm_build_187, dm_build_188, dm_build_189)
	return dm_build_186.Dm_build_167(dm_build_190)
}

func (dm_build_192 *Dm_build_78) Dm_build_191(dm_build_193 string, dm_build_194 string, dm_build_195 *DmConnection) int {
	dm_build_196 := Dm_build_1298.Dm_build_1511(dm_build_193, dm_build_194, dm_build_195)
	return dm_build_192.Dm_build_170(dm_build_196)
}

func (dm_build_198 *Dm_build_78) Dm_build_197() byte {
	dm_build_199 := Dm_build_1298.Dm_build_1391(dm_build_198.dm_build_79, dm_build_198.dm_build_80)
	dm_build_198.dm_build_80++
	return dm_build_199
}

func (dm_build_201 *Dm_build_78) Dm_build_200() int16 {
	dm_build_202 := Dm_build_1298.Dm_build_1395(dm_build_201.dm_build_79, dm_build_201.dm_build_80)
	dm_build_201.dm_build_80 += 2
	return dm_build_202
}

func (dm_build_204 *Dm_build_78) Dm_build_203() int32 {
	dm_build_205 := Dm_build_1298.Dm_build_1400(dm_build_204.dm_build_79, dm_build_204.dm_build_80)
	dm_build_204.dm_build_80 += 4
	return dm_build_205
}

func (dm_build_207 *Dm_build_78) Dm_build_206() int64 {
	dm_build_208 := Dm_build_1298.Dm_build_1405(dm_build_207.dm_build_79, dm_build_207.dm_build_80)
	dm_build_207.dm_build_80 += 8
	return dm_build_208
}

func (dm_build_210 *Dm_build_78) Dm_build_209() float32 {
	dm_build_211 := Dm_build_1298.Dm_build_1410(dm_build_210.dm_build_79, dm_build_210.dm_build_80)
	dm_build_210.dm_build_80 += 4
	return dm_build_211
}

func (dm_build_213 *Dm_build_78) Dm_build_212() float64 {
	dm_build_214 := Dm_build_1298.Dm_build_1414(dm_build_213.dm_build_79, dm_build_213.dm_build_80)
	dm_build_213.dm_build_80 += 8
	return dm_build_214
}

func (dm_build_216 *Dm_build_78) Dm_build_215() uint8 {
	dm_build_217 := Dm_build_1298.Dm_build_1418(dm_build_216.dm_build_79, dm_build_216.dm_build_80)
	dm_build_216.dm_build_80 += 1
	return dm_build_217
}

func (dm_build_219 *Dm_build_78) Dm_build_218() uint16 {
	dm_build_220 := Dm_build_1298.Dm_build_1422(dm_build_219.dm_build_79, dm_build_219.dm_build_80)
	dm_build_219.dm_build_80 += 2
	return dm_build_220
}

func (dm_build_222 *Dm_build_78) Dm_build_221() uint32 {
	dm_build_223 := Dm_build_1298.Dm_build_1427(dm_build_222.dm_build_79, dm_build_222.dm_build_80)
	dm_build_222.dm_build_80 += 4
	return dm_build_223
}

func (dm_build_225 *Dm_build_78) Dm_build_224(dm_build_226 int) []byte {
	dm_build_227 := Dm_build_1298.Dm_build_1449(dm_build_225.dm_build_79, dm_build_225.dm_build_80, dm_build_226)
	dm_build_225.dm_build_80 += dm_build_226
	return dm_build_227
}

func (dm_build_229 *Dm_build_78) Dm_build_228() []byte {
	return dm_build_229.Dm_build_224(int(dm_build_229.Dm_build_203()))
}

func (dm_build_231 *Dm_build_78) Dm_build_230() []byte {
	return dm_build_231.Dm_build_224(int(dm_build_231.Dm_build_197()))
}

func (dm_build_233 *Dm_build_78) Dm_build_232() []byte {
	return dm_build_233.Dm_build_224(int(dm_build_233.Dm_build_200()))
}

func (dm_build_235 *Dm_build_78) Dm_build_234(dm_build_236 int) []byte {
	return dm_build_235.Dm_build_224(dm_build_236)
}

func (dm_build_238 *Dm_build_78) Dm_build_237() []byte {
	dm_build_239 := 0
	for dm_build_238.Dm_build_197() != 0 {
		dm_build_239++
	}
	dm_build_238.Dm_build_103(dm_build_239, false, false)
	return dm_build_238.Dm_build_224(dm_build_239)
}

func (dm_build_241 *Dm_build_78) Dm_build_240(dm_build_242 int, dm_build_243 string, dm_build_244 *DmConnection) string {
	return Dm_build_1298.Dm_build_1548(dm_build_241.Dm_build_224(dm_build_242), dm_build_243, dm_build_244)
}

func (dm_build_246 *Dm_build_78) Dm_build_245(dm_build_247 string, dm_build_248 *DmConnection) string {
	return Dm_build_1298.Dm_build_1548(dm_build_246.Dm_build_228(), dm_build_247, dm_build_248)
}

func (dm_build_250 *Dm_build_78) Dm_build_249(dm_build_251 string, dm_build_252 *DmConnection) string {
	return Dm_build_1298.Dm_build_1548(dm_build_250.Dm_build_230(), dm_build_251, dm_build_252)
}

func (dm_build_254 *Dm_build_78) Dm_build_253(dm_build_255 string, dm_build_256 *DmConnection) string {
	return Dm_build_1298.Dm_build_1548(dm_build_254.Dm_build_232(), dm_build_255, dm_build_256)
}

func (dm_build_258 *Dm_build_78) Dm_build_257(dm_build_259 string, dm_build_260 *DmConnection) string {
	return Dm_build_1298.Dm_build_1548(dm_build_258.Dm_build_237(), dm_build_259, dm_build_260)
}

func (dm_build_262 *Dm_build_78) Dm_build_261(dm_build_263 int, dm_build_264 byte) int {
	return dm_build_262.Dm_build_297(dm_build_263, Dm_build_1298.Dm_build_1476(dm_build_264))
}

func (dm_build_266 *Dm_build_78) Dm_build_265(dm_build_267 int, dm_build_268 int16) int {
	return dm_build_266.Dm_build_297(dm_build_267, Dm_build_1298.Dm_build_1479(dm_build_268))
}

func (dm_build_270 *Dm_build_78) Dm_build_269(dm_build_271 int, dm_build_272 int32) int {
	return dm_build_270.Dm_build_297(dm_build_271, Dm_build_1298.Dm_build_1482(dm_build_272))
}

func (dm_build_274 *Dm_build_78) Dm_build_273(dm_build_275 int, dm_build_276 int64) int {
	return dm_build_274.Dm_build_297(dm_build_275, Dm_build_1298.Dm_build_1485(dm_build_276))
}

func (dm_build_278 *Dm_build_78) Dm_build_277(dm_build_279 int, dm_build_280 float32) int {
	return dm_build_278.Dm_build_297(dm_build_279, Dm_build_1298.Dm_build_1488(dm_build_280))
}

func (dm_build_282 *Dm_build_78) Dm_build_281(dm_build_283 int, dm_build_284 float64) int {
	return dm_build_282.Dm_build_297(dm_build_283, Dm_build_1298.Dm_build_1491(dm_build_284))
}

func (dm_build_286 *Dm_build_78) Dm_build_285(dm_build_287 int, dm_build_288 uint8) int {
	return dm_build_286.Dm_build_297(dm_build_287, Dm_build_1298.Dm_build_1494(dm_build_288))
}

func (dm_build_290 *Dm_build_78) Dm_build_289(dm_build_291 int, dm_build_292 uint16) int {
	return dm_build_290.Dm_build_297(dm_build_291, Dm_build_1298.Dm_build_1497(dm_build_292))
}

func (dm_build_294 *Dm_build_78) Dm_build_293(dm_build_295 int, dm_build_296 uint32) int {
	return dm_build_294.Dm_build_297(dm_build_295, Dm_build_1298.Dm_build_1500(dm_build_296))
}

func (dm_build_298 *Dm_build_78) Dm_build_297(dm_build_299 int, dm_build_300 []byte) int {
	return copy(dm_build_298.dm_build_79[dm_build_299:], dm_build_300)
}

func (dm_build_302 *Dm_build_78) Dm_build_301(dm_build_303 int, dm_build_304 []byte) int {
	return dm_build_302.Dm_build_269(dm_build_303, int32(len(dm_build_304))) + dm_build_302.Dm_build_297(dm_build_303+4, dm_build_304)
}

func (dm_build_306 *Dm_build_78) Dm_build_305(dm_build_307 int, dm_build_308 []byte) int {
	return dm_build_306.Dm_build_261(dm_build_307, byte(len(dm_build_308))) + dm_build_306.Dm_build_297(dm_build_307+1, dm_build_308)
}

func (dm_build_310 *Dm_build_78) Dm_build_309(dm_build_311 int, dm_build_312 []byte) int {
	return dm_build_310.Dm_build_265(dm_build_311, int16(len(dm_build_312))) + dm_build_310.Dm_build_297(dm_build_311+2, dm_build_312)
}

func (dm_build_314 *Dm_build_78) Dm_build_313(dm_build_315 int, dm_build_316 []byte) int {
	return dm_build_314.Dm_build_297(dm_build_315, dm_build_316) + dm_build_314.Dm_build_261(dm_build_315+len(dm_build_316), 0)
}

func (dm_build_318 *Dm_build_78) Dm_build_317(dm_build_319 int, dm_build_320 string, dm_build_321 string, dm_build_322 *DmConnection) int {
	return dm_build_318.Dm_build_301(dm_build_319, Dm_build_1298.Dm_build_1511(dm_build_320, dm_build_321, dm_build_322))
}

func (dm_build_324 *Dm_build_78) Dm_build_323(dm_build_325 int, dm_build_326 string, dm_build_327 string, dm_build_328 *DmConnection) int {
	return dm_build_324.Dm_build_305(dm_build_325, Dm_build_1298.Dm_build_1511(dm_build_326, dm_build_327, dm_build_328))
}

func (dm_build_330 *Dm_build_78) Dm_build_329(dm_build_331 int, dm_build_332 string, dm_build_333 string, dm_build_334 *DmConnection) int {
	return dm_build_330.Dm_build_309(dm_build_331, Dm_build_1298.Dm_build_1511(dm_build_332, dm_build_333, dm_build_334))
}

func (dm_build_336 *Dm_build_78) Dm_build_335(dm_build_337 int, dm_build_338 string, dm_build_339 string, dm_build_340 *DmConnection) int {
	return dm_build_336.Dm_build_313(dm_build_337, Dm_build_1298.Dm_build_1511(dm_build_338, dm_build_339, dm_build_340))
}

func (dm_build_342 *Dm_build_78) Dm_build_341(dm_build_343 int) byte {
	return Dm_build_1298.Dm_build_1516(dm_build_342.Dm_build_368(dm_build_343, 1))
}

func (dm_build_345 *Dm_build_78) Dm_build_344(dm_build_346 int) int16 {
	return Dm_build_1298.Dm_build_1519(dm_build_345.Dm_build_368(dm_build_346, 2))
}

func (dm_build_348 *Dm_build_78) Dm_build_347(dm_build_349 int) int32 {
	return Dm_build_1298.Dm_build_1522(dm_build_348.Dm_build_368(dm_build_349, 4))
}

func (dm_build_351 *Dm_build_78) Dm_build_350(dm_build_352 int) int64 {
	return Dm_build_1298.Dm_build_1525(dm_build_351.Dm_build_368(dm_build_352, 8))
}

func (dm_build_354 *Dm_build_78) Dm_build_353(dm_build_355 int) float32 {
	return Dm_build_1298.Dm_build_1528(dm_build_354.Dm_build_368(dm_build_355, 4))
}

func (dm_build_357 *Dm_build_78) Dm_build_356(dm_build_358 int) float64 {
	return Dm_build_1298.Dm_build_1531(dm_build_357.Dm_build_368(dm_build_358, 8))
}

func (dm_build_360 *Dm_build_78) Dm_build_359(dm_build_361 int) uint8 {
	return Dm_build_1298.Dm_build_1534(dm_build_360.Dm_build_368(dm_build_361, 1))
}

func (dm_build_363 *Dm_build_78) Dm_build_362(dm_build_364 int) uint16 {
	return Dm_build_1298.Dm_build_1537(dm_build_363.Dm_build_368(dm_build_364, 2))
}

func (dm_build_366 *Dm_build_78) Dm_build_365(dm_build_367 int) uint32 {
	return Dm_build_1298.Dm_build_1540(dm_build_366.Dm_build_368(dm_build_367, 4))
}

func (dm_build_369 *Dm_build_78) Dm_build_368(dm_build_370 int, dm_build_371 int) []byte {
	return dm_build_369.dm_build_79[dm_build_370 : dm_build_370+dm_build_371]
}

func (dm_build_373 *Dm_build_78) Dm_build_372(dm_build_374 int) []byte {
	dm_build_375 := dm_build_373.Dm_build_347(dm_build_374)
	return dm_build_373.Dm_build_368(dm_build_374+4, int(dm_build_375))
}

func (dm_build_377 *Dm_build_78) Dm_build_376(dm_build_378 int) []byte {
	dm_build_379 := dm_build_377.Dm_build_341(dm_build_378)
	return dm_build_377.Dm_build_368(dm_build_378+1, int(dm_build_379))
}

func (dm_build_381 *Dm_build_78) Dm_build_380(dm_build_382 int) []byte {
	dm_build_383 := dm_build_381.Dm_build_344(dm_build_382)
	return dm_build_381.Dm_build_368(dm_build_382+2, int(dm_build_383))
}

func (dm_build_385 *Dm_build_78) Dm_build_384(dm_build_386 int) []byte {
	dm_build_387 := 0
	for dm_build_385.Dm_build_341(dm_build_386) != 0 {
		dm_build_386++
		dm_build_387++
	}

	return dm_build_385.Dm_build_368(dm_build_386-dm_build_387, int(dm_build_387))
}

func (dm_build_389 *Dm_build_78) Dm_build_388(dm_build_390 int, dm_build_391 string, dm_build_392 *DmConnection) string {
	return Dm_build_1298.Dm_build_1548(dm_build_389.Dm_build_372(dm_build_390), dm_build_391, dm_build_392)
}

func (dm_build_394 *Dm_build_78) Dm_build_393(dm_build_395 int, dm_build_396 string, dm_build_397 *DmConnection) string {
	return Dm_build_1298.Dm_build_1548(dm_build_394.Dm_build_376(dm_build_395), dm_build_396, dm_build_397)
}

func (dm_build_399 *Dm_build_78) Dm_build_398(dm_build_400 int, dm_build_401 string, dm_build_402 *DmConnection) string {
	return Dm_build_1298.Dm_build_1548(dm_build_399.Dm_build_380(dm_build_400), dm_build_401, dm_build_402)
}

func (dm_build_404 *Dm_build_78) Dm_build_403(dm_build_405 int, dm_build_406 string, dm_build_407 *DmConnection) string {
	return Dm_build_1298.Dm_build_1548(dm_build_404.Dm_build_384(dm_build_405), dm_build_406, dm_build_407)
}
