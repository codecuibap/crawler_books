# BOOK SITE CRAWLER DATA

Load list page then go to detail page that using query selector to get data detail

Format output by json

## Setting config by query selector (css) for each site.

By Descending pority in array

```json
    {
      "scrap_site": "https://nxbkimdong.com.vn/collections/all?page=1", // first site
      "site": "nxbkimdong.com.vn",  // domain
      "collection": "collections/all?page=",    //paging
      "url_detail": "nxbkimdong.com.vn/products",   //url patern detail
      "product": ".product-item > .product-img > a", // collection in product list page
      "section": "section[id=product-wrapper]", // detail card product
      "title": ["div.header_wishlist > h1"],
      "price": [".ProductPrice"], //number
      "page": [ //number
        "div.pro-short-desc>ul>li:nth-child(5)",
        "spen.field-name-field-product-sotrang span.field-items > span"
      ],
      "author": [
        "ul>li:nth-child(2) > a",
        "span.field-name-field-product-tacgia span.field-items > span"
      ],
      "isbn": [
        "div.pro-short-desc ul>li:nth-child(1) > strong",
        "div.pro-short-desc ul>li:nth-child(1)",
        "span.field-name-field-product--isbn span.field-items > span"
      ],
      "category": [
        "section[id=breadcrumb-wrapper] div.breadcrumb-content div.breadcrumb-small a:nth-child(2)"
      ],
      "name": [
        "section[id=breadcrumb-wrapper] div.breadcrumb-content div.breadcrumb-small a:nth-child(3)"
      ],
      "group": ["ul>li:nth-child(8) > a"],
      "desc": [],
      "rating": []
    },
```

## RUN

```shell
#BUILD
make build

#RUN
make crawl site="nxbkimdong.com.vn"
```

## CLEAR CACHE

Delete cache folder to clear cache

## TODO

| Site                         | Test    |
| ---------------------------- | ------- |
| http://www.nxbkimdong.com.vn | ✅ DONE |
| https://www.nxbtre.com.vn    |         |
| https://www.nxbhcm.com.vn    |         |
| https://www.nxbgd.vn         |         |
| https://nxbhoinhavan.vn      | ✅ DONE |
| http://www.nxblaodong.com.vn |         |
| https://nxbdantri.com.vn     |         |
| https://nxbphunu.com.vn      |         |
| https://dinhtibooks.com.vn   |         |
| https://nhanam.vn            | ✅ DONE |
| https://sachdonga.vn         |         |
| https://www.fahasa.com       |         |
| https://alphabooks.vn        | ✅ DONE |
| https://thaihabooks.com      |         |
| https://firstnews.com.vn/vi  | ✅ DONE |

## Version 2

> Fix performance && prety error

> Add image url, thumbnail list

> Category collection

> Fix Author & ISBN

> Add more format (csv, json)

> Add store provider (mongo, elaticsearch)
