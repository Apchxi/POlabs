package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/shirou/gopsutil/disk"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nВыберите опцию:")
		fmt.Println("1. Вывести информацию о логических дисках")
		fmt.Println("2. Работа с файлами")
		fmt.Println("3. Работа с JSON")
		fmt.Println("4. Работа с XML")
		fmt.Println("5. Работа с ZIP-архивами")
		fmt.Println("0. Выход")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			displayDiskInfo()
		case "2":
			handleFileOperations("", reader)
		case "3":
			handleJSON(reader)
		case "4":
			handleXML(reader)
		case "5":
			handleZip(reader)
		case "0":
			fmt.Println("Выход из программы.")
			return
		default:
			fmt.Println("Некорректный ввод. Попробуйте снова.")
		}
	}
}

// Функция для вывода информации о дисках
func displayDiskInfo() {
	partitions, err := disk.Partitions(true)
	if err != nil {
		fmt.Println("Ошибка получения информации о разделах дисков:", err)
		return
	}

	fmt.Println("=== Информация о логических дисках ===")

	for _, partition := range partitions {
		fmt.Printf("Диск: %s\n", partition.Device)
		fmt.Printf("  Монтирован: %s\n", partition.Mountpoint)
		fmt.Printf("  Тип файловой системы: %s\n", partition.Fstype)
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			fmt.Println("Ошибка получения информации о диске:", err)
			continue
		}
		fmt.Printf("  Общий размер: %.2f GB\n", float64(usage.Total)/1e9)
		fmt.Printf("  Свободно: %.2f GB\n", float64(usage.Free)/1e9)
		fmt.Printf("  Использовано: %.2f GB\n", float64(usage.Used)/1e9)
		fmt.Printf("  Процент использования: %.2f%%\n\n", usage.UsedPercent)
	}

	disks, err := disk.IOCounters()
	if err != nil {
		fmt.Println("Ошибка получения статистики устройств:", err)
		return
	}

	fmt.Println("=== Информация о подключенных устройствах ===")
	for device, ioStats := range disks {
		fmt.Printf("Устройство: %s\n", device)
		fmt.Printf("  Прочитано байт: %d\n", ioStats.ReadBytes)
		fmt.Printf("  Записано байт: %d\n", ioStats.WriteBytes)
		fmt.Printf("  Прочитано операций: %d\n", ioStats.ReadCount)
		fmt.Printf("  Записано операций: %d\n", ioStats.WriteCount)
		fmt.Printf("  Время ожидания операций (в секундах): %.2f\n", ioStats.IoTime/1000.0)
		fmt.Println()
	}
}

// 2. Работа с текстовым файлом
func handleFileOperations(fileType string, reader *bufio.Reader) {
	for {
		fmt.Printf("\n=== Работа с %s файлом ===\n", fileType)
		fmt.Println("1. Создать файл")
		fmt.Println("2. Записать данные в файл")
		fmt.Println("3. Прочитать файл")
		fmt.Println("4. Удалить файл")
		fmt.Println("0. Вернуться в главное меню")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			createFile(reader, fileType)
		case "2":
			writeToFile(reader, fileType)
		case "3":
			readFile(reader, fileType)
		case "4":
			deleteFile(reader, fileType)
		case "0":
			return
		default:
			fmt.Println("Некорректный ввод. Попробуйте снова.")
		}
	}
}

func createFile(reader *bufio.Reader, fileType string) {
	fmt.Printf("Введите имя %s файла: ", fileType)
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	fmt.Printf("%s файл '%s' успешно создан.\n", fileType, fileName)
}

func writeToFile(reader *bufio.Reader, fileType string) {
	fmt.Printf("Введите имя %s файла: ", fileType)
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Printf("Файл '%s' не существует. Создайте его сначала.\n", fileName)
		return
	}

	fmt.Printf("Введите данные для записи в %s файл: ", fileType)
	content, _ := reader.ReadString('\n')

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
	} else {
		fmt.Printf("Данные успешно записаны в %s файл '%s'.\n", fileType, fileName)
	}
}

func readFile(reader *bufio.Reader, fileType string) {
	fmt.Printf("Введите имя %s файла для чтения: ", fileType)
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
	} else {
		fmt.Printf("Содержимое %s файла:\n", fileType)
		fmt.Println(string(content))
	}
}

func deleteFile(reader *bufio.Reader, fileType string) {
	fmt.Printf("Введите имя %s файла для удаления: ", fileType)
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	err := os.Remove(fileName)
	if err != nil {
		fmt.Println("Ошибка удаления файла:", err)
	} else {
		fmt.Printf("%s файл '%s' успешно удален.\n", fileType, fileName)
	}
}

// 3. Работа с JSON
func handleJSON(reader *bufio.Reader) {
	fmt.Println("\n=== Работа с JSON ===")
	for {
		fmt.Println("\n1. Создать JSON файл")
		fmt.Println("2. Записать данные в JSON файл")
		fmt.Println("3. Прочитать JSON файл")
		fmt.Println("4. Удалить JSON файл")
		fmt.Println("0. Вернуться в главное меню")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			createJSONFile(reader)
		case "2":
			writeToJSON(reader)
		case "3":
			readJSONFile(reader)
		case "4":
			deleteFile(reader, "JSON")
		case "0":
			return
		default:
			fmt.Println("Некорректный ввод. Попробуйте снова.")
		}
	}
}

func createJSONFile(reader *bufio.Reader) {
	fmt.Print("Введите имя JSON файла: ")
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	if !strings.HasSuffix(fileName, ".json") {
		fmt.Println("Ошибка: Неверное расширение файла. Ожидается .json")
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	fmt.Printf("JSON файл '%s' успешно создан.\n", fileName)
}

func writeToJSON(reader *bufio.Reader) {
	fmt.Print("Введите имя JSON файла: ")
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	if !strings.HasSuffix(fileName, ".json") {
		fmt.Println("Ошибка: Неверное расширение файла. Ожидается .json")
		return
	}

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	fmt.Print("Введите имя: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Введите возраст: ")
	var age int
	fmt.Scan(&age)

	person := Person{Name: name, Age: age}
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(person)
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
	} else {
		fmt.Printf("Данные успешно записаны в JSON файл '%s'.\n", fileName)
	}
}

func readJSONFile(reader *bufio.Reader) {
	fmt.Print("Введите имя JSON файла для чтения: ")
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	if !strings.HasSuffix(fileName, ".json") {
		fmt.Println("Ошибка: Неверное расширение файла. Ожидается .json")
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	var person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&person)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	fmt.Printf("Данные из JSON файла:\nИмя: %s\nВозраст: %d\n", person.Name, person.Age)
}

// 4. Работа с XML
func handleXML(reader *bufio.Reader) {
	fmt.Println("\n=== Работа с XML ===")
	for {
		fmt.Println("\n1. Создать XML файл")
		fmt.Println("2. Записать данные в XML файл")
		fmt.Println("3. Прочитать XML файл")
		fmt.Println("4. Удалить XML файл")
		fmt.Println("0. Вернуться в главное меню")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			createXMLFile(reader)
		case "2":
			writeToXML(reader)
		case "3":
			readXMLFile(reader)
		case "4":
			deleteFile(reader, "XML")
		case "0":
			return
		default:
			fmt.Println("Некорректный ввод. Попробуйте снова.")
		}
	}
}

func createXMLFile(reader *bufio.Reader) {
	fmt.Print("Введите имя XML файла: ")
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	if !strings.HasSuffix(fileName, ".xml") {
		fmt.Println("Ошибка: Неверное расширение файла. Ожидается .xml")
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	fmt.Printf("XML файл '%s' успешно создан.\n", fileName)
}

func writeToXML(reader *bufio.Reader) {
	fmt.Print("Введите имя XML файла: ")
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	if !strings.HasSuffix(fileName, ".xml") {
		fmt.Println("Ошибка: Неверное расширение файла. Ожидается .xml")
		return
	}

	type Person struct {
		XMLName xml.Name `xml:"person"`
		Name    string   `xml:"name"`
		Age     int      `xml:"age"`
	}

	fmt.Print("Введите имя: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Введите возраст: ")
	var age int
	fmt.Scan(&age)

	person := Person{Name: name, Age: age}
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	err = encoder.Encode(person)
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
	} else {
		fmt.Printf("Данные успешно записаны в XML файл '%s'.\n", fileName)
	}
}

func readXMLFile(reader *bufio.Reader) {
	fmt.Print("Введите имя XML файла для чтения: ")
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	if !strings.HasSuffix(fileName, ".xml") {
		fmt.Println("Ошибка: Неверное расширение файла. Ожидается .xml")
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	var person struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
	}

	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&person)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	fmt.Printf("Данные из XML файла:\nИмя: %s\nВозраст: %d\n", person.Name, person.Age)
}

// 5. Работа с ZIP-архивами
func handleZip(reader *bufio.Reader) {
	for {
		fmt.Println("\n=== Работа с ZIP-архивами ===")
		fmt.Println("1. Создать ZIP-архив")
		fmt.Println("2. Добавить файл в архив")
		fmt.Println("3. Разархивировать архив")
		fmt.Println("4. Удалить архив и файл")
		fmt.Println("0. Вернуться в главное меню")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			createZip(reader)
		case "2":
			addFileToZip(reader)
		case "3":
			extractZip(reader)
		case "4":
			deleteZipArchive(reader)
		case "0":
			return
		default:
			fmt.Println("Некорректный ввод. Попробуйте снова.")
		}
	}
}

func createZip(reader *bufio.Reader) {
	fmt.Print("Введите имя ZIP архива: ")
	archiveName, _ := reader.ReadString('\n')
	archiveName = strings.TrimSpace(archiveName)

	if !strings.HasSuffix(archiveName, ".zip") {
		fmt.Println("Ошибка: Неверное расширение архива. Ожидается .zip")
		return
	}

	archiveFile, err := os.Create(archiveName)
	if err != nil {
		fmt.Println("Ошибка создания архива:", err)
		return
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	fmt.Printf("ZIP архив '%s' успешно создан.\n", archiveName)
}

func addFileToZip(reader *bufio.Reader) {
	fmt.Print("Введите имя ZIP архива: ")
	archiveName, _ := reader.ReadString('\n')
	archiveName = strings.TrimSpace(archiveName)

	// Проверка расширения архива
	if !strings.HasSuffix(archiveName, ".zip") {
		fmt.Println("Ошибка: Неверное расширение архива. Ожидается .zip")
		return
	}

	// Проверка существует ли файл
	if _, err := os.Stat(archiveName); os.IsNotExist(err) {
		// Если файл не существует, создаем новый архив
		fmt.Printf("Архив '%s' не найден. Создается новый архив...\n", archiveName)
		archiveFile, err := os.Create(archiveName)
		if err != nil {
			fmt.Println("Ошибка создания архива:", err)
			return
		}
		defer archiveFile.Close()

		// Создаем новый архивный писатель
		zipWriter := zip.NewWriter(archiveFile)
		defer zipWriter.Close()

		fmt.Printf("ZIP архив '%s' успешно создан.\n", archiveName)

		// Запрос на выбор файла для добавления в архив
		fmt.Print("Введите имя файла, который хотите добавить в архив: ")
		fileName, _ := reader.ReadString('\n')
		fileName = strings.TrimSpace(fileName)

		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			fmt.Println("Файл не существует.")
			return
		}

		// Открытие файла для добавления в архив
		fileToAdd, err := os.Open(fileName)
		if err != nil {
			fmt.Println("Ошибка открытия файла для добавления в архив:", err)
			return
		}
		defer fileToAdd.Close()

		// Получение информации о файле
		fileInfo, err := fileToAdd.Stat()
		if err != nil {
			fmt.Println("Ошибка получения информации о файле:", err)
			return
		}

		// Добавление файла в архив
		zipFile, err := zipWriter.Create(fileInfo.Name())
		if err != nil {
			fmt.Println("Ошибка добавления файла в архив:", err)
			return
		}

		_, err = io.Copy(zipFile, fileToAdd)
		if err != nil {
			fmt.Println("Ошибка копирования файла в архив:", err)
			return
		}

		fmt.Printf("Файл '%s' успешно добавлен в архив '%s'.\n", fileName, archiveName)

	} else {
		// Если архив существует, открываем его для записи
		archiveFile, err := os.OpenFile(archiveName, os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("Ошибка открытия архива:", err)
			return
		}
		defer archiveFile.Close()

		// Создаем новый архивный писатель для добавления файлов
		zipWriter := zip.NewWriter(archiveFile)
		defer zipWriter.Close()

		// Запрос на выбор файла для добавления в архив
		fmt.Print("Введите имя файла, который хотите добавить в архив: ")
		fileName, _ := reader.ReadString('\n')
		fileName = strings.TrimSpace(fileName)

		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			fmt.Println("Файл не существует.")
			return
		}

		// Открытие файла для добавления в архив
		fileToAdd, err := os.Open(fileName)
		if err != nil {
			fmt.Println("Ошибка открытия файла для добавления в архив:", err)
			return
		}
		defer fileToAdd.Close()

		// Получение информации о файле
		fileInfo, err := fileToAdd.Stat()
		if err != nil {
			fmt.Println("Ошибка получения информации о файле:", err)
			return
		}

		// Добавление файла в архив
		zipFile, err := zipWriter.Create(fileInfo.Name())
		if err != nil {
			fmt.Println("Ошибка добавления файла в архив:", err)
			return
		}

		_, err = io.Copy(zipFile, fileToAdd)
		if err != nil {
			fmt.Println("Ошибка копирования файла в архив:", err)
			return
		}

		fmt.Printf("Файл '%s' успешно добавлен в архив '%s'.\n", fileName, archiveName)
	}
}

