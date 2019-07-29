package util

import (
	"fmt"
	"sort"
)

func DeleteSliceItem(src []uint32, target []uint32) []uint32 {
	var location = []int{}
	var result = []uint32{}
	nonDuplicateTarget := RmDuplicate(target)
	for _, deleteValue := range nonDuplicateTarget {
		for index, value := range src {
			if deleteValue == value {
				location = append(location, index)
			}
		}
	}
	sort.Ints(location)
	if 0 == len(location) {
		//说明要删除的元素在src[]中不存在,需要返回原始src[]即可
		result = src
	} else {
		var preValue int
		for index, value := range location {
			result = append(result, src[preValue:value]...)
			if index == len(location)-1 {
				result = append(result, src[value+1:]...)
			}
			preValue = value + 1
		}
	}
	return result
}

func RmDuplicate(list []uint32) []uint32 {
	var x = []uint32{}
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

func DifferenceInt(arr1, arr2 []uint32) []uint32 {
	var diff []uint32
	for _, v1 := range arr1 {
		found := false
		for _, v2 := range arr2 {
			if v1 == v2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, v1)
		}
	}
	return diff
}

func EqualInt(oldIds, newIds []uint32) bool {
	excludedIds, _, addedIds := IntersectionInt(oldIds, newIds)
	if len(excludedIds) == 0 && len(addedIds) == 0 {
		return true
	}
	return false
}

func ToIntSlice(src []interface{}) []int {

	_, ok := src[0].(float64)
	if ok {
		strV0Int := fmt.Sprintf("%d", int(src[0].(float64)))
		strV0 := fmt.Sprintf("%v", src[0])

		if strV0 == strV0Int {
			var iarry []int
			for _, sv := range src {
				iarry = append(iarry, int(sv.(float64)))
			}
			return iarry
		}

		return nil
	}

	_, ok64 := src[0].(int64)
	_, okint := src[0].(int)
	if ok64 {
		var iarry []int
		for _, sv := range src {
			iarry = append(iarry, int(sv.(int64)))
		}
		return iarry
	}

	if okint {
		var iarry []int
		for _, sv := range src {
			iarry = append(iarry, sv.(int))
		}
		return iarry
	}

	return nil
}

func ToStringSlice(src []interface{}) []string {
	_, ok := src[0].(string)
	if ok {
		var sarry []string
		for _, sv := range src {
			sarry = append(sarry, sv.(string))
		}
		return sarry
	}

	return nil
}

/*
	This method allows to compute the intersection of two list.

	input:
		oldIds			:	a list of old ids
		newIds			:	a list of new ids

	return:
		excludedIds		:	the ids which are in old ids but not in new ids
		keptIds			:	the ids which are both in old ids and new ids
		addedIds		:	the ids which are in new ids but not in old olds
*/
func IntersectionInt(oldIds, newIds []uint32) (excludedIds, keptIds, addedIds []uint32) {
	for _, oldId := range oldIds {
		found := false
		for _, newId := range newIds {
			if oldId == newId {
				found = true
				keptIds = append(keptIds, oldId)
				break
			}
		}
		if !found {
			excludedIds = append(excludedIds, oldId)
		}
	}

	for _, newId := range newIds {
		found := false
		for _, keptId := range keptIds {
			if newId == keptId {
				found = true
				break
			}
		}
		if !found {
			addedIds = append(addedIds, newId)
		}
	}
	return excludedIds, keptIds, addedIds
}

func ContainsInt(list []uint32, e uint32) bool {
	for _, a := range list {
		if a == e {
			return true
		}
	}
	return false
}
