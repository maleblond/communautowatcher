#!/bin/sh

curl 'https://www.reservauto.net/Scripts/Client/Ajax/PublicCall/Get_Car_DisponibilityJSON.asp' \
  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36' \
  --data-raw 'CurrentLanguageID=1&CityID=90&StartDate=01%2F10%2F2022+11%3A00&EndDate=01%2F10%2F2022+11%3A30&Accessories=0&FeeType=80&Latitude=46.8046123&Longitude=-71.2342123'