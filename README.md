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
   
---
## 執行該專案方法
本專案使用 Golang 的 Gin 開發 RESTful API，使用PostgreSQL作為後端資料庫，以Redis做為快取資料庫。

### Use Docker
使用 Docker 可以直接運行該專案會用到的 ` golang, postgres, pg-admin, redis, redis-admin ` ，執行該專案請確認測試前有安裝 **Git, Docker 以及 Docker compose**
1. ```sh
   #透過github clone將此專案下載至本機

   git clone https://github.com/WeiWeiCheng123/URLshortener.git 
   ```
2. ```sh
   #在專案資料夾底下，透過 shell script 來執行此專案
   #shell script 內是 docker-compose up

   cd URLshortener
   ./start-project.sh
   ```
   在docker compose的部分啟動了以下
   - postgres 
  
      Postgres 為此專案的 backend database

   - postgres-admin
    
      使用 pgadmin 可以直接查看目前 Postgres 內的資訊，透過開啟 **localhost:81**，並在帳號的地方輸入 **dcard123@dcard.com**，密碼輸入 **dcard123**，登入後再連線至 postgres 資料庫即可。

   - redis 
  
      Redis 為此專案的 Cache

   - redis-admin
  
      使用 phpredisadmin 可以查看目前 Redis 內的資訊，透過開啟 **localhost:82**，即可查看目前 redis 的資料。

   - backend
  
      使用Golang撰寫的Bankend server，透Dockerfile進行包裝，連線的網址為**localhost:8080**，根據不同的API，後面的網址會有不同

---
## 測試方法
在這個專案底下，設計了兩個API

- 第一個API為POST，給定一個網址和過期時間，將原網址變成縮網址
  
  - Example request
    ```sh
    curl -X POST -H "Content-Type:application/json" -d '{"url":"https://www.dcard.tw/f","expireAt":"2023-01-01T09:00:41Z"}' http://localhost:8083/api/v1/urls
    ```
   
  - Example response
    ```sh
    {"id":"KF4eAy9","shortURL":"http://localhost:8083/KF4eAy9"}
    ```

- 第二個API為GET，若該縮網址沒過期，輸入縮網址，重新導向到原網址; 若縮網址過期，則顯示404
  
    - Example request
      ```sh
      # use the return shortID
      # http://localhost:8080/ + shortID
      curl -L -X GET "http://localhost:8083/KF4eAy9"
      ```

  - Example response
      ```
      該網站HTML form
      ```

---
## 專案目錄結構
```c
|--URLshortener 
    |-- .env                 // 環境變數，裡面宣告了資料庫連線的資料和程式中會使用到的參數
    |-- docker-compose.yaml  // Docker compose file，用來執行此專案
    |-- Dockerfile           // Dockerfile，用來包裝API
    |-- go.mod
    |-- go.sum
    |-- main.go              // 主程式
    |-- main_parse_test.go   // 測試主程式的 GET API
    |-- main_shorten_test.go // 測試主程式的 POST API  
    |-- README.md             
    |-- start-project.sh     // 用來啟動此專案的 shell script
    |-- docker-pg-init       // 用來 Initialize postgres，裡面包含創建使用者、創建DB及Table
    |   |-- init.sql
    |-- handler              // API 的功能
    |   |-- handler.go
    |-- lib
    |   |-- config           // 取得 .env 檔的環境變數
    |   |   |-- config.go
    |   |-- cron             // 排程任務，用於刪除Postgres的過期資料
    |   |   |-- cron.go
    |   |-- function         // 檢查輸入格式及產生shortID
    |   |   |-- checker.go
    |   |   |-- checker_test.go
    |   |   |-- id_generator.go
    |   |-- lua              // 用於IP limit function，確保在維持原子性的情況下操作Redis
    |   |   |-- lua.go
    |   |-- middleware       // 中間層，用來處理資料及限制IP不得在規定時間內超過最大值
    |       |-- middleware.go
    |-- model                // 宣告 short URL structure
        |-- model.go

```

---
## 設計發想
### 縮網址設計
- 縮網址的長度

  - 在取決縮網址長度的部分，我首先先查看了 http://www.worldwidewebsize.com，這個網站上會顯示目前全世界網站的數量有多少，假設網站的數字為真的話，統計至 (2022/02) 的數量是 1.84 億個網站
  
  - 知道目前數量後，就需要來決定縮網址的長度了，首先縮網址的內容我設定只有數字 0~9 和英文的大小寫，一共 62 個字母，$62^N\geq1.84*10^8\implies N\geq4.61$
  
  - 經過計算之後，可以得知如果要 Cover 到全世界的網站的話縮網址的長度至少要有 5 ，又怕因為長度只有 5 太少，導致重複的機率變高，因此我將縮網址的長度定成7
  
