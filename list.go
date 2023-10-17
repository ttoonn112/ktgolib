package ktgolib

import (
	"sort"
)

func ListSort(slist []map[string]interface{}, key string) []map[string]interface{} {
	keys := []string{}
	resultset := map[string]map[string]interface{}{}
	for _, currMap := range slist {
		text := T(currMap,key)
		resultset[text] = currMap
		keys = append(keys, text)
	}
	sort.Strings(keys)
	list := []map[string]interface{}{}
	for _, kv := range keys {
		if rdata, ok := resultset[kv]; ok {
			list = append(list, rdata)
		}
	}
	return list
}

func ListSortDesc(slist []map[string]interface{}, key string) []map[string]interface{} {
	compare := func(i, j int) bool {
		key1 := slist[i][key].(string)
		key2 := slist[j][key].(string)
		return key1 > key2 // Change to < for ascending order
	}

	sort.SliceStable(slist, func(i, j int) bool {
		return compare(i, j)
	})

	return slist
}
