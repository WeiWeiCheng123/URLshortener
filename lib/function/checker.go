package function

import (
	"errors"
	"net/url"
	"time"
)

//Valid the URL, if URL not correct then return false
func IsURL(OriginalURL string) bool {
	u, err := url.Parse(OriginalURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

//Time formater, if time format is wrong or expired, return false.
//Otherwise, return true and convert to time type
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
		return time.Time{}, errors.New("Time expired")
	}

	return expireTime, nil
}

//Change the time to Taiwan time zone
func Time_to_Taiwanzone(expTime time.Time) (time.Time, error) {
	var localLocation *time.Location
	localLocation, _ = time.LoadLocation("Asia/Shanghai")

	expTime = expTime.In(localLocation)
	expTime = expTime.Add(-8 * time.Hour)
	//Time expired
	if time.Now().After(expTime) {
		return time.Time{}, errors.New("Time expired")
	}

	return expTime, nil
}

//Check the shortID is legal or not
//In this application, after decoding the shortID it must  be multiples of 2
func ShortID_legal(shortID string) error {
	if i := decode(shortID); i % 2 == 1 {
		return errors.New("Not a shortID")
	}

	return nil
}