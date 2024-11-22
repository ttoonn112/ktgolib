package ktgolib

import (
	"sort"
)

// sortType asc or desc
func ListSort(list []map[string]interface{}, key string, sortType string) []map[string]interface{} {
	if sortType == "desc" {
		sort.Slice(list, func(i, j int) bool {
			return T(list[i],key) > T(list[j],key)
		})
	}else{
		sort.Slice(list, func(i, j int) bool {
			return T(list[i],key) < T(list[j],key)
		})
	}
	return list
}

func ListSortNumber(list []map[string]interface{}, key string, sortType string) []map[string]interface{} {
	if sortType == "desc" {
		sort.Slice(list, func(i, j int) bool {
			return SF64(list[i],key) > SF64(list[j],key)
		})
	}else{
		sort.Slice(list, func(i, j int) bool {
			return SF64(list[i],key) < SF64(list[j],key)
		})
	}
	return list
}

func ListReverse(arr []map[string]interface{}) []map[string]interface{} {
    reversed := make([]map[string]interface{}, len(arr))
    j := 0
    for i := len(arr) - 1; i >= 0; i-- {
        reversed[j] = arr[i]
        j++
    }
    return reversed
}
