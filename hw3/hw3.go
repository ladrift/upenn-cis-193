// Homework 3: Interfaces
package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	fmt.Println("========== Problem 1: Sorting Names ==========")
	people := PersonSlice{NewPerson("Ziyi", "Yan"), NewPerson("Xuefei", "Chen"), NewPerson("Zida", "Liu")}
	fmt.Printf("Before: %v\n", people)
	sort.Sort(people)
	fmt.Printf("After: %v\n", people)

	fmt.Println("========== Problem 2: IsPalindrome Redux ==========")
	first := NewPerson("Mr.", "First")
	second := NewPerson("Mr.", "Second")
	them := PersonSlice{first, second, first}
	fmt.Printf("IsPalindrome(%v) = %v\n", them, IsPalindrome(them))

	fmt.Println("========== Problem 3: Functional Programming ==========")
	add := func(x, y int) int { return x + y }
	fmt.Printf("Fold([]int{1, 2, 3, 4, 5}, 0, add) = %v", Fold([]int{1, 2, 3, 4, 5}, 0, add))
	mult := func(x, y int) int { return x * y }
	fmt.Printf("Fold([]int{1, 2, 3, 4, 5}, 1, mult) = %v", Fold([]int{1, 2, 3, 4, 5}, 1, mult))
}

// Problem 1: Sorting Names
// Sorting in Go is done through interfaces!
// To sort a collection (such as a slice), the type must satisfy sort.Interface,
// which requires 3 methods: Len() int, Less(i, j int) bool, and Swap(i, j int).
// To actually sort a slice, you need to first implement all 3 methods on a
// custom type, and then call sort.Sort on your the PersonSlice type.
// See the Go documentation: https://golang.org/pkg/sort/ for full details.

// Person stores a simple profile. These should be sorted by alphabetical order
// by last name, followed by the first name, followed by the ID. You can assume
// the ID will be unique, but the names need not be unique.
// Sorting should be case-sensitive and UTF-8 aware.
type Person struct {
	ID        int
	FirstName string
	LastName  string
}

// NewPerson is a constructor for Person. ID should be assigned automatically in
// sequential order, starting at 1 for the first Person created.
func NewPerson(first, last string) *Person {
	p := new(Person)
	p.ID = nextPersonID
	p.FirstName = first
	p.LastName = last
	nextPersonID++
	return p
}

var nextPersonID int = 1

func (p *Person) String() string {
	return fmt.Sprintf("%s %s(%v)", p.FirstName, p.LastName, p.ID)
}

type PersonSlice []*Person

func (ps PersonSlice) Len() int {
	return len([]*Person(ps))
}

func (ps PersonSlice) Less(i, j int) bool {
	if c := strings.Compare(ps[i].LastName, ps[j].LastName); c != 0 {
		return c < 0
	}
	if c := strings.Compare(ps[i].FirstName, ps[j].FirstName); c != 0 {
		return c < 0
	}
	return ps[i].ID < ps[j].ID
}

func (ps PersonSlice) Swap(i, j int) {
	p := ps[j]
	ps[j] = ps[i]
	ps[i] = p
}

// Problem 2: IsPalindrome Redux
// Using a function that simply requires sort.Interface, you should be able to
// check if a sequence is a palindrome. You may use, adapt, or modify your code
// from HW0. Note that the input does not need to be a string: any type which
// satisfies sort.Interface can (and will) be used to test. This means that the
// only functionality you should use should come from the sort.Interface methods
// Ex: [1, 2, 1] => true

// IsPalindrome checks if the string is a palindrome.
// A palindrome is a string that reads the same backward as forward.
func IsPalindrome(s sort.Interface) bool {
	for i := 0; i < s.Len()/2; i++ {
		if s.Less(i, s.Len()-1-i) || s.Less(s.Len()-1-i, i) {
			return false
		}
	}
	return true
}

// Problem 3: Functional Programming
// Write a function Fold which applies a function repeatedly on a slice,
// producing a single value via repeated application of an input function.
// The behavior of Fold should be as follows:
//   - When s is empty, return v (default value)
//   - When s has 1 value (x0), apply f once: f(v, x0)
//   - When s has 2 values (x0, x1), apply f twice, from left to right: f(f(v, x0), x1)
//   - Continue this pattern recursively to obtain the final result.

// Fold applies a left to right application of f on s starting with v.
// Note the argument signature of f - func(int, int) int.
// This means f is a function which has 2 int arguments and returns an int.
func Fold(s []int, v int, f func(int, int) int) int {
	if len(s) == 0 {
		return v
	}

	return Fold(s[1:], f(v, s[0]), f)
}
