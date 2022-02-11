package function

import "testing"

func Test_IsUrl_is_URL(t *testing.T) {
	url := "https://www.dcard.tw/f"
	if IsURL(url) != true {
		t.Error("Wrong result")
	}
}

func Test_IsURL_not_a_URL(t *testing.T) {
	url := "https//www.dcard.tw/f"
	if IsURL(url) != false {
		t.Error("Wrong result")
	}
}

func Test_IsUrl_not_a_URL(t *testing.T) {
	url := "https:/www.dcard.tw/f"
	if IsURL(url) != false {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_Pass(t *testing.T) {
	time := "2022-02-12T15:04:05Z"
	_, err := TimeFormater(time)
	if err != nil {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_Fail_in_Day(t *testing.T) {
	time := "2022-01-T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_Fail_in_Month(t *testing.T) {
	time := "2022-13-01T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_Fail_in_Year(t *testing.T) {
	time := "1-13-01T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_is_Expired(t *testing.T) {
	time := "2022-02-10T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_is_not_Expired(t *testing.T) {
	time := "2022-02-10T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}
