package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func jsonError(w http.ResponseWriter, text string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	text_error := fmt.Sprintf("%s [%s]", text, err)
	//fmt.Println(text_error)

	json_text, err := json.Marshal(map[string]string{"error": text_error})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Ошибка генерации JSON для jsonError:", err)
		return
	}

	_, err = w.Write(json_text)
	if err != nil {
		fmt.Println("Ошибка записи JSON:", err)
		return
	}
}
