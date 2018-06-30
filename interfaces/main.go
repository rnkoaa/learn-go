package main

import "fmt"

type User struct {
	firstName, lastName string
}

// Namer
type Namer interface {
	Name() string
}

func (u *User) Name() string {
	return fmt.Sprintf("%s %s", u.firstName, u.lastName)
}

func greet(n Namer) string {
	mName := n.Name()
	return fmt.Sprintf("Dear %s", mName)
}

func main() {
	nm := &User{
		firstName: "Richard",
		lastName:  "Amoako Agyei",
	}

	fmt.Println(greet(nm))
}
