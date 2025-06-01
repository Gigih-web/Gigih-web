package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const maksKampanye int = 100

type Kampanye struct {
	Pembuat       string `json:"pembuat"`
	Kategori      string `json:"kategori"`
	Judul         string `json:"judul"`
	Deskripsi     string `json:"deskripsi"`
	TargetDana    int    `json: "targetDana"`
	Tenggat       string `json: "tenggat"`
	JumlahDonasi  int    `json:"jumlahDonasi"`
	JumlahDonatur int    `json:"jumlahDonatur"`
	Status        string `json: "status"`
}

var kampanye [maksKampanye]Kampanye
var nKamp int
var reader *bufio.Reader = bufio.NewReader(os.Stdin)

const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Red    = "\033[31m"
)

func simpanKeJSON() error {
	file, err := os.Create("kampanye.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(kampanye[:nKamp])
}

func muatDariJSON() error {
	file, err := os.Open("kampanye.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&kampanye)
	if err != nil {
		return err
	}

	nKamp = 0
	for _, k := range kampanye {
		if k.Judul != "" {
			nKamp++
		}
	}
	return nil
}

func selectionJumlahDonatur(arr *[maksKampanye]Kampanye, n int) {
	var i, j int
	for i = 0; i < n-1; i++ {
		maxIdx := i
		for j = i + 1; j < n; j++ {
			if arr[j].JumlahDonatur > arr[maxIdx].JumlahDonatur {
				maxIdx = j
			}
		}
		if maxIdx != i {
			arr[i], arr[maxIdx] = arr[maxIdx], arr[i]
		}
	}
}

func insertionSortDescending(arr *[maksKampanye]Kampanye, n int) {
	var i, j int
	var temp Kampanye
	for i = 1; i < n; i++ {
		temp = arr[i]
		j = i - 1
		for j >= 0 && arr[j].JumlahDonasi < temp.JumlahDonasi {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = temp
	}
}

func sortByTitleAscending(arr *[maksKampanye]Kampanye, n int) {
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if strings.ToLower(arr[i].Judul) > strings.ToLower(arr[j].Judul) {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

func binarySearchByTitle(arr *[maksKampanye]Kampanye, n int, target string) int {
	low := 0
	high := n - 1

	for low <= high {
		mid := (low + high) / 2
		if strings.ToLower(arr[mid].Judul) == strings.ToLower(target) {
			return mid
		} else if strings.ToLower(arr[mid].Judul) < strings.ToLower(target) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

func bacaInputString(prompt string) string {
	var input string
	fmt.Print(Yellow + prompt + Reset)
	input, _ = reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func garis() {
	fmt.Println(Blue + strings.Repeat("=", 50) + Reset)
}

func updateStatusKampanye() {
	var i int
	var now, tenggat time.Time
	var err error
	now = time.Now()

	for i = 0; i < nKamp; i++ {
		if kampanye[i].Status != "Aktif" {
			return
		}
		tenggat, err = time.Parse("2006-01-01", kampanye[i].Tenggat)
		if err != nil {
			return
		}

		if kampanye[i].JumlahDonasi >= kampanye[i].TargetDana {
			kampanye[i].Status = "Tercapai"
		} else if now.After(tenggat) {
			kampanye[i].Status = "Berakhir"
		}
	}
}

func ubahKampanye(judulCari string, deskripsiBaru string, targetDanaBaru int, tenggatBaru string) {
	var index int

	judulCari = bacaInputString("Masukkan judul kampanye yang ingin diubah: ")
	sortByTitleAscending(&kampanye, nKamp)
	index = binarySearchByTitle(&kampanye, nKamp, judulCari)

	if index == -1 {
		fmt.Println(Red + "Kampanye tidak ditemukan." + Reset)
		return
	}

	deskripsiBaru = bacaInputString("Masukkan deskripsi baru (kosongkan jika tidak ingin mengubah): ")
	if deskripsiBaru != "" {
		kampanye[index].Deskripsi = deskripsiBaru
	}

	fmt.Print("Masukkan target dana baru: ")
	_, err := fmt.Scan(&targetDanaBaru)
	reader.ReadString('\n')
	if err == nil && targetDanaBaru > 0 {
		kampanye[index].TargetDana = targetDanaBaru
	}

	tenggatBaru = bacaInputString("Masukkan tenggat baru (YYYY-MM-DD) (kosongkan jika tidak ingin mengubah): ")
	if tenggatBaru != "" {
		kampanye[index].Tenggat = tenggatBaru
	}

	fmt.Println(Green + "✅ Kampanye berhasil diubah." + Reset)
}

func hapusKampanye(judulCari string) {
	var index, i int

	judulCari = bacaInputString("Masukkan judul kampanye yang ingin dihapus: ")
	index = binarySearchByTitle(&kampanye, nKamp, judulCari)
	if index == -1 {
		fmt.Println(Red + "Kampanye tidak ditemukan." + Reset)
		return
	}
	for i = index; i < nKamp-1; i++ {
		kampanye[i] = kampanye[i+1]
	}
	nKamp--
	fmt.Println(Green + "✅ Kampanye berhasil dihapus." + Reset)
}

func main() {
	var err error
	var role, lanjut, ktg, donasi, dntr, targetDana int
	var kategori, judul, deskripsi, nama, tenggat string

	err = muatDariJSON()
	if err != nil {
		fmt.Println(Red + "❌ Gagal memuat data:" + err.Error() + Reset)
		return
	}
	updateStatusKampanye()

	defer func() {
		err = simpanKeJSON()
		if err != nil {
			fmt.Println(Red + "❌ Gagal menyimpan data:" + err.Error() + Reset)
		}
	}()

	garis()
	fmt.Println(Green + "🌱 Selamat datang di Platform Donasi 'Satu Tangan'" + Reset)
	fmt.Println("Membantu satu sama lain, satu tangan dalam perubahan.")
	garis()

	fmt.Println("\n--- HOME PAGE ---")
	nama = bacaInputString("Masukkan nama Anda: ")

	fmt.Println("\nSilahkan pilih role Anda:")
	fmt.Println("1. ✏  Pembuat Kampanye")
	fmt.Println("2. 💰 Donatur")
	fmt.Println("3. ✏  Edit Kampanye")
	fmt.Println("4. 🚪 Keluar")
	fmt.Print("Pilihan: ")
	fmt.Scan(&role)
	reader.ReadString('\n')

	if role == 1 {
		fmt.Println(Green + "\nHai " + nama + ", mari kita mulai membuat kampanye baru!" + Reset)
		fmt.Println("1. Ya")
		fmt.Println("2. Tidak")
		fmt.Scan(&lanjut)
		reader.ReadString('\n')

		if lanjut != 1 {
			fmt.Println("Terima kasih, sampai jumpa 👋")
			return
		}

		fmt.Println("\nPilih kategori kampanye:")
		fmt.Println("1. Kesehatan")
		fmt.Println("2. Pendidikan")
		fmt.Println("3. Bantuan Sosial")
		fmt.Scan(&ktg)
		reader.ReadString('\n')

		switch ktg {
		case 1:
			kategori = "Kesehatan"
		case 2:
			kategori = "Pendidikan"
		case 3:
			kategori = "Bantuan Sosial"
		default:
			fmt.Println(Red + "Kategori tidak valid." + Reset)
			return
		}

		judul = bacaInputString("Masukkan judul kampanye: ")
		deskripsi = bacaInputString("Masukkan deskripsi kampanye: ")
		fmt.Print("Masukan target dana: \n")
		fmt.Scan(&targetDana)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		tenggat = bacaInputString("Masukan tenggat kampanye (YYYY-MM-DD): ")

		if nKamp < maksKampanye {
			kampanye[nKamp] = Kampanye{
				Pembuat:      nama,
				Kategori:     kategori,
				Judul:        judul,
				Deskripsi:    deskripsi,
				TargetDana:   targetDana,
				Tenggat:      tenggat,
				JumlahDonasi: 0,
				Status:       "Sedang berjalan",
			}
			nKamp++
			fmt.Println(Green + "✅ Kampanye berhasil dibuat!" + Reset)
		} else {
			fmt.Println(Red + "❌ Kapasitas kampanye penuh." + Reset)
		}

	} else if role == 2 {
		var i int

		if nKamp == 0 {
			fmt.Println(Red + "❌ Belum ada kampanye tersedia." + Reset)
			return
		}

		fmt.Println(Green + "\nHai " + nama + ", berikut daftar kampanye yang tersedia:" + Reset)
		for i = 0; i < nKamp; i++ {
			fmt.Printf("%d. [%s] %s - %s (Donasi: Rp%d, Donatur: %d orang)\n", i+1, kampanye[i].Kategori, kampanye[i].Judul, kampanye[i].Deskripsi, kampanye[i].JumlahDonasi, kampanye[i].JumlahDonatur)
		}

		fmt.Println()
		fmt.Println("Menu untuk mengurutkan kampanye: ")
		fmt.Println("1. Total tertinggi donasi.")
		fmt.Println("2. Terbanyak yang mendonasi.")
		fmt.Println("3. Tidak")
		var opsiUrut int
		fmt.Scan(&opsiUrut)
		reader.ReadString('\n')

		switch opsiUrut {
		case 1:
			insertionSortDescending(&kampanye, nKamp)
			fmt.Println(Green + "\nKampanye diurutkan berdasarkan total donasi tertinggi:" + Reset)
			for i = 0; i < nKamp; i++ {
				fmt.Printf("%d. [%s] %s - %s (Donasi: Rp%d, Donatur: %d)\n",
					i+1, kampanye[i].Kategori, kampanye[i].Judul,
					kampanye[i].Deskripsi, kampanye[i].JumlahDonasi, kampanye[i].JumlahDonatur)
			}
		case 2:
			selectionJumlahDonatur(&kampanye, nKamp)
			fmt.Println(Green + "\nKampanye diurutkan berdasarkan jumlah donatur terbanyak:" + Reset)
			for i := 0; i < nKamp; i++ {
				fmt.Printf("%d. [%s] %s - %s (Donasi: Rp%d, Donatur: %d)\n",
					i+1, kampanye[i].Kategori, kampanye[i].Judul,
					kampanye[i].Deskripsi, kampanye[i].JumlahDonasi, kampanye[i].JumlahDonatur)
			}
		default:
			fmt.Println(Red + "Menampilkan kampanye tanpa urutan." + Reset)
			for i := 0; i < nKamp; i++ {
				fmt.Printf("%d. [%s] %s - %s (Donasi: Rp%d, Donatur: %d)\n",
					i+1, kampanye[i].Kategori, kampanye[i].Judul,
					kampanye[i].Deskripsi, kampanye[i].JumlahDonasi, kampanye[i].JumlahDonatur)
			}
		}

		fmt.Println("\nIngin mencari kampanye berdasarkan:")
		fmt.Println("1. Judul 🔍")
		fmt.Println("2. Kategori 📂")
		fmt.Println("3. Tidak ➡ langsung donasi")
		fmt.Print("Pilihan: ")

		var cari int
		fmt.Scan(&cari)
		reader.ReadString('\n')

		var target, kategori string
		switch cari {
		case 1:
			sortByTitleAscending(&kampanye, nKamp)
			target = bacaInputString("Masukkan judul kampanye: ")
			index := binarySearchByTitle(&kampanye, nKamp, target)

			if index != -1 {
				fmt.Printf(Green+"Kampanye ditemukan: [%s] %s - %s (Donasi: Rp%d, Donatur: %d orang)\n"+Reset,
					kampanye[index].Kategori, kampanye[index].Judul, kampanye[index].Deskripsi,
					kampanye[index].JumlahDonasi, kampanye[index].JumlahDonatur)
			} else {
				fmt.Println(Red + "❌ Kampanye tidak ditemukan." + Reset)
			}

		case 2:
			kategori = bacaInputString("Masukkan nama kategori (Kesehatan/Pendidikan/Bantuan Sosial): ")
			found := false
			fmt.Println(Green + "\n📂 Kampanye dengan kategori: " + kategori + Reset)
			for i = 0; i < nKamp; i++ {
				if strings.EqualFold(kampanye[i].Kategori, kategori) {
					fmt.Printf("%d [%s] %s - %s (Donasi: Rp%d, Donatur: %d orang)\n", 1+i,
						kampanye[i].Kategori, kampanye[i].Judul, kampanye[i].Deskripsi,
						kampanye[i].JumlahDonasi, kampanye[i].JumlahDonatur)
					found = true
				}
			}

			if !found {
				fmt.Println(Red + "❌ Tidak ada kampanye dengan kategori tersebut." + Reset)
			}
		}

		if cari == 1 {
			sortByTitleAscending(&kampanye, nKamp)
			target := bacaInputString("Masukkan judul kampanye: ")
			index := binarySearchByTitle(&kampanye, nKamp, target)

			if index != -1 {
				fmt.Printf(Green+"Kampanye ditemukan: [%s] %s - %s (Donasi: Rp%d, Donatur: %d orang)\n"+Reset,
					kampanye[index].Kategori, kampanye[index].Judul, kampanye[index].Deskripsi, kampanye[index].JumlahDonasi, kampanye[index].JumlahDonatur)
			} else {
				fmt.Println(Red + "❌ Kampanye tidak ditemukan." + Reset)
			}
		}

		fmt.Print("\nMasukkan nomor kampanye yang ingin didonasikan: ")
		fmt.Scan(&dntr)
		reader.ReadString('\n')

		if dntr < 1 || dntr > nKamp {
			fmt.Println(Red + "❌ Pilihan tidak valid." + Reset)
			return
		}

		if kampanye[dntr-1].Status != "Sedang berjalan" {
			fmt.Println("Kampanye ini sudah tidak aktif atau sudah selesai.")
			return
		}

		fmt.Print("Masukkan jumlah donasi: Rp.")
		fmt.Scan(&donasi)
		reader.ReadString('\n')

		kampanye[dntr-1].JumlahDonasi += donasi
		kampanye[dntr-1].JumlahDonatur++
		fmt.Printf(Green+"🙏 Terima kasih %s telah berdonasi di kampanye '%s'!\n"+Reset, nama, kampanye[dntr-1].Judul)
		fmt.Printf("💰 Total donasi terkumpul sekarang: Rp%d\n", kampanye[dntr-1].JumlahDonasi)

	} else if role == 3 {
		var judulBaru, deskripsiBaru, tenggatBaru string
		var targetDanaBaru, pilihan, i int
		var namaEdit string

		namaEdit = bacaInputString("Masukkan nama Anda sebagai pembuat kampanye: ")
		fmt.Println(Green + "\nKampanye yang Anda buat:" + Reset)
		found := false

		for i = 0; i < nKamp; i++ {
			if strings.EqualFold(kampanye[i].Pembuat, namaEdit) {
				fmt.Printf("%d. [%s] %s - %s (Target: %d, Donasi: %d, Tenggat: %s, Status: %s)\n",
					i+1, kampanye[i].Kategori, kampanye[i].Judul, kampanye[i].Deskripsi,
					kampanye[i].TargetDana, kampanye[i].JumlahDonasi, kampanye[i].Tenggat, kampanye[i].Status)
				found = true
			}
		}

		if !found {
			fmt.Println(Red + "❌ Tidak ada kampanye yang ditemukan atas nama tersebut." + Reset)
			return
		}
		fmt.Println("1. Ubah kampanye")
		fmt.Println("2. Hapus kampanye")
		fmt.Println("3. Keluar")
		fmt.Print("Pilihan: ")
		fmt.Scan(&pilihan)
		reader.ReadString('\n')

		switch pilihan {
		case 1:
			ubahKampanye(judulBaru, deskripsiBaru, targetDanaBaru, tenggatBaru)
		case 2:
			hapusKampanye(judul)
		default:
			fmt.Println("Keluar...")
		}

	} else {
		fmt.Println("Keluar dari program. Sampai jumpa 👋")
		return
	}
}
