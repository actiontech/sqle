/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */
package dm

import (
	"container/list"
	"io"
)

type Dm_build_0 struct {
	dm_build_1 *list.List
	dm_build_2 *dm_build_54
	dm_build_3 int
}

func Dm_build_4() *Dm_build_0 {
	return &Dm_build_0{
		dm_build_1: list.New(),
		dm_build_3: 0,
	}
}

func (dm_build_6 *Dm_build_0) Dm_build_5() int {
	return dm_build_6.dm_build_3
}

func (dm_build_8 *Dm_build_0) Dm_build_7(dm_build_9 *Dm_build_78, dm_build_10 int) int {
	var dm_build_11 = 0
	var dm_build_12 = 0
	for dm_build_11 < dm_build_10 && dm_build_8.dm_build_2 != nil {
		dm_build_12 = dm_build_8.dm_build_2.dm_build_62(dm_build_9, dm_build_10-dm_build_11)
		if dm_build_8.dm_build_2.dm_build_57 == 0 {
			dm_build_8.dm_build_44()
		}
		dm_build_11 += dm_build_12
		dm_build_8.dm_build_3 -= dm_build_12
	}
	return dm_build_11
}

func (dm_build_14 *Dm_build_0) Dm_build_13(dm_build_15 []byte, dm_build_16 int, dm_build_17 int) int {
	var dm_build_18 = 0
	var dm_build_19 = 0
	for dm_build_18 < dm_build_17 && dm_build_14.dm_build_2 != nil {
		dm_build_19 = dm_build_14.dm_build_2.dm_build_66(dm_build_15, dm_build_16, dm_build_17-dm_build_18)
		if dm_build_14.dm_build_2.dm_build_57 == 0 {
			dm_build_14.dm_build_44()
		}
		dm_build_18 += dm_build_19
		dm_build_14.dm_build_3 -= dm_build_19
		dm_build_16 += dm_build_19
	}
	return dm_build_18
}

func (dm_build_21 *Dm_build_0) Dm_build_20(dm_build_22 io.Writer, dm_build_23 int) int {
	var dm_build_24 = 0
	var dm_build_25 = 0
	for dm_build_24 < dm_build_23 && dm_build_21.dm_build_2 != nil {
		dm_build_25 = dm_build_21.dm_build_2.dm_build_71(dm_build_22, dm_build_23-dm_build_24)
		if dm_build_21.dm_build_2.dm_build_57 == 0 {
			dm_build_21.dm_build_44()
		}
		dm_build_24 += dm_build_25
		dm_build_21.dm_build_3 -= dm_build_25
	}
	return dm_build_24
}

func (dm_build_27 *Dm_build_0) Dm_build_26(dm_build_28 []byte, dm_build_29 int, dm_build_30 int) {
	if dm_build_30 == 0 {
		return
	}
	var dm_build_31 = dm_build_58(dm_build_28, dm_build_29, dm_build_30)
	if dm_build_27.dm_build_2 == nil {
		dm_build_27.dm_build_2 = dm_build_31
	} else {
		dm_build_27.dm_build_1.PushBack(dm_build_31)
	}
	dm_build_27.dm_build_3 += dm_build_30
}

func (dm_build_33 *Dm_build_0) dm_build_32(dm_build_34 int) byte {
	var dm_build_35 = dm_build_34
	var dm_build_36 = dm_build_33.dm_build_2
	for dm_build_35 > 0 && dm_build_36 != nil {
		if dm_build_36.dm_build_57 == 0 {
			continue
		}
		if dm_build_35 > dm_build_36.dm_build_57-1 {
			dm_build_35 -= dm_build_36.dm_build_57
			dm_build_36 = dm_build_33.dm_build_1.Front().Value.(*dm_build_54)
		} else {
			break
		}
	}
	return dm_build_36.dm_build_75(dm_build_35)
}
func (dm_build_38 *Dm_build_0) Dm_build_37(dm_build_39 *Dm_build_0) {
	if dm_build_39.dm_build_3 == 0 {
		return
	}
	var dm_build_40 = dm_build_39.dm_build_2
	for dm_build_40 != nil {
		dm_build_38.dm_build_41(dm_build_40)
		dm_build_39.dm_build_44()
		dm_build_40 = dm_build_39.dm_build_2
	}
	dm_build_39.dm_build_3 = 0
}
func (dm_build_42 *Dm_build_0) dm_build_41(dm_build_43 *dm_build_54) {
	if dm_build_43.dm_build_57 == 0 {
		return
	}
	if dm_build_42.dm_build_2 == nil {
		dm_build_42.dm_build_2 = dm_build_43
	} else {
		dm_build_42.dm_build_1.PushBack(dm_build_43)
	}
	dm_build_42.dm_build_3 += dm_build_43.dm_build_57
}

