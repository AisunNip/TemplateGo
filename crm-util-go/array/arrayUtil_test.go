package array_test

import (
	"crm-util-go/array"
	"fmt"
	"testing"
)

type Books struct {
	Title   string `json:"title,omitempty"`
	Author  string `json:"author,omitempty"`
	Subject string `json:"subject,omitempty"`
	BookID  int    `json:"bookID,omitempty"`
}

func TestContainString(t *testing.T) {
	stringArray := []string{"ccc", "aaa", "zzz", "bbb"}
	isContain := array.ContainString("bbb", stringArray)
	fmt.Println("ContainString", isContain)
	if !isContain {
		t.Error("TestContainString Error")
	}
}

func TestSortStrings(t *testing.T) {
	stringArray := []string{"ccc", "aaa", "zzz", "bbb"}
	array.SortStrings(stringArray)
	fmt.Println("SortStrings", stringArray)
}

func TestSortInt(t *testing.T) {
	intArray := []int{10, 5, 25, 351, 14, 9}
	array.SortInt(intArray)
	fmt.Println("SortInt", intArray)
}

func TestSortFloat64(t *testing.T) {
	floatArray := []float64{18787677.878716, 565435.321, 7888.545, 8787677.8716, 987654.252}
	array.SortFloat64(floatArray)
	fmt.Println("SortFloat64", floatArray)
}

func TestAddRemoveArray(t *testing.T) {
	var dataList []interface{}
	dataList = array.Add(dataList, "Thai1", "Thai2")
	dataList = array.Add(dataList, "Thai3")
	dataList = array.Add(dataList, "Thai4")

	dataList = array.Remove(dataList, 2)
	fmt.Println(dataList)
	fmt.Println("#################################")

	var book1 Books
	book1.Title = "Title1"
	book1.Author = "A1"

	var book2 Books
	book2.Title = "Title2"
	book2.Author = "A2"

	var book3 Books
	book3.Title = "Title3"
	book3.Author = "A3"

	var bookList []interface{}
	bookList = array.Add(bookList, book1)
	bookList = array.Add(bookList, book2)
	bookList = array.Add(bookList, book3)
	bookList = array.Add(bookList, nil)
	bookList = array.Add(bookList, nil)
	fmt.Println(bookList)
}