package main

import (
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type charity struct {
	Name string
	City string
}

func postScrape(category string, url string, charities []charity) []charity {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Panic(err)
	}

	doc.Find(".DNNModuleContent .Normal h4").Each(func(index int, item *goquery.Selection) {
		title := item.Text()
		re := regexp.MustCompile(`\(.*?\)`)
		if strings.Contains(title, "ON") {
			region := re.FindAllStringSubmatch(title, -1)
			r := strings.NewReplacer(
				"(", "",
				"National:", "",
				"Regional:", "",
				"Local:", "",
				"City:", "",
				"ON", "",
				",", "",
				")", "",
			)
			city := r.Replace(region[len(region)-1][0])
			name := strings.Replace(title, region[len(region)-1][0], "", 1)

			charities = append(charities, charity{Name: name, City: city})
		}
	})

	return charities
}

func writeRecords(charities []charity) error {
	headers := []string{"Charity", "City"}

	file, err := os.Create("charities.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write(headers)

	for _, c := range charities {
		writer.Write([]string{c.Name, c.City})
	}
	return nil
}

func filterCharities(charities []charity) []charity {
	u := make([]charity, 0, len(charities))
	m := make(map[charity]bool)

	for _, val := range charities {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

func main() {

	charities := []charity{}

	categories := map[string]string{
		"addictions":            "https://charityvillage.com/cms/organizations/addictions-and-substance-abuse",
		"children-youth-family": "https://charityvillage.com/cms/organizations/children-youth-and-family",
	}
	for cat, url := range categories {
		charities = postScrape(cat, url, charities)
	}

	charities = filterCharities(charities)

	writeRecords(charities)

}
