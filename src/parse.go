package euchef

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// MenuItem describes a single menu item
type MenuItem struct {
	Date     string
	Title    []string
	ImageURL string
}

func ParseData(data []byte) ([]MenuItem, error) {

	buf := bytes.NewBuffer(data)
	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		log.Println("NewDocumentFromReader error:", err)
		return nil, err
	}

	log.Println("Parsing")

	// Find the menus
	sel := doc.Find(".category-products")
	sel = sel.First()
	sel = sel.Find(".w-saledate")
	menuItems := []MenuItem{}

	now := time.Now()
	currentDateString := Format(now)
	log.Println("Current date:", currentDateString)

	for i := range sel.Nodes {
		log.Printf("Menu %d", i+1)
		menuItem := MenuItem{Title: []string{}}

		menu := sel.Eq(i)

		// Find the date this is valid
		date := menu.Find(".salgstid")
		if date != nil {
			dateText := date.Text()
			dateText = strings.TrimSpace(dateText)
			log.Println("Menu date:", dateText)

			if currentDateString != dateText {
				log.Println("Menu date doesn't match current date, skipping")
				continue
			}

			menuItem.Date = dateText
		}

		// Find the menu description text
		productDetails := menu.Find(".desc")
		if productDetails != nil {
			text := productDetails.Text()
			text = strings.TrimSpace(text)
			text = strings.ReplaceAll(text, "\n", "")
			log.Println("Found menu title:", text)
			menuItem.Title = append(menuItem.Title, text)
		} else {
			log.Println("menu.Find(\".productDetails\") is nil")
		}

		// Find the menu picture link
		productPhotoLink := menu.Find(".product-image")
		if productPhotoLink != nil {
			link := productPhotoLink.First()
			if link != nil {
				imgElement := link.Find("img")
				if imgElement != nil {
					imgNode := imgElement.Eq(0)
					src, found := imgNode.Attr("src")
					if found {
						menuItem.ImageURL = src
						log.Println("Found menu image url:", src)
					} else {
						log.Println(`imgNode.Attr("src") not found`)
					}
				} else {
					log.Println(`link.Find("img") is nil`)
				}
			} else {
				log.Println(`productPhotoLink.First() is nil`)
			}
		} else {
			log.Println(`menu.Find(".productPhotoLink") is nil`)
		}
		menuItems = append(menuItems, menuItem)

		// Only take the first menu
		//break
	}

	log.Println("Parsing done")

	return menuItems, nil

}

func Format(t time.Time) string {
	return fmt.Sprintf("%s %02d. %s",
		days[t.Weekday()], t.Day(), months[t.Month()-1],
	)
}

var days = [...]string{
	"søndag", "mandag", "tirsdag", "onsdag", "torsdag", "fredag", "lørdag"}

var months = [...]string{
	"januar", "feburar", "marts", "april", "maj", "juni",
	"juli", "august", "september", "oktober", "november", "december",
}
