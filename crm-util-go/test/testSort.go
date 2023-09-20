package test

import (
	"fmt"
	"sort"
)

type employee struct {
	name   string
	salary int
}

func sortEmployeeSalary(e []employee, descending bool) {
	sort.SliceStable(e, func(i, j int) bool {
		if descending {
			return e[i].salary > e[j].salary
		} else {
			return e[i].salary < e[j].salary
		}
	})
}

func sortEmployeeName(e []employee, descending bool) {
	sort.SliceStable(e, func(i, j int) bool {
		if descending {
			return e[i].name > e[j].name
		} else {
			return e[i].name < e[j].name
		}
	})
}

func sortEmployeeSalaryName(e []employee, descending bool) {
	sort.SliceStable(e, func(i, j int) bool {
		var sortSalary bool

		if descending {
			sortSalary = e[i].salary > e[j].salary
		} else {
			sortSalary = e[i].salary < e[j].salary
		}

		if e[i].salary == e[j].salary {
			var sortName bool
			if descending {
				sortName = e[i].name > e[j].name
			} else {
				sortName = e[i].name < e[j].name
			}
			return sortName
		}

		return sortSalary
	})
}

func TestSortStruct() {
	empList := []employee{
		employee{name: "X", salary: 100},
		employee{name: "John", salary: 3000},
		employee{name: "C", salary: 100},
		employee{name: "Bill", salary: 4000},
		employee{name: "Pui", salary: 200},
		employee{name: "Sam", salary: 1000},
		employee{name: "A", salary: 100},
		employee{name: "B", salary: 100},
		employee{name: "Z", salary: 100},
	}

	//sortEmployeeSalary(empList, false)
	//sortEmployeeName(empList, false)
	sortEmployeeSalaryName(empList, true)

	for _, employee := range empList {
		fmt.Printf("Name: %s Salary %d\n", employee.name, employee.salary)
	}
}