func (dm_build_45 *Dm_build_0) dm_build_44() {
	var dm_build_46 = dm_build_45.dm_build_1.Front()
	if dm_build_46 == nil {
		dm_build_45.dm_build_2 = nil
	} else {
		dm_build_45.dm_build_2 = dm_build_46.Value.(*dm_build_54)
		dm_build_45.dm_build_1.Remove(dm_build_46)
	}
}

func (dm_build_48 *Dm_build_0) Dm_build_47() []byte {
	var dm_build_49 = make([]byte, dm_build_48.dm_build_3)
	var dm_build_50 = dm_build_48.dm_build_2
	var dm_build_51 = 0
	var dm_build_52 = len(dm_build_49)
	var dm_build_53 = 0
	for dm_build_50 != nil {
		if dm_build_50.dm_build_57 > 0 {
			if dm_build_52 > dm_build_50.dm_build_57 {
				dm_build_53 = dm_build_50.dm_build_57
			} else {
				dm_build_53 = dm_build_52
			}
			copy(dm_build_49[dm_build_51:dm_build_51+dm_build_53], dm_build_50.dm_build_55[dm_build_50.dm_build_56:dm_build_50.dm_build_56+dm_build_53])
			dm_build_51 += dm_build_53
			dm_build_52 -= dm_build_53
		}
		if dm_build_48.dm_build_1.Front() == nil {
			dm_build_50 = nil
		} else {
			dm_build_50 = dm_build_48.dm_build_1.Front().Value.(*dm_build_54)
		}
	}
	return dm_build_49
}

type dm_build_54 struct {
	dm_build_55 []byte
	dm_build_56 int
	dm_build_57 int
}

func dm_build_58(dm_build_59 []byte, dm_build_60 int, dm_build_61 int) *dm_build_54 {
	return &dm_build_54{
		dm_build_59,
		dm_build_60,
		dm_build_61,
	}
}

func (dm_build_63 *dm_build_54) dm_build_62(dm_build_64 *Dm_build_78, dm_build_65 int) int {
	if dm_build_63.dm_build_57 <= dm_build_65 {
		dm_build_65 = dm_build_63.dm_build_57
	}
	dm_build_64.Dm_build_157(dm_build_63.dm_build_55[dm_build_63.dm_build_56 : dm_build_63.dm_build_56+dm_build_65])
	dm_build_63.dm_build_56 += dm_build_65
	dm_build_63.dm_build_57 -= dm_build_65
	return dm_build_65
}

func (dm_build_67 *dm_build_54) dm_build_66(dm_build_68 []byte, dm_build_69 int, dm_build_70 int) int {
	if dm_build_67.dm_build_57 <= dm_build_70 {
		dm_build_70 = dm_build_67.dm_build_57
	}
	copy(dm_build_68[dm_build_69:dm_build_69+dm_build_70], dm_build_67.dm_build_55[dm_build_67.dm_build_56:dm_build_67.dm_build_56+dm_build_70])
	dm_build_67.dm_build_56 += dm_build_70
	dm_build_67.dm_build_57 -= dm_build_70
	return dm_build_70
}

func (dm_build_72 *dm_build_54) dm_build_71(dm_build_73 io.Writer, dm_build_74 int) int {
	if dm_build_72.dm_build_57 <= dm_build_74 {
		dm_build_74 = dm_build_72.dm_build_57
	}
	dm_build_73.Write(dm_build_72.dm_build_55[dm_build_72.dm_build_56 : dm_build_72.dm_build_56+dm_build_74])
	dm_build_72.dm_build_56 += dm_build_74
	dm_build_72.dm_build_57 -= dm_build_74
	return dm_build_74
}
func (dm_build_76 *dm_build_54) dm_build_75(dm_build_77 int) byte {
	return dm_build_76.dm_build_55[dm_build_76.dm_build_56+dm_build_77]
}
