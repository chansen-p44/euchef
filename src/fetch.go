package euchef

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func FetchData() ([]byte, error) {

	req, err := http.NewRequest("GET", "https://rtx.takeaway.cheval-blanc.dk/", nil)
	if err != nil {
		log.Println("NewRequest error: ", err)
		return nil, err
	}

	client := &http.Client{Timeout: time.Second * 5}
	res, err := client.Do(req)
	if err != nil {
		log.Println("client Do error: ", err)
		return nil, err
	}

	log.Printf("Request statuscode: %d", res.StatusCode)

	// Make sure we close the body stream when we are done
	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error: ", err)
		return nil, err
	}

	log.Printf("Downloaded %d bytes", len(buf))

	ioutil.WriteFile("menu_data.html", buf, 0644)

	return buf, nil
}