- 產生縮網址的方式 
  
  - Random string
  
    其實最直觀的方式就是隨機產生出一組的英文 + 數字的組合
    
    優點:

    - 容易實現，而且可以預先產生一大堆的 random string
    - 執行速度很快，可以降低反應的時間    
    
    缺點:

    - 會有機會產生出重複的 random string

  - Hash function
  
    Hash function 會產生出唯一的 value ，再使用這個 value 轉換成 base62 的格式，在依照縮網址長度的需求去取需要的長度

    優點:

    - Hash function 可以產生出唯一的 value

    缺點:

    - Hash function 在實現上如果只使用原網址進行 hash ，每一個人 hash 的結果都會一模一樣，因此 hash function 的 input 還需要再加上其他的數值來避免衝突，在實現上會比較麻煩
    - Hash function 再產生 value 後，假如只取前/後面的幾個數值當縮網址，有機會發生重複的情況
    - Hash function 比起 Random string 來說，速度較慢
  
  - 此專案的方式
  
    在本專案中實現縮網址的方式流程如下

    1. 將現在的時間轉成UnixNano去產生random seed `s1` 
    2. 利用 `s1` 產生出一個隨機數字 `r1` 
    3. 再用現在時間的UnixNano加上剛剛產生的隨機數 `r1` 去產生的random seed `s2` 
    4. 利用 `s2` 產生出一個隨機數字 `r2` 
    5. 將現在時間的Unix加上 `r1` 和 `r2` 
    6. 判斷該產生的數字是否為偶數，如不是偶數則加 1
    7. 將處理完的數字轉換成 base62 後取最後7位元並回傳

    引入當前的時間當作產生縮網址的參數，可以避免去產生到重複的縮網址，且只產生偶數數字來當作輸入的檢查機制

### Cache
在設計Cache的時候，常見的問題有

1. Caching Penetration
2. Cache avalanche
3. Cache Stampede 

在這個專案中，容易發生問題的地方是縮網址 Redirect 到原網址的 API

假設現在有一群人要使用原先不在 Cache 的縮網址或是一大堆根本不存在的縮網址的request時該怎麼做

解決方法: 

- 在縮網址進入後端前，先進行長度的檢查，查看長度是否為 7
- 檢查此縮網址經過解碼後是否為偶數
- 在GET的這個地方加入 IP Limit 的方式去限制同一個IP在短時間內大量的 request
- 當使用者輸入一個不存在的縮網址，會在 Redis 存入該縮網址，並將 Value 設為 **NotExist** ，過期時間設置為 150 秒，當使用者想在短時間內重複輸入該不存在的縮網址就不會每次都去後端資料庫找資料

假設使用者輸入了一個正確並存在於後端資料庫的縮網址 ID 時，程式會將此縮網址 ID 存入 Redis 並判斷過期時間，假如該縮網址的過期時間大於 20分鐘 的話就將該縮網址存入 Redis 的 TTL 設為 20分鐘，避免出現該縮網址過期時間還有一年所以這筆資料就會放在 Cache 裡一年的時間，反之小於 20 分鐘則直接存入 Redis

### 資料庫Table
當資料庫內的資料量極大時需要去考慮到如何設計去進行優化，讓搜尋速度可以往上提升

本專案使用的是 Postgres ，設計 Table 時是使用 URL shortID 解碼後的數字當作 Primary key(PK) ，在 Postgres 內部會幫我們自動把 PK 當作 index ，當資料庫內資料量很多時，這麼做可以加快在 GET 的這個 method 從資料庫搜尋出資料的速度。

### IP Limit
由於在這個專案中需要去考慮到有人會在短時間多次的透過短網址重新導向到原網址，因此我這邊設計了一個 IP limit 去限制同一個IP下在一段時間內的 Request 數量，預設是在 300 秒內的 Request 數量為 2000 ，當超過限制後就會回傳 429

又因為怕產生 Race condition ，為確保原子性，因此使用 lua 腳本來進行操作

lua 腳本如下

```lua

local ip = KEYS[1]
local ipLimit = tonumber(ARGV[1])
local period = tonumber(ARGV[2])
local count = redis.call('GET', ip)
if (not count) then
    redis.call('SET', ip, 1)
    redis.call('EXPIRE', ip, period)
    return 1
end
if (tonumber(count) < ipLimit) then
    redis.call('INCR', ip)
    redis.call('EXPIRE', ip, period)
    return 1
end
redis.call('EXPIRE', ip, period)
return -1

```
說明 : 
- **KEYS[1], ARGV[1], ARGV[2]** 分別是使用者的 IP 、 IP 流量的最大值和限制 IP 的時間
1. 在最一開始會先查看 Redis 是否有這個 IP 的紀錄
     - 如果沒有，則 SET 該 IP 的數值為 1 ，TTL 為限制時間，並回傳 1 
    - 如果有則查看該 IP 內的數值
2. 如果 IP 內數值小於最大值，則將數值增加 1 ，刷新過期時間，並回傳 1
3. 如果 IP 內數值等於最大值則回傳 -1 ，並刷新過期時間

- lua 執行完畢後會回傳 1 or -1
  - 1 代表該 IP 未超過限制
  - -1 代表該 IP 超過限制

### Cron Job
在這個專案中我使用 Postgres 用來存放 縮網址 ID 、原網址、過期時間 ，由於資料具有時間性，因此我使用Cron Job的方式來刪除過期的資料，在 Demo 的時候我預設是每 5 分鐘進行刪除，假如在真正的應用上，我會選擇設定在冷門時段，像是凌晨 3 點

Cron Job 設定時間的格式如下

		*     *     *      *      *
		分    時    日     月     星期
		0-59  0-23  1-31  1-12  0-6 (週日~週六)
For example
```
   */5 * * * * == 每當分鐘數為 5 的倍數時執行一次 (0, 5, 10, 15, 20...)
    0  3 * * * == 凌晨 3 點整執行
```

---
## 錯誤處理
在設計這個專案的時候我考慮了以下這些項目

1. POST API 中輸入不符合格式的URL
  
   使用 `net/url` 的 Parse來檢查該URL的結構是否有錯誤，正常的 URL 格式應該要為 `http://... or https://...` ，假如格式不符，回傳 400

2. POST API 中輸入錯誤的時間格式或是已經過期的時間
   
   使用 `Time` 來檢查時間格式是否正確，以及去比對現在時間，確認有沒有過期，假如格式不符或是過期，回傳400

3. GET API 中嘗試輸入錯誤的縮網址
   
   - 先確認該縮網址長度是否為 7，假如不符合，回傳 404
   - 確認該縮網址經過解碼之後是否為2的倍數，假如不是則直接回傳 404
   - 假如該縮網址不存在 Cache 以及 Database 中，則在 Cache 存入該縮網址，給予 Value  **NotExist** 並設定 150 秒後過期，假如有惡意的想要在短時間內多次的連線該不存在的網址就可以避免再次到後端的 Database 中查找，可以在 Cache 直接找到並回傳 404

4. 某一個IP嘗試在短時間內大量連線
   
   在 middleware 中使用 IP limit 的方式去計算該 IP 的 Request 次數，使用了 lua 搭配 Redis 來達成 IP limit 的功能，當達到上限時，回傳 429

---
## Test
在這個專案中分成兩個測試，一個是將原網址變成縮網址，一個是重新導向縮網址

在原網址變成縮網址的測試裡面又分成
1. 輸入正確無誤的格式，回傳 200 以及 short URL 的內容
2. 輸入錯誤的網址格式，回傳 400 以及 invalid URL
3. 輸入錯誤的時間格式，回傳 400 以及 error time format or time is expired
4. 輸入過期的時間格式，回傳 400 以及 error time format or time is expired

在重新導向縮網址的測試裡面又分成
1. 輸入實際存在的縮網址，回傳 302
2. 輸入不存在的縮網址，回傳 404 以及 this shortid is not existed or expired
3. 輸入過期的縮網址，回傳 404 以及 this shortid is not existed or expired
    
---
## ApacheBench
設計完成後我有使用 ApacheBench 進行測試，在我的 Blog 有介紹到 ApacheBench https://weiweicheng123.github.io/2022/02/26/ab-test/#more

測試傳入API的總數量為 1000 ，並模擬 200 個 Request 同時進行

```sh
#POST

ab -n 1000 -c 200 -p post.json -T 'application/x-www-form-urlencoded' http://localhost:8083/api/v1/urls

Concurrency Level:      200
Time taken for tests:   0.841 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      182000 bytes
Total body sent:        235000
HTML transferred:       59000 bytes
Requests per second:    1188.66 [#/sec] (mean)
Time per request:       168.257 [ms] (mean)
Time per request:       0.841 [ms] (mean, across all concurrent requests)
Transfer rate:          211.26 [Kbytes/sec] received
                        272.79 kb/s sent
                        484.05 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    3   4.6      1      32
Processing:     9  150 117.1    118     790
Waiting:        3  148 117.2    118     790
Total:          9  153 116.9    121     801

Percentage of the requests served within a certain time (ms)
  50%    121
  66%    155
  75%    183
  80%    207
  90%    285
  95%    372
  98%    546
  99%    669
 100%    801 (longest request)

```

```sh
#GET

ab -n 1000 -c 200 http://localhost:8083/L4qOJ2a

Concurrency Level:      200
Time taken for tests:   0.441 seconds
Complete requests:      1000
Failed requests:        0
Non-2xx responses:      1000
Total transferred:      198000 bytes
HTML transferred:       45000 bytes
Requests per second:    2268.50 [#/sec] (mean)
Time per request:       88.164 [ms] (mean)
Time per request:       0.441 [ms] (mean, across all concurrent requests)
Transfer rate:          438.64 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    9   8.2      8      34
Processing:     3   76  63.4     54     213
Waiting:        2   70  62.7     50     211
Total:          4   85  63.0     63     223

Percentage of the requests served within a certain time (ms)
  50%     63
  66%     78
  75%     97
  80%    142
  90%    203
  95%    217
  98%    220
  99%    220
 100%    223 (longest request)

```