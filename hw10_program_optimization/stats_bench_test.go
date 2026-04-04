package hw10programoptimization

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
)

var result DomainStat

func BenchmarkGetDomainStat(b *testing.B) {
	userData := generateUsersJSON(generateUsers())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, _ = GetDomainStat(bytes.NewBufferString(userData), "com")
	}
}

func generateUsersJSON(u users) string {
	var builder strings.Builder

	for _, v := range u {
		jsonBytes, _ := json.Marshal(v)
		builder.Write(jsonBytes)
		builder.WriteByte('\n')
	}

	return builder.String()
}

func generateUsers() users {
	var u users
	for i := 0; i < 100_000; i++ {
		u[i] = generateRandomUser(i + 1)
	}

	return u
}

func generateRandomUser(id int) User {
	return User{
		ID:       id,
		Name:     generateRandomName(),
		Username: generateRandomUsername(id),
		Email:    generateRandomEmail(id),
		Phone:    generateRandomPhone(),
		Password: generateRandomPassword(),
		Address:  generateRandomAddress(),
	}
}

func generateRandomName() string {
	firstNames := []string{
		"Alexander", "Dmitry", "Maxim", "Sergey", "Andrey",
		"Alexey", "Mikhail", "Vladimir", "Ivan", "Nikolai",
		"Ekaterina", "Olga", "Anna", "Elena", "Natalia",
		"Irina", "Maria", "Anastasia", "Yulia", "Svetlana",
	}

	lastNames := []string{
		"Ivanov", "Smirnov", "Kuznetsov", "Popov", "Vasiliev",
		"Petrov", "Sokolov", "Mikhailov", "Novikov", "Fedorov",
		"Morozov", "Volkov", "Alekseev", "Lebedev", "Semenov",
		"Egorov", "Pavlov", "Kozlov", "Stepanov", "Nikolayev",
	}

	return firstNames[rand.IntN(len(firstNames))] + " " + lastNames[rand.IntN(len(lastNames))]
}

func generateRandomUsername(id int) string {
	return fmt.Sprintf("user%d", id)
}

func generateRandomEmail(id int) string {
	domains := []string{"gmail.com", "yahoo.com", "outlook.com", "example.com", "test.com", "yandex.ru", "mail.ru"}
	return fmt.Sprintf("user%d@%s", id, domains[rand.IntN(len(domains))])
}

func generateRandomPhone() string {
	return fmt.Sprintf("+7-%d-%d-%d-%d",
		rand.IntN(900)+100, // 100-999
		rand.IntN(900)+100, // 100-999
		rand.IntN(90)+10,   // 10-99
		rand.IntN(90)+10)   // 10-99
}

func generateRandomPassword() string {
	data := make([]byte, 16)
	for i := range data {
		data[i] = byte(rand.IntN(256))
	}

	hash := md5.Sum(data)

	return hex.EncodeToString(hash[:])
}

func generateRandomAddress() string {
	streets := []string{
		"Lenina St", "Pushkina St", "Gagarina St", "Mira St",
		"Lenina Ave", "Pobedy Ave", "Mira Ave", "Sovetsky Ave",
		"Kirova Lane", "Sovetsky Lane", "Shkolny Lane",
		"Tsentralny Blvd", "Molodezhny Blvd", "Rechnaya Embankment",
		"Tverskaya St", "Arbat St", "Nevsky Ave", "Tverskoy Blvd",
	}

	cities := []string{
		"Moscow", "Saint Petersburg", "Novosibirsk", "Yekaterinburg",
		"Kazan", "Nizhny Novgorod", "Chelyabinsk", "Samara",
		"Omsk", "Rostov-on-Don", "Ufa", "Krasnoyarsk",
		"Voronezh", "Perm", "Volgograd", "Krasnodar",
	}

	return fmt.Sprintf("%s, %s, bld. %d, apt. %d",
		cities[rand.IntN(len(cities))],
		streets[rand.IntN(len(streets))],
		rand.IntN(200)+1, // building 1-200
		rand.IntN(500)+1) // apartment 1-500
}
