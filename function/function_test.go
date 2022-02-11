package function

import "testing"

func Test_url_Pass(t *testing.T) {
	url := "https://www.dcard.tw/f"
	if IsUrl(url) != true {
		t.Error("Wrong result")
	}
}

func Test_url_Fail1(t *testing.T) {
	url := "https//www.dcard.tw/f"
	if IsUrl(url) != false {
		t.Error("Wrong result")
	}
}

func Test_url_Fail2(t *testing.T) {
	url := "https:/www.dcard.tw/f"
	if IsUrl(url) != false {
		t.Error("Wrong result")
	}
}

func Test_time_Pass(t *testing.T) {
	time := "2022-02-12T15:04:05Z"
	_, err := TimeFormater(time)
	if err != nil {
		t.Error("Wrong result")
	}
}

func Test_time_Fail_in_Day(t *testing.T) {
	time := "2022-01-T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_time_Fail_in_Month(t *testing.T) {
	time := "2022-13-01T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_time_Fail_in_Year(t *testing.T) {
	time := "1-13-01T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_time_is_Expire(t *testing.T) {
	time := "2022-02-10T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_time_is_not_Expire(t *testing.T) {
	time := "2022-02-10T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}
