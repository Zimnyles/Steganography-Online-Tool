package transcribe

import (
	"bufio"
	"fmt"
	"image/png"
	"os"

	"github.com/auyer/steganography"
)

func TranscribeImage(filepath string) (string, error) {
	inFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer inFile.Close()


	reader := bufio.NewReader(inFile)
	img, err := png.Decode(reader)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	sizeOfMessage := steganography.GetMessageSizeFromImage(img)

	msg := steganography.Decode(sizeOfMessage, img)
	stringMsg := btos(msg)

	return stringMsg, nil
}

func btos(c []byte) string {
	// Создаём новый слайс для ненулевых байтов
	var result []byte
	for _, b := range c {
		if b != 0 {
			result = append(result, b)
		}
	}
	return string(result)
}

