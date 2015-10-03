package cmds

import (
	"appengine"
	"appengine/urlfetch"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
	"strings"
)

func Payup(c appengine.Context, p []string, w io.Writer) {
	link := p[0]
	link = link[1 : len(link)-1]

	client := urlfetch.Client(c)
	resp, err := client.Get(link)
	if err != nil {
		fmt.Fprintf(w, "Error getting URL: %v %v", link, err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	prices := make([]float64, 0)
	names := make([]string, 0)

	doc.Find("#order-confirmation-page").Find("tbody").Find("tr").Find(".price-table").Each(func(i int, s *goquery.Selection) {
		p := strings.TrimSpace(s.Text())
		if price, err := strconv.ParseFloat(p[1:], 64); err == nil {
			prices = append(prices, price)
		}
	})

	doc.Find("#order-confirmation-page").Find("tbody").Find("tr").Find("strong").Each(func(i int, s *goquery.Selection) {
		names = append(names, strings.TrimSpace(s.Text()))
	})

	fees := float64(0)

	doc.Find(".order-information").Find(".column-2:nth-child(2)").Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		item := strings.TrimSpace(tds.First().Text())
		p := strings.Replace(strings.TrimSpace(tds.Last().Text()), "$", "", -1)
		price, _ := strconv.ParseFloat(p, 64)

		if strings.Contains(item, "Fee") || strings.Contains(item, "Tip") || strings.Contains(item, "Tax") {
			fees = fees + price
		}
	})

	due := make(map[string]float64)

	for i := range names {
		due[names[i]] = due[names[i]] + prices[i]
	}

	// for key, value := range due {
	// 	due[key] = value + fees/float64(len(due))
	// }

	feePerPerson := fees / float64(len(due))

	fmt.Fprintf(w, "Well then, here we go. Fees, taxes, and tip evenly divided across labels:\n\n")
	for key, value := range due {
		fmt.Fprintf(w, "*%v*: $%.2f + $%.2f = *$%.2f*\n", key, value, feePerPerson, value+feePerPerson)
	}
	fmt.Fprintf(w, "\n_Total Taxes, Fees, and Tip: $%.2f, Per label: $%.2f_", fees, feePerPerson)
	// fmt.Fprintf(w, "```")
}
