package euchef

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func FetchData(when time.Time) ([]byte, error) {

	formData := []string{
		"PageSize=1000",
		"PageNumber=1",
		"SearchMenus=True",
		"SearchProducts=False",
		"SearchCoffeeBreakTimeSlots=false",
		"deliveryText=15-01-2019",
		"searchText=",
		"performanceTimeslotId=",
		"amountOfCustomer=",
	}

	formData[5] = "deliveryText=" + when.Format("02-01-2006")

	postData := bytes.NewBufferString(strings.Join(formData, "&"))

	log.Println("Requesting data for date:", formData[5])

	req, err := http.NewRequest("POST", "https://www.wipaway.dk/5137", postData)
	if err != nil {
		log.Println("NewRequest error: ", err)
		return nil, err
	}

	req.Header.Add("content-type", `application/x-www-form-urlencoded; charset=UTF-8`)
	req.Header.Add("origin", `https://www.wipaway.dk`)
	req.Header.Add("referer", ` https://www.wipaway.dk/5137`)

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

	//ioutil.WriteFile("menu_data.html", buf, 0644)

	return buf, nil
}
