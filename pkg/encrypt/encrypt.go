package encrypt

import (
	"bufio"
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"
	filenamegenerator "stegano-webapp/steagano-webapp/pkg/filenameGenerator"
	"strconv"
	"time"

	"github.com/auyer/steganography"
)

func EncryptImage(filepath string, secretMessage string) (string, string, error) {
	inFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	fmt.Println(filepath, secretMessage)

	reader := bufio.NewReader(inFile)
	img, err := png.Decode(reader)
	if err != nil {
		fmt.Println("0")
	}

	w := new(bytes.Buffer)
	err = steganography.Encode(w, img, []byte(secretMessage))
	if err != nil {
		log.Printf("Error Encoding file %v", err)
		return "", "", err
	}

	uniqueFilenameCode := filenamegenerator.GenerateFilename()
	uniqueFilename := "encrypted_image_" + uniqueFilenameCode + strconv.FormatInt((time.Now().Unix()), 10) + ".png"
	encryptedfilepath := "static/images/encryimages/" + uniqueFilename

	fmt.Println(uniqueFilename, encryptedfilepath+"1")

	outFile, _ := os.Create(encryptedfilepath)
	w.WriteTo(outFile)
	outFile.Close()

	fmt.Println(uniqueFilename, encryptedfilepath+"2")

	return encryptedfilepath, uniqueFilename, nil
}
