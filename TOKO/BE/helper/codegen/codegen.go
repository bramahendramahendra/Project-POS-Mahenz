package codegen

import (
	"fmt"
	"unicode"

	"pos_api/errors"
)

// BuildLetterPrefix menghasilkan prefix huruf besar sepanjang length dari name,
// dipad dengan 'X' bila jumlah huruf pada name kurang dari length.
func BuildLetterPrefix(name string, length int) string {
	letters := make([]rune, 0, length)
	for _, r := range name {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToUpper(r))
		}
		if len(letters) >= length {
			break
		}
	}
	for len(letters) < length {
		letters = append(letters, 'X')
	}
	return string(letters)
}

// UniqueByPrefix mencoba prefix apa adanya, lalu prefix2, prefix3, ... sampai
// checkExists mengembalikan false, maksimal 99 percobaan.
func UniqueByPrefix(prefix string, checkExists func(code string) (bool, error)) (string, error) {
	exists, err := checkExists(prefix)
	if err != nil {
		return "", err
	}
	if !exists {
		return prefix, nil
	}

	for i := 2; i <= 99; i++ {
		candidate := fmt.Sprintf("%s%d", prefix, i)
		exists, err := checkExists(candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}

	return "", &errors.InternalServerError{Message: fmt.Sprintf("tidak bisa generate kode unik untuk prefix %q", prefix)}
}

// Sequential menghasilkan kode "prefix%0Nd" mulai dari startFrom, retry sampai
// checkExists mengembalikan false, dibatasi maxAttempts percobaan agar tidak infinite loop.
func Sequential(prefix string, digits int, startFrom int, checkExists func(code string) (bool, error)) (string, error) {
	const maxAttempts = 10000
	format := fmt.Sprintf("%s%%0%dd", prefix, digits)

	for i := 0; i < maxAttempts; i++ {
		code := fmt.Sprintf(format, startFrom+i)
		exists, err := checkExists(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}

	return "", &errors.InternalServerError{Message: fmt.Sprintf("tidak bisa generate kode unik sequential untuk prefix %q", prefix)}
}
