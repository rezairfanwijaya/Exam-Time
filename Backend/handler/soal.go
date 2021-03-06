package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rg-km/final-project-engineering-46/helper"
	"github.com/rg-km/final-project-engineering-46/soal"
	tokensoal "github.com/rg-km/final-project-engineering-46/token-soal"
)

// struct handlerSoal
type handlerSoal struct {
	soal      soal.Service
	tokenSoal tokensoal.Service
}

// function NewHandlerSoal
func NewHandlerSoal(soal soal.Service, tokenSoal tokensoal.Service) *handlerSoal {
	return &handlerSoal{soal, tokenSoal}
}

// function buat handler create soal
func (h *handlerSoal) CreateSoal(c *gin.Context) {

	// authorization
	user := helper.IsGuru(c)
	if user.Role != "guru" {
		return
	}

	// inisiasi input soal
	var inputSoal soal.InputSoal

	// binding
	err := c.ShouldBindJSON(&inputSoal)
	if err != nil {
		// ambil error binding
		myErr := helper.ErrorBinding(err)
		response := helper.ResponsAPI("Gagal binding", "Gagal", http.StatusBadRequest, myErr)
		c.JSON(http.StatusBadRequest, response)
		return

	}

	// panggil function create soal
	err = h.soal.CreateSoal(inputSoal, user.Id_users)
	if err != nil {
		myErr := gin.H{
			"error": err.Error(),
		}
		response := helper.ResponsAPI("Gagal membuat soal", "Gagal", http.StatusBadRequest, myErr)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// response sukses
	message := gin.H{
		"message": "Berhasil membuat soal",
	}

	response := helper.ResponsAPI("Berhasil membuat soal", "Berhasil", http.StatusOK, message)
	c.JSON(http.StatusOK, response)

}

// function handler untuk menampilkan semua soal ketika token valid
func (h *handlerSoal) ShowAllSoalSiswa(c *gin.Context) {
	// authorization
	currentUser := helper.IsSiswa(c)
	if currentUser.Role != "siswa" {
		return
	}

	// input token dari siswa
	var input tokensoal.InputTokenSiswa

	// binding
	err := c.ShouldBindJSON(&input)
	if err != nil {
		myErr := helper.ErrorBinding(err)
		response := helper.ResponsAPI("gagal mengambil token", "gagal", http.StatusBadRequest, myErr)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// validasi token
	mapel, err := h.tokenSoal.ValidasiTokenSoal(currentUser.Id_users, input.Token)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		response := helper.ResponsAPI("Gagal validasi token", "gagal", http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// tampilkan semua soal
	soals, err := h.soal.ShowAllSoalSiswa(mapel.IdMataPelajaran)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		response := helper.ResponsAPI("Gagal mengambil soal", "gagal", http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"mapel":    mapel.MataPelajaran,
		"id_mapel": mapel.IdMataPelajaran,
		"durasi":   mapel.Durasi,
		"soal":     soals,
	}

	response := helper.ResponsAPI("Sukses mengambil soal", "sukses", http.StatusOK, data)
	c.JSON(http.StatusOK, response)
}

// handler untuk show all soal guru (bank soal)
func (h *handlerSoal) ShowAllSoalGuru(c *gin.Context) {
	// authorization
	user := helper.IsGuru(c)
	if user.Role != "guru" {
		return
	}

	// input
	var input soal.InputBankSoal

	// binding
	err := c.ShouldBindJSON(&input)
	if err != nil {
		myErr := helper.ErrorBinding(err)
		response := helper.ResponsAPI("gagal binding", "gagal", http.StatusBadRequest, myErr)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// panggil function diservice
	soals, err := h.soal.ShowAllSoalGuru(input.IdMataPelajaran)
	if err != nil {
		data := gin.H{
			"error": err.Error(),
		}
		response := helper.ResponsAPI("gagal mengambil data soal", "gagal", http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(soals) == 0 {
		response := helper.ResponsAPI("id mata pelajaran tidak terdaftar", "gagal", http.StatusBadRequest, soal.FormatterBankSoals(soals))
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// return sukses
	response := helper.ResponsAPI("berhasil mengambil soal", "sukses", http.StatusOK, soal.FormatterBankSoals(soals))
	c.JSON(http.StatusOK, response)

}
