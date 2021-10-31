package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var (
	inputSnapshotPath = flag.String("input", "", "Snapshot input directory")
	csvOutputPath     = flag.String("csv", "dump.csv", "Path to the output CSV file")
	plotOutputPath    = flag.String("plot", "plot.png", "Path to the output plot (only PNG is supported)")
)

func main() {
	flag.Parse()

	if *inputSnapshotPath == "" {
		log.Fatal("no \"input\" argument was provided")
	}

	csvOutputFile, err := os.OpenFile(*csvOutputPath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer csvOutputFile.Close()

	csvWriter := newCSVWriter(*csvOutputFile)
	csvWriter.Write([]string{"sound", "illuminance", "voltage"})

	var snapshot snapshotContents
	if err := readJSONFile(*inputSnapshotPath, &snapshot); err != nil {
		log.Fatal(err)
	}

	for _, snapshotEntry := range snapshot.Data {
		if err := csvWriter.Write([]string{
			fmt.Sprintf("%.2f", snapshotEntry.Sound),
			fmt.Sprintf("%.2f", snapshotEntry.Illuminance),
			fmt.Sprintf("%.2f", snapshotEntry.Voltage),
		}); err != nil {
			log.Fatal(err)
		}
	}

	csvWriter.Flush()

	if err := makePlot(snapshot, *plotOutputPath); err != nil {
		log.Fatal(err)
	}
}

// Чтение содержимого JSON файла и создание структуры языка на его основе (маппинг)
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

// Создание объекта для записи в CSV файл
func newCSVWriter(file os.File) *csv.Writer {
	writer := csv.NewWriter(&file)
	writer.Comma = ','
	writer.UseCRLF = false
	return writer
}

// Создание графика по данным из файла. График сохраняется по указанному в аргументе пути
func makePlot(snapshot snapshotContents, filepath string) error {
	p := plot.New()

	p.Title.Text = "WB Snapshot Data"

	var illuminancePts, voltagePts, soundPts plotter.XYs

	for i, entry := range snapshot.Data {
		illuminancePts = append(illuminancePts, plotter.XY{X: float64(i), Y: entry.Illuminance})
		voltagePts = append(voltagePts, plotter.XY{X: float64(i), Y: entry.Voltage})
		soundPts = append(soundPts, plotter.XY{X: float64(i), Y: entry.Sound})
	}

	p.Legend.Top = true

	if err := plotutil.AddLines(p,
		"Illuminance", illuminancePts,
		"Voltage", voltagePts,
		"Sound", soundPts); err != nil {
		return err
	}

	if err := p.Save(vg.Inch*8, vg.Inch*8, filepath); err != nil {
		return err
	}

	return nil
}