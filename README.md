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

## 設計想法

### 縮網址設計

- 縮網址的長度 ?

  - 在取決縮網址長度的部分，我首先先查看了 http://www.worldwidewebsize.com，這個網站上會顯示目前全世界網站的數量有多少，假設網站的數字為真的話，目前(2022/02/17)的數量是1.84億個網站
  
  - 知道目前數量後，就需要來決定縮網址的長度了，首先縮網址的內容我設定只有數字0~9和英文的大小寫，一共62個字母，**接下來就是可怕的數學了**，$62^N\geq1.84*10^8\implies N\geq4.61$
  
  - 經過可怕的運算之後(計算機救我)，可以得知如果要Cover到全世界的網站的話縮網址的長度至少要有5
  
- 產生縮網址的方式 ?
  
  - Random string
  
    其實最直觀的方式就是隨機產生出一組的英文+數字的組合
    
    優點:

    - 容易實現，而且可以預先產生一大堆的random string
    - 執行速度很快，可以降低反應的時間    
    
    缺點:  
    - 會有機會產生出重複的random string (假想一個情況，有一個是Youtube的縮網址，一個是Dcard的縮網址，但很不幸的縮網址不小心重複了，那請問這個人連進這個縮網址是要看影片還是看文章呢XD)，但其實可以預先產生一大堆random string來進行檢查，就可以避免掉產生重複縮網址的情況。

  - Hash function
    Hash function 會產生出唯一的value，再使用這個value轉換成base62的格式，在依照縮網址長度的需求去取需要的長度

    優點:
    - Hash function 可以產生出唯一的value

    缺點:
    - Hash function在實現上如果只使用原網址進行hash，每一個人hash的結果都會一模一樣，因此hash function的input還需要再加上其他的數值來避免衝突，在實現上會比較麻煩
    - Hash function再產生value後，假如只取前/後面的幾個數值當縮網址，有機會發生重複的情況
    - Hash function比起Random string來說，速度較慢

**結論**

在本專案中實現縮網址的方式是採用，待補

### Cache
在設計Cache的時候，常見的問題有

1. Caching Penetration
2. Cache avalanche
3. Cache Stampede 

在這個專案中，容易發生問題的地方是縮網址Redirect到原網址的API

假設現在有一群人要使用原先不在Cache的縮網址或是一大堆根本不存在的縮網址的request時該怎麼做 ?

解決方法: 

- 123
- 123
- 123
- 123
- 123
- 123

### 系統