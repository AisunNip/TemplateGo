package array

import "sort"

/* interface{} --> An empty interface may hold values of any type */
func Add(dataList []interface{}, data ...interface{}) []interface{} {
	return append(dataList, data...)
}

func Remove(dataList []interface{}, index int) []interface{} {
	return append(dataList[:index], dataList[index+1:]...)
}

func ContainString(val string, array []string) bool {
	for _, element := range array {
		if element == val {
			return true
		}
	}

	return false
}

func SortStrings(data []string) {
	sort.Strings(data)
}

func SortInt(data []int) {
	sort.Ints(data)
}

func SortFloat64(data []float64) {
	sort.Float64s(data)
}