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
