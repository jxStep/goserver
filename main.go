package main

import "fmt"

func main() {

	var whatToSay string
	var i int

	whatToSay = "Goodbye, cruel world!"

	fmt.Println(whatToSay)

	i = 7

	fmt.Println("i is set to", i)
	whatWasSaid, otherThingThatWasSaid := saySomething()

	fmt.Println("The function returned", whatWasSaid, otherThingThatWasSaid)
}

func saySomething() (string, string) {
	return "Something", "else"
}