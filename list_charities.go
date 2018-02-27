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
				"International:", "",
				"National:", "",
				"Regional:", "",
				"Local:", "",
				"City:", "",
				"ON", "",
				",", "",
				")", "",
			)
			if len(region) > 0 {
				city := r.Replace(region[len(region)-1][0])
				name := strings.Replace(title, region[len(region)-1][0], "", 1)
				charities = append(charities, charity{Name: name, City: city})
			}

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
		"addictions":                "https://charityvillage.com/cms/organizations/addictions-and-substance-abuse",
		"children-youth-family":     "https://charityvillage.com/cms/organizations/children-youth-and-family",
		"community-social-services": "https://charityvillage.com/cms/organizations/community-and-social-services",
		"criminal-justice":          "https://charityvillage.com/cms/organizations/criminal-justice",
		"culture-heritage":          "https://charityvillage.com/cms/organizations/culture-and-heritage",
		"disabilities":              "https://charityvillage.com/cms/organizations/disabilities",
		"lgbtq":                     "https://charityvillage.com/cms/organizations/lesbian-gay-bisexual-transgender-lgbt",
		"health-diseases":           "https://charityvillage.com/cms/organizations/health-and-diseases",
		"human-rights":              "https://charityvillage.com/cms/organizations/human-rights-and-civil-liberties",
		"itl-relief":                "https://charityvillage.com/cms/organizations/international-relief-development-peace",
		"poverty":                   "https://charityvillage.com/cms/organizations/poverty-social-justice",
		"public":                    "https://charityvillage.com/cms/organizations/public-society-benefit",
		"senior-citizens":           "https://charityvillage.com/cms/organizations/senior-citizens",
		"sports-recreation":         "https://charityvillage.com/cms/organizations/sports-and-recreation",
		"women":                     "https://charityvillage.com/cms/organizations/women",
	}
	for cat, url := range categories {
		charities = postScrape(cat, url, charities)
	}

	charities = filterCharities(charities)

	writeRecords(charities)

}
