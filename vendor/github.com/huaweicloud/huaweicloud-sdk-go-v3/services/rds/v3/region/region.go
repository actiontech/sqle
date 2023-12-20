package region

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
)

var (
	AF_SOUTH_1 = region.NewRegion("af-south-1",
		"https://rds.af-south-1.myhuaweicloud.com")
	CN_NORTH_4 = region.NewRegion("cn-north-4",
		"https://rds.cn-north-4.myhuaweicloud.com")
	CN_NORTH_1 = region.NewRegion("cn-north-1",
		"https://rds.cn-north-1.myhuaweicloud.com")
	CN_EAST_2 = region.NewRegion("cn-east-2",
		"https://rds.cn-east-2.myhuaweicloud.com")
	CN_EAST_3 = region.NewRegion("cn-east-3",
		"https://rds.cn-east-3.myhuaweicloud.com")
	CN_SOUTH_1 = region.NewRegion("cn-south-1",
		"https://rds.cn-south-1.myhuaweicloud.com")
	CN_SOUTHWEST_2 = region.NewRegion("cn-southwest-2",
		"https://rds.cn-southwest-2.myhuaweicloud.com")
	AP_SOUTHEAST_2 = region.NewRegion("ap-southeast-2",
		"https://rds.ap-southeast-2.myhuaweicloud.com")
	AP_SOUTHEAST_1 = region.NewRegion("ap-southeast-1",
		"https://rds.ap-southeast-1.myhuaweicloud.com")
	AP_SOUTHEAST_3 = region.NewRegion("ap-southeast-3",
		"https://rds.ap-southeast-3.myhuaweicloud.com")
	RU_NORTHWEST_2 = region.NewRegion("ru-northwest-2",
		"https://rds.ru-northwest-2.myhuaweicloud.com")
	SA_BRAZIL_1 = region.NewRegion("sa-brazil-1",
		"https://rds.sa-brazil-1.myhuaweicloud.com")
	LA_NORTH_2 = region.NewRegion("la-north-2",
		"https://rds.la-north-2.myhuaweicloud.com")
	CN_SOUTH_2 = region.NewRegion("cn-south-2",
		"https://rds.cn-south-2.myhuaweicloud.com")
	NA_MEXICO_1 = region.NewRegion("na-mexico-1",
		"https://rds.na-mexico-1.myhuaweicloud.com")
	LA_SOUTH_2 = region.NewRegion("la-south-2",
		"https://rds.la-south-2.myhuaweicloud.com")
	CN_NORTH_9 = region.NewRegion("cn-north-9",
		"https://rds.cn-north-9.myhuaweicloud.com")
	CN_NORTH_2 = region.NewRegion("cn-north-2",
		"https://rds.cn-north-2.myhuaweicloud.com")
	TR_WEST_1 = region.NewRegion("tr-west-1",
		"https://rds.tr-west-1.myhuaweicloud.com")
	AP_SOUTHEAST_4 = region.NewRegion("ap-southeast-4",
		"https://rds.ap-southeast-4.myhuaweicloud.com")
	AE_AD_1 = region.NewRegion("ae-ad-1",
		"https://rds.ae-ad-1.myhuaweicloud.com")
	EU_WEST_101 = region.NewRegion("eu-west-101",
		"https://rds.eu-west-101.myhuaweicloud.eu")
)

var staticFields = map[string]*region.Region{
	"af-south-1":     AF_SOUTH_1,
	"cn-north-4":     CN_NORTH_4,
	"cn-north-1":     CN_NORTH_1,
	"cn-east-2":      CN_EAST_2,
	"cn-east-3":      CN_EAST_3,
	"cn-south-1":     CN_SOUTH_1,
	"cn-southwest-2": CN_SOUTHWEST_2,
	"ap-southeast-2": AP_SOUTHEAST_2,
	"ap-southeast-1": AP_SOUTHEAST_1,
	"ap-southeast-3": AP_SOUTHEAST_3,
	"ru-northwest-2": RU_NORTHWEST_2,
	"sa-brazil-1":    SA_BRAZIL_1,
	"la-north-2":     LA_NORTH_2,
	"cn-south-2":     CN_SOUTH_2,
	"na-mexico-1":    NA_MEXICO_1,
	"la-south-2":     LA_SOUTH_2,
	"cn-north-9":     CN_NORTH_9,
	"cn-north-2":     CN_NORTH_2,
	"tr-west-1":      TR_WEST_1,
	"ap-southeast-4": AP_SOUTHEAST_4,
	"ae-ad-1":        AE_AD_1,
	"eu-west-101":    EU_WEST_101,
}

var provider = region.DefaultProviderChain("RDS")

func ValueOf(regionId string) *region.Region {
	if regionId == "" {
		panic("unexpected empty parameter: regionId")
	}

	reg := provider.GetRegion(regionId)
	if reg != nil {
		return reg
	}

	if _, ok := staticFields[regionId]; ok {
		return staticFields[regionId]
	}
	panic(fmt.Sprintf("unexpected regionId: %s", regionId))
}
