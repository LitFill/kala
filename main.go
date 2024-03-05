package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Tanggal time.Time
	Nama    string
	Kelas   string
	Jam     int
	Ket     string
}

type CatatanKehadiran struct {
	Tanggal time.Time
	Ket     string
}

type DataKehadiran struct {
	Ketidakhadiran []CatatanKehadiran
	JumlahGhoib    int
	JumlahIzin     int
}

type Data map[string]*DataKehadiran

func main() {
	catatan := []Record{}
	if len(os.Args) < 2 {
		fmt.Println("Usage: kala <csv filename>")
		return
	}
	namaFile := os.Args[1]
	file, err := os.Open(namaFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return
		}
		input := Record{}
		for i, field := range record {
			switch i {
			case 0:
				input.Tanggal, err = stringToDate(field)
				if err != nil {
					fmt.Printf("Error parsing date: %v\n", err)
					return
				}
			case 1:
				input.Nama = field
			case 2:
				input.Kelas = field
			case 3:
				input.Jam, _ = strconv.Atoi(field)
			case 4:
				input.Ket = field
			}
		}
		catatan = append(catatan, input)
	}
	fmt.Println("Jumlah data: ", len(catatan))
	fmt.Println()
	data := Data{}
	for _, v := range catatan {
		nama := v.Nama
		if _, ok := data[nama]; !ok {
			data[nama] = &DataKehadiran{}
		}
		catatKetidakhadiran := CatatanKehadiran{Tanggal: v.Tanggal, Ket: v.Ket}
		data[nama].Ketidakhadiran = append(data[nama].Ketidakhadiran, catatKetidakhadiran)
		switch v.Ket {
		case "IZIN":
			data[nama].JumlahIzin++
		case "GHOIB":
			data[nama].JumlahGhoib++
		}
	}
	displayData(data)
}

func stringToDate(date string) (time.Time, error) {
	layouts := []string{"2/1/2006", "2 1 2006"}
	var galatAkhir error
	for _, layout := range layouts {
		t, err := time.Parse(layout, date)
		if err == nil {
			return t, nil
		}
		galatAkhir = err
	}
	if strings.Contains(galatAkhir.Error(), "cannot parse") {
		return time.Time{}, fmt.Errorf("invalid date format: %s", date)
	}
	return time.Time{}, galatAkhir
}

func displayData(data Data) {
	for nama, kehadiran := range data {
		fmt.Println("Nama:", nama)
		fmt.Println("Jumlah Izin:", kehadiran.JumlahIzin)
		fmt.Println("Jumlah Ghoib:", kehadiran.JumlahGhoib)
		fmt.Println("Ketidakhadiran:")
		for _, catatan := range kehadiran.Ketidakhadiran {
			fmt.Printf("\t%s: %s\n", catatan.Tanggal.Format("2006-01-02"), catatan.Ket)
		}
		fmt.Println()
	}
}
