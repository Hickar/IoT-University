package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

var (
	snapshotInputPath  = flag.String("input", "", "Path to the snapshot dataset directory")
	snapshotOutputPath = flag.String("output", "./snapshot_dump", "Path to the snapshot output directory")
)

func main() {
	flag.Parse()

	if *snapshotInputPath == "" {
		log.Fatal("no \"input\" argument was provided")
	}

	// Создание директории для JSON файлов, если её ещё не существует
	if _, err := os.Stat(*snapshotOutputPath); os.IsNotExist(err) {
		if err := os.Mkdir(*snapshotOutputPath, 0777); err != nil {
			log.Fatalf("unexpected error during directory creation: %s", err)
		}
	}

	// Строки 39-72 – чтение json файлов из датасета с последующим "вытягиванием" нужных свойств
	// и их записью в отдельный CSV файл
	datasetFiles, err := readDirFiles(*snapshotInputPath)
	if err != nil {
		log.Fatal(err)
	}

	sortFilesByName(datasetFiles)

	for _, datasetFile := range datasetFiles {
		inputFilepath := path.Join(*snapshotInputPath, datasetFile.Name())
		dataset := &snapshotDataset{}
		if err := readJSONFile(inputFilepath, &dataset); err != nil {
			log.Fatal(err)
		}

		var snapshotEntries []snapshotEntry
		for _, datasetEntry := range dataset.Data {
			snapshotEntries = append(snapshotEntries, snapshotEntry{
				Motion:      datasetEntry.Motion,
				Sound:       datasetEntry.Sound,
				Illuminance: datasetEntry.Illuminance,
				Temperature: datasetEntry.Temperature,
			})
		}

		snapshot := &snapshot{snapshotEntries}
		outputFilename := path.Join(*snapshotOutputPath, basename(datasetFile.Name()))
		if err := writeJSONFile(outputFilename+".json", &snapshot); err != nil {
			log.Fatal(err)
		}

		if err := writeXMLFile(outputFilename+".xml", &snapshot); err != nil {
			log.Fatal(err)
		}
	}

	// Строки 75-92 – чтение данных из JSON файлов, созданных ранее
	snapshotFileEntries, _ := readDirFiles(*snapshotOutputPath)
	for _, fileEntry := range snapshotFileEntries {
		if strings.Contains(fileEntry.Name(), ".json") {
			fmt.Printf("========%s========\n", fileEntry.Name())

			filepath := path.Join(*snapshotOutputPath, fileEntry.Name())

			var snapshot snapshot
			if err := readJSONFile(filepath, &snapshot); err != nil {
				log.Fatal(err)
			}

			for i, snapshotEntry := range snapshot.Data {
				fmt.Printf("[%s/%d]: %s\n", fileEntry.Name(), i, snapshotEntry.String())
			}
			fmt.Print("\n\n")
		}
	}
}

// Получение списка объектов файлов в директории по указанному пути.
// path – путь к директории.
func readDirFiles(path string) ([]fs.DirEntry, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open directory at %s: %s", path, err)
	}

	if stat, _ := dir.Stat(); !stat.IsDir() {
		return nil, errors.New("file path specified instead of directory")
	}

	files, err := dir.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("unexpecte error during directory reading: %s", err)
	}

	return files, nil
}

// Чтение содержимого JSON файла в объект по ссылке.
// filepath – путь к файлу, value – ссылка на объект.
func readJSONFile(filepath string, value interface{}) error {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read json file contents: %s", err)
	}

	return json.Unmarshal(bytes, &value)
}

// Запись структуры языка в JSON файл.
// filepath – путь к новому или существующему файлу, value – ссылка на структуру данных.
func writeJSONFile(filepath string, value interface{}) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := json.MarshalIndent(&value, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	return err
}

// Запись структуры языка в XML файл.
// filepath – путь к новому или существующему файлу, value – ссылка на структуру данных.
func writeXMLFile(filepath string, value interface{}) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := xml.MarshalIndent(&value, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	return err
}

// Сортировка массива объектов файлов по названию
func sortFilesByName(files []os.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
}

// Возвращает название файла без расширения, если же расширения нет – возвращает изначальное название
func basename(filename string) string {
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex == -1 {
		return filename
	}

	return filename[:dotIndex]
}
