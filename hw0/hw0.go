// Homework 0: Hello Go!
package main

import "fmt"

func main() {
	// Feel free to use the main function for testing your functions
	fmt.Println("Hello, दुनिया!")

	fmt.Println("Fizzbuzz() test")
	fmt.Printf("Fizzbuzz(%v) = %v\n", 30, Fizzbuzz(30))
	fmt.Printf("Fizzbuzz(%v) = %v\n", 35, Fizzbuzz(35))
	fmt.Printf("Fizzbuzz(%v) = %v\n", 36, Fizzbuzz(36))

	fmt.Println("IsPrime() test")
	fmt.Printf("IsPrime(%v) = %v\n", 1, IsPrime(1))
	fmt.Printf("IsPrime(%v) = %v\n", 2, IsPrime(2))
	fmt.Printf("IsPrime(%v) = %v\n", 22, IsPrime(22))
	fmt.Printf("IsPrime(%v) = %v\n", 32, IsPrime(32))

	fmt.Println("IsPalindrome() test")
	fmt.Printf("IsPalindrome(%v) = %v\n", "123321", IsPalindrome("123321"))
	fmt.Printf("IsPalindrome(%v) = %v\n", "123322", IsPalindrome("123322"))
	fmt.Printf("IsPalindrome(%v) = %v\n", "123221", IsPalindrome("123221"))
}

// Fizzbuzz is a classic introductory programming problem.
// If n is divisible by 3, return "Fizz"
// If n is divisible by 5, return "Buzz"
// If n is divisible by 3 and 5, return "FizzBuzz"
// Otherwise, return the empty string
func Fizzbuzz(n int) string {
	var s string

	if n%3 == 0 {
		s += "Fizz"
	}

	if n%5 == 0 {
		s += "Buzz"
	}

	return s
}

// IsPrime checks if the number is prime.
// You may use any prime algorithm, but you may NOT use the standard library.
func IsPrime(n int) bool {
	if n <= 2 {
		return true
	}

	for i := 2; i*i < n; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

// IsPalindrome checks if the string is a palindrome.
// A palindrome is a string that reads the same backward as forward.
func IsPalindrome(s string) bool {
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}
