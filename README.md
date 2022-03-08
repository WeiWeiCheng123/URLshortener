# Dcard Backend Intern Homework
## 介紹
這個專案是Dcard在2022年的Backend Intern作業，作業的目標是設計並使用Golang做出一個URL Shortener的API

專案的要求如下

1. URL shortener has 2 APIs, please follow API example to implement:
   1. A RESTful API to upload a URL with its expired date and response with a shorten URL.
   
   2. An API to serve shorten URLs responded by upload API, and redirect to original URL. If URL is expired, please response with status 404.

2. Please feel free to use any external libs if needed.

3. It is also free to use following external storage including:

   1. Relational database (MySQL, PostgreSQL, SQLite)

   2. Cache storage (Redis, Memcached)

4. Please implement reasonable constrains and error handling of these 2 APIs.

5. You do not need to consider auth.

6. Many clients might access shorten URL simultaneously or try to access with non-existent shorten URL, please take performance into account.
----
## 執行該專案方法
### Use Docker
使用 docker 可以直接運行該專案會用到的golang, postgres, redis，請確認測試前有安裝docker
1. ``` 
   #透過github clone將專案下載至本機
   git clone https://github.com/WeiWeiCheng123/URLshortener.git 
   ```
2. ```
   #在專案資料夾底下，透過docker來執行此專案
   docker-compose up -d
   ```
   在docker compose的部分啟動了以下
   1. Postgres ( 作為此專案的backend database )
   2. Redis ( 作為此專案的Cache )
   3. 此專案 ( 使用Dockerfile進行包裝 )
----
## 測試方法
在這個專案底下，設計了兩個API

- 第一個API為POST，給定過期時間並將原網址變成縮網址
  - Example request
    ```
    curl -X POST -H "Content-Type:application/json" -d '{"url":"https://www.dcard.tw/f","expireAt":"2023-01-01T09:00:41Z"}' http://localhost:8080/api/v1/urls
    ```
  - Example response
    ```
    {"id":"m42er1M","shortURL":"http://localhost:8080/m42er1M"}
    ```
- 第二個API為GET，若該縮網址沒過期，輸入縮網址，重新導向到原網址; 若縮網址過期，則顯示404
    - Example request
      ```
      curl -L -X GET "http://localhost:8080/m42er1M"
      ```
  - Example response
      ```
      該網站html from
      ```
----
## 設計想法
### 縮網址設計
- 縮網址的長度 ?

  - 在取決縮網址長度的部分，我首先先查看了 http://www.worldwidewebsize.com，這個網站上會顯示目前全世界網站的數量有多少，假設網站的數字為真的話，目前(2022/02/17)的數量是1.84億個網站
  
  - 知道目前數量後，就需要來決定縮網址的長度了，首先縮網址的內容我設定只有數字0~9和英文的大小寫，一共62個字母，$62^N\geq1.84*10^8\implies N\geq4.61$
  
  - 經過計算之後，可以得知如果要Cover到全世界的網站的話縮網址的長度至少要有5，又怕因為長度只有5太少，導致重複的機率變高，因此我將縮網址的長度定成7
  
- 產生縮網址的方式 
  
  - Random string
  
    其實最直觀的方式就是隨機產生出一組的英文+數字的組合
    
    優點:

    - 容易實現，而且可以預先產生一大堆的random string
    - 執行速度很快，可以降低反應的時間    
    
    缺點:  
    - 會有機會產生出重複的random string

  - Hash function
  
    Hash function 會產生出唯一的value，再使用這個value轉換成base62的格式，在依照縮網址長度的需求去取需要的長度

    優點:
    - Hash function 可以產生出唯一的value

    缺點:
    - Hash function在實現上如果只使用原網址進行hash，每一個人hash的結果都會一模一樣，因此hash function的input還需要再加上其他的數值來避免衝突，在實現上會比較麻煩
    - Hash function再產生value後，假如只取前/後面的幾個數值當縮網址，有機會發生重複的情況
    - Hash function比起Random string來說，速度較慢
  
  - 我的方式
  
    在本專案中實現縮網址的方式流程如下
    1. 將現在的時間轉成UnixNano去產生random seed `s1` 
    2. 利用 `s1` 產生出一個隨機數字 `r1` 
    3. 再用現在時間的UnixNano加上剛剛產生的隨機數 `r1` 去產生的random seed `s2` 
    4. 利用 `s2` 產生出一個隨機數字 `r2` 
    5. 將現在時間的Unix加上 `s1` 和 `s2` 之後將其轉為base62

    引入當前的時間當作產生縮網址的參數，可以避免去產生到重複的縮網址

### Cache
在設計Cache的時候，常見的問題有

1. Caching Penetration
2. Cache avalanche
3. Cache Stampede 

在這個專案中，容易發生問題的地方是縮網址Redirect到原網址的API

假設現在有一群人要使用原先不在Cache的縮網址或是一大堆根本不存在的縮網址的request時該怎麼做 ?

解決方法: 

- 在縮網址進入後端前，先進行長度的檢查，查看長度是否為7
- 在GET的這個地方加入IP Limit的方式去限制同一個IP在短時間內大量的request
- 當使用者輸入一個不存在的縮網址，會在Redis存入該縮網址，並將Value設為NotExist，當使用者想再次輸入該不存在的縮網址不會每次都跟後端資料庫找資料
- 當Cache內沒有該縮網址時，取得一個Lock以確保不會有大量的request同時間進入後端資料庫內，當結束後端的查找工作之後再Unlock

### IP Limit
由於在這個專案中需要去考慮到有人會在短時間多次的透過短網址重新導向到原網址，因此我這邊設計了一個IP limit 去限制同一個IP下在一段時間內的Request數量，預設是在300秒內的Request數量為500，當超過限制後就會回傳429

### Cron Job
在這個專案中我使用Postgres用來存放 縮網址ID、原網址、過期時間 ，由於資料具有時間性，因此我使用Cron Job的方式來刪除過期的資料，在Demo的時候我預設是每3分鐘進行刪除，假如在真正的應用上，我會選擇設定在冷門時段，像是凌晨3點

----
## 錯誤處理
在設計這個專案的時候我考慮了以下這些項目
1. POST API 中輸入不符合格式的URL
  
   使用 `net/url` 的 Parse來檢查該URL的結構是否有錯誤，正常的URL格式應該要為 `http://<domain> or https://<domain>` ，假如格式不符，回傳400

2. POST API 中輸入錯誤格式或是已經過期的時間
   
   使用 `Time` 來檢查時間格式是否正確，以及去比對現在時間，確認有沒有過期，假如格式不符或是過期，回傳400

3. GET API 中嘗試輸入錯誤的縮網址
   
   - 確認該縮網址長度是否為7，假如不符合，回傳404
   - 假如該縮網址不存在Cache以及Database中，則在Cache存入該縮網址，給予Value ` NotExist ` 並設定300秒後過期，假如下次再輸入一樣的縮網址就可以避免再次到後端的Database中查找，可以在Cache直接找到並回傳404

4. 某一個IP嘗試在短時間內大量連線
   
   在middleware中使用IP limit的方式去計算該IP的Request次數，使用了lua搭配Redis來達成IP limit的功能，當達到上限時，回傳429
----
## 系統架構