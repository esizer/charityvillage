package main

import (
	"encoding/csv"
	"fmt"
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

func postScrape(category string, url string) error {

	charities := []charity{}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
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
	writeRecords(charities, category)

	return nil
}

func writeRecords(charities []charity, category string) error {
	headers := []string{"Charity", "City"}

	file, err := os.Create(fmt.Sprintf("%s_charities.csv", category))
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

func main() {
	categories := map[string]string{
		"addictions":            "https://charityvillage.com/cms/organizations/addictions-and-substance-abuse",
		"children-youth-family": "https://charityvillage.com/cms/organizations/children-youth-and-family",
	}
	for cat, url := range categories {
		postScrape(cat, url)
	}
}
