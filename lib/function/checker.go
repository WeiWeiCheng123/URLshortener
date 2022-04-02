package function

import (
	"errors"
	"net/url"
	"time"
)

//valid the URL, if url not correct then return false
func IsURL(OriginalURL string) bool {
	u, err := url.Parse(OriginalURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

//time formater, if time format is wrong or expired, return false.
//otherwise, return true and convert to time type
func TimeFormater(expTime string) (time.Time, error) {
	var localLocation *time.Location
	localLocation, _ = time.LoadLocation("Asia/Shanghai")
	layout := "2006-01-02T15:04:05Z"
	expireTime, err := time.Parse(layout, expTime)
	//Time format error
	if err != nil {
		return time.Time{}, err
	}

	expireTime = expireTime.In(localLocation)
	expireTime = expireTime.Add(-8 * time.Hour)
	//Time expired
	if time.Now().After(expireTime) {
		return time.Time{}, errors.New("data expired")
	}

	return expireTime, nil
}

//change the time to Taiwan time zone
func Time_to_Taiwanzone(expTime time.Time) (time.Time, error) {
	var localLocation *time.Location
	localLocation, _ = time.LoadLocation("Asia/Shanghai")

	expTime = expTime.In(localLocation)
	expTime = expTime.Add(-8 * time.Hour)

	//time expired
	if time.Now().After(expTime) {
		return time.Time{}, errors.New("data expired")
	}

	return expTime, nil
}

//check the shortID is legal or not
//in this application, after decoding the shortID it must be multiples of 2
func ShortID_legal(shortID string) error {
	if i := Decode(shortID); i%2 == 1 {
		return errors.New("not a shortID")
	}

	return nil
}
