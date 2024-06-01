package main

import (
	"math/rand"
	"strings"
	"time"
)

type Person struct {
	Name        string
	Birthday    time.Time
	Description string
	Location    string
	Email       string
	Password    string
	Gender      string
}

func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generateRandomEmail generates a random email address
func generateRandomEmail() string {
	return generateRandomString(10) + "@" + generateRandomString(5) + ".com"
}

// generateRandomPassword generates a random password
func generateRandomPassword() string {
	return generateRandomString(10)
}

// generateRandomDescription generates a random description
func generateRandomDescription(n int) string {
	words := []string{"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit", "sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore", "magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud", "exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea", "commodo", "consequat", "duis", "aute", "irure", "dolor", "in", "reprehenderit", "voluptate", "velit", "esse", "cillum", "dolore", "eu", "fugiat", "nulla", "pariatur", "excepteur", "sint", "occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia", "deserunt", "mollit", "anim", "id", "est", "laborum"}

	var description strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			description.WriteString(" ")
		}
		description.WriteString(words[rand.Intn(len(words))])
	}

	return description.String()
}

// generateRandomLocation generates a random location
func generateRandomLocation() string {
	return generateRandomString(20)
}

// generateRandomGender generates a random gender
func generateRandomGender() string {
	genders := []string{"male", "female"}
	return genders[rand.Intn(len(genders))]
}

// generateRandomPerson generates a random person
func generateRandomPerson() Person {
	return Person{
		Name:        generateRandomString(10),
		Birthday:    time.Now().AddDate(-rand.Intn(50), -rand.Intn(12), -rand.Intn(30)),
		Description: generateRandomDescription(10),
		Location:    generateRandomLocation(),
		Email:       generateRandomEmail(),
		Password:    generateRandomPassword(),
		Gender:      generateRandomGender(),
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	// connect to the database

}
