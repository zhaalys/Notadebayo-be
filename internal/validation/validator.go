package validation

import (
	"fmt"
	"reflect"
	"strings"
	"tasklybe/internal/dto"

	"github.com/go-playground/validator/v10" // install package validator
	"github.com/gofiber/fiber/v2"
)

// inisialisasi variable global Validate
var Validate = validator.New()

/** 
fungsi init() mirip dengan main() namun ada perbedaan yaitu:
- init() selalu dipanggil pertama sebelum main()
- main() hanya dipanggil 1x saat program jalan, namun init() dipanggil berkali -
  kali. misal kita panggil fungsi "BindAndValidate" maka yang terjadi
	A. init() --> dipanggil terlebih dahulu
	B. BindAndValidate() --> dipanggil kemudian
	jika terdapat main() maka main() tidak akan dipanggil
- main() hanya bisa pada package main, init() bisa di package mana pun
- init() bisa di definisikan di berbagai file walaupun 1 package yang sama
  sedangkan main() tidak bisa, 1 package hanya 1 main()
**/
func init() {

	/** 
	Fungsi init() ini digunakan untuk koversi penamaan struct yang sebelumnya
	"UserId" menjadi "userId" agar bisa dibaca dengan mudah oleh frontend.
	**/
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.Split(fld.Tag.Get("json"), ",")[0]

		if name == "-" {
			return ""
		}
		return name
	})
}

/**
Berfungsi untuk mencocokan request user dengan struct kita, jika user
memiliki struktur body yang berbeda maka otomatis akan response 400
**/
func BindAndValidate[T any](c *fiber.Ctx, dst *T) error {
	if err := c.BodyParser(dst); err != nil {
		return err
	}
	return Validate.Struct(dst)
}

/**
Berfungsi untuk memformat ulang hasil validasi dari package validator/v10.
**/
func FormatValidationError(err error) *[]dto.ResponseError {
	errors := []dto.ResponseError{}
	
	// jika ada field yang tidak sesuai maka blok kode ini jalan
	if ve, ok := err.(validator.ValidationErrors); ok {
	
		// Ini adalah looping field dari request yang error. Disini kita hanya
		// menambahkan data ke variable errors
		for _, e := range ve {
			errors = append(errors, dto.ResponseError{
				Field:   e.Field(),
				Value:   fmt.Sprint(e.Value()),
				Tag:     e.Tag(),
				Message: fmt.Sprintf("Field '%s' is %s", e.Field(), e.Tag()),
				Target:  "/task",
			})
		}
	}

	return &errors
}