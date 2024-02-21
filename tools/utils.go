package tools

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"latihandatabasegolang/models"
	"math/big"
)

func CetakDataDiTerminal(dataDiterima models.User) {
	dataFormatJSON, err := json.MarshalIndent(dataDiterima, "", "  ")
	if err != nil {
		fmt.Println("Error mencetak JSON:", err)
		return
	}

	fmt.Println(string(dataFormatJSON))
}

func GenerateRandomString(panjang int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	charsetLen := big.NewInt(int64(len(charset)))
	result := make([]byte, panjang)

	for i := range result {
		index, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return ""
		}
		result[i] = charset[index.Int64()]
	}

	return string(result)

}
