package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type StrategyParse struct {
	ScrapSite   string `json:"scrap_site"` //https://nxbkimdong.com.vn/collections/all?page=1
	Site        string `json:"site"`       //nxbkimdong.com.vn
	Collection  string `json:"collection"` //collections/all?page=
	UrlDetail   string `json:"url_detail"` //nxbkimdong.com.vn/products
	Product     string `json:"product"`    //.product-item > .product-img > a
	Section     string `json:"section"`    //section[id=product-wrapper]
	Title       string `json:"title"`      //div.header_wishlist > h1
	Price       string `json:"price"`      //.ProductPrice
	Page        string `json:"page"`       //ul>li:nth-child(5)
	Author      string `json:"author"`     //ul>li:nth-child(2) > a
	ISBN        string `json:"isbn"`       //ul>li:nth-child(1) > strong
	Category    string `json:"category"`   //"ul>li:nth-child(8) > a
	Name        string `json:"name"`       //"ul>li:nth-child(8) > a
	Group       string `json:"group"`      //"ul>li:nth-child(8) > a
	Description string `json:"desc"`       //
	Rating      string `json:"rating"`
}

type Strategy struct {
	Strategies []StrategyParse `json:"strategies"`
}

func LoadConfig(key string) (*StrategyParse, error) {
	jsonFile, err := os.Open("strategy.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened strategy.json")
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var result Strategy
	json.Unmarshal([]byte(byteValue), &result)
	for _, s := range result.Strategies {
		if s.Site == key {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("Can not find key in config [%s]", key)
}