func extractZip(reader *bufio.Reader) {
	fmt.Print("Введите имя ZIP архива: ")
	archiveName, _ := reader.ReadString('\n')
	archiveName = strings.TrimSpace(archiveName)

	// Проверка расширения архива
	if !strings.HasSuffix(archiveName, ".zip") {
		fmt.Println("Ошибка: Неверное расширение архива. Ожидается .zip")
		return
	}

	// Открытие архива
	archiveFile, err := os.Open(archiveName)
	if err != nil {
		fmt.Println("Ошибка открытия архива:", err)
		return
	}
	defer archiveFile.Close()

	// Получаем информацию о файле
	stat, err := archiveFile.Stat()
	if err != nil {
		fmt.Println("Ошибка получения информации о файле:", err)
		return
	}

	// Чтение архива
	zipReader, err := zip.NewReader(archiveFile, stat.Size())
	if err != nil {
		fmt.Println("Ошибка чтения архива:", err)
		return
	}

	// Вывод информации о содержимом архива
	fmt.Printf("Содержимое архива '%s':\n", archiveName)
	for _, file := range zipReader.File {
		fmt.Printf("Имя файла: %s\n", file.Name)
		fmt.Printf("Размер файла: %d байт\n", file.UncompressedSize)
		fmt.Printf("Тип сжатия: %d\n", file.Method)
		fmt.Println("----------")
	}

	// Запрос на разархивирование и вывод содержимого
	fmt.Print("Введите имя файла из архива для извлечения: ")
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)

	// Поиск выбранного файла в архиве
	var fileToExtract *zip.File
	for _, file := range zipReader.File {
		if file.Name == fileName {
			fileToExtract = file
			break
		}
	}

	if fileToExtract == nil {
		fmt.Println("Файл не найден в архиве.")
		return
	}

	// Вывод данных о выбранном файле
	fmt.Printf("Извлекаем файл: %s\n", fileToExtract.Name)
	fmt.Printf("Размер извлекаемого файла: %d байт\n", fileToExtract.UncompressedSize)
	fmt.Println("Данные о файле будут распакованы...")

	// Разархивирование файла
	outputFile, err := os.Create(fileToExtract.Name)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer outputFile.Close()

	fileReader, err := fileToExtract.Open()
	if err != nil {
		fmt.Println("Ошибка открытия файла для чтения:", err)
		return
	}
	defer fileReader.Close()

	_, err = io.Copy(outputFile, fileReader)
	if err != nil {
		fmt.Println("Ошибка разархивирования файла:", err)
		return
	}

	fmt.Printf("Файл '%s' успешно извлечен.\n", fileToExtract.Name)
}

func deleteZipArchive(reader *bufio.Reader) {
	fmt.Print("Введите имя ZIP архива для удаления: ")
	archiveName, _ := reader.ReadString('\n')
	archiveName = strings.TrimSpace(archiveName)

	if !strings.HasSuffix(archiveName, ".zip") {
		fmt.Println("Ошибка: Неверное расширение архива. Ожидается .zip")
		return
	}

	err := os.Remove(archiveName)
	if err != nil {
		fmt.Println("Ошибка удаления архива:", err)
	} else {
		fmt.Printf("ZIP архив '%s' успешно удален.\n", archiveName)
	}
}
