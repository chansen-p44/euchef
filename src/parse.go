package euchef

import (
	"bytes"
	"log"
	"strings"

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
	sel := doc.Find(".product")

	menuItems := []MenuItem{}

	for i := range sel.Nodes {
		log.Printf("Menu %d", i+1)
		menuItem := MenuItem{Title: []string{}}

		menu := sel.Eq(i)

		date := menu.Find("h4 span")
		if date != nil {
			dateText := date.Text()
			log.Println("Menu date:", dateText)
			menuItem.Date = dateText
		}

		// Find the menu description text
		productDetails := menu.Find(".productDetails")
		if productDetails != nil {
			productDetailsItem := productDetails.First()
			if productDetailsItem != nil {
				paragraph := productDetailsItem.Find("p")
				if paragraph != nil {
					for i := range paragraph.Nodes {
						pNode := paragraph.Eq(i)
						text := pNode.Text()

						// Skip the silly "Maks 5 kuverter pr. person:-)" line
						if strings.Contains(text, "kuvert") && strings.Contains(text, "person") {
							continue
						}

						menuItem.Title = append(menuItem.Title, text)
						log.Println("Found menu title:", text)
					}
				} else {
					log.Println("productDetailsItem.Find(\"p\") is nil")
				}
			} else {
				log.Println("productDetailsItem.First() is nil")
			}
		} else {
			log.Println("menu.Find(\".productDetails\") is nil")
		}

		// Find the menu picture link
		productPhotoLink := menu.Find(".productPhotoLink")
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
	}

	log.Println("Parsing done")

	return menuItems, nil

}
