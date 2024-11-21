package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Функция для получения MD5 хеша
func md5Hash(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// Функция для получения SHA-256 хеша
func sha256Hash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// Проверка пароля на соответствие хэшам
func checkPassword(password string, hashes map[string]struct{}) bool {
	md5HashResult := md5Hash(password)
	sha256HashResult := sha256Hash(password)
	_, exists := hashes[md5HashResult]
	if exists {
		return true
	}
	_, exists = hashes[sha256HashResult]
	return exists
}

// Генерация всех пятибуквенных паролей
func generatePasswords() []string {
	var passwords []string
	letters := "abcdefghijklmnopqrstuvwxyz"
	for a := 0; a < len(letters); a++ {
		for b := 0; b < len(letters); b++ {
			for c := 0; c < len(letters); c++ {
				for d := 0; d < len(letters); d++ {
					for e := 0; e < len(letters); e++ {
						passwords = append(passwords, string(letters[a])+string(letters[b])+string(letters[c])+string(letters[d])+string(letters[e]))
					}
				}
			}
		}
	}
	return passwords
}

// Многопоточная версия алгоритма полного перебора
func bruteForceMultiThread(hashes map[string]struct{}, numThreads int) {
	startTime := time.Now()
	passwords := generatePasswords()

	var wg sync.WaitGroup
	passwordsPerThread := len(passwords) / numThreads
	ch := make(chan string)

	// Создаем пул потоков
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			threadStartTime := time.Now() // Засекаем время на выполнение этого потока
			start := threadID * passwordsPerThread
			end := (threadID + 1) * passwordsPerThread
			if threadID == numThreads-1 {
				end = len(passwords)
			}

			// Перебор паролей в потоке
			for _, password := range passwords[start:end] {
				if checkPassword(password, hashes) {
					elapsed := time.Since(threadStartTime) // Время на нахождение пароля
					ch <- fmt.Sprintf("Поток %d - Пароль найден: %s (Время поиска: %d мс)", threadID+1, password, elapsed.Milliseconds())
				}
			}
		}(i)
	}

	// Ожидаем завершения всех горутин
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Выводим результаты из канала
	for result := range ch {
		fmt.Println(result)
	}

	// Выводим время на весь процесс
	fmt.Printf("Общее время выполнения (многопоточность): %s\n", time.Since(startTime))
}

func main() {
	// Пример хэшей из задания
	hashes := map[string]struct{}{
		"1115dd800feaacefdf481f1f9070374a2a81e27880f187396db67958b207cbad": {},
		"3a7bd3e2360a3d29eea436fcfb7e44c735d117c42d1c1835420b6b9942dd4f1b": {},
		"74e1bb62f8dabb8125a58852b63bdf6eaef667cb56ac7f7cdba6d7305c50a22f": {},
		"7a68f09bd992671bb3b19a5e70b7827e":                                 {},
	}

	var numThreads int
	fmt.Println("Введите количество потоков:")
	fmt.Scanln(&numThreads)

	if numThreads < 1 {
		fmt.Println("Количество потоков должно быть не менее 1.")
		return
	}

	bruteForceMultiThread(hashes, numThreads)
}
