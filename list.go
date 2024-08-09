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

func ListReverse(arr []map[string]interface{}) []map[string]interface{} {
    reversed := make([]map[string]interface{}, len(arr))
    j := 0
    for i := len(arr) - 1; i >= 0; i-- {
        reversed[j] = arr[i]
        j++
    }
    return reversed
}

func MapSortByKey(olist map[string]map[string]interface{}) map[string]map[string]interface{} {
	keys := []string{}
	for key, _ := range olist {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	resultset := map[string]map[string]interface{}{}
	for _, kv := range keys {
		if rdata, ok := olist[kv]; ok {
			resultset[kv] = rdata
		}
	}
	return resultset
}
