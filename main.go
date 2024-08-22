package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	godotenv "github.com/joho/godotenv"
)

var (
	PATH_UPLOADER_MOD string
	OCULUS_APP_ID     string
	OCULUS_APP_SECRET string
	PATH_APK          string
	PATH_OBB          string
	listArgs          = make(map[string]string, 0)
)

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PATH_UPLOADER_MOD = os.Getenv("PATH_UPLOADER_MOD")
	if PATH_UPLOADER_MOD == "" {
		log.Fatalln("PATH_UPLOADER_MOD not found")
	} else {
		listArgs["PATH_UPLOADER_MOD"] = PATH_UPLOADER_MOD
	}

	OCULUS_APP_ID = os.Getenv("OCULUS_APP_ID")
	if OCULUS_APP_ID == "" {
		log.Fatalln("OCULUS_APP_ID not found")
	} else {
		listArgs["OCULUS_APP_ID"] = OCULUS_APP_ID
	}

	OCULUS_APP_SECRET = os.Getenv("OCULUS_APP_SECRET")
	if OCULUS_APP_ID == "" {
		log.Fatalln("OCULUS_APP_SECRET not found")
	} else {
		listArgs["OCULUS_APP_SECRET"] = OCULUS_APP_SECRET
	}

	PATH_APK = os.Getenv("PATH_APK")
	if PATH_APK == "" {
		log.Fatalln("PATH_APK not found")
	} else {
		listArgs["PATH_APK"] = PATH_APK
	}

	PATH_OBB = os.Getenv("PATH_OBB")
	if PATH_OBB == "" {
		log.Fatalln("PATH_OBB not found")
	} else {
		listArgs["PATH_OBB"] = PATH_OBB
	}

	for i, k := range listArgs {

		if i == "OCULUS_APP_SECRET" {
			log.Println(i, " = ", "********")
		} else {
			log.Println(i, " = ", k)
		}

	}

}

func main() {

	http.HandleFunc("/upload", handelUpload)

	server := &http.Server{
		Addr: ":8080",
	}

	log.Println("<<SERVER START>>")
	log.Println(server.ListenAndServe())

}

func handelUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Выполнил работу по загрузке")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Errorf("read request body: %w", err))
	}

	fmt.Println(string(body))

	args := []string{
		"OCULUS",
		OCULUS_APP_ID,
		OCULUS_APP_SECRET,
		PATH_APK,
		PATH_OBB,
		"ALPHA",
	}

	cmd := exec.Command(PATH_UPLOADER_MOD, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Errorf("run uploader mod, path: %s, ERROR: %w", PATH_UPLOADER_MOD, err))
		echoToBot(false)
	}

	log.Println(string(output))
	echoToBot(true)
}

type Status struct {
	STATUS bool `json:"status"`
}

func echoToBot(answer bool) {
	tr := &http.Transport{}

	client := &http.Client{Transport: tr}

	jsonData, err := json.Marshal(Status{
		STATUS: answer,
	})
	if err != nil {
		log.Println(fmt.Errorf("marsahal: %w", err))
	}

	data := bytes.NewReader(jsonData)

	resp, err := client.Post("http://localhost:3030", "application/json", data)

	if err != nil {
		log.Println(fmt.Errorf("response request: %w", err))
	} else {
		log.Println(resp)
	}

}
