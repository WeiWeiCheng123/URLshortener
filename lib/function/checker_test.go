package function

import "testing"

func Test_IsUrl_is_URL(t *testing.T) {
	//Set a correct URL format
	//it should return true
	url := "https://www.dcard.tw/f"
	if IsURL(url) != true {
		t.Error("Wrong result")
	}
}

func Test_IsURL_not_a_URL1(t *testing.T) {
	//Set a wrong URL format (loss :)
	//it should return false
	url := "https//www.dcard.tw/f"
	if IsURL(url) != false {
		t.Error("Wrong result")
	}
}

func Test_IsURL_not_a_URL2(t *testing.T) {
	//Set a wrong URL format (loss /)
	//it should return false
	url := "https:/www.dcard.tw/f"
	if IsURL(url) != false {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_Pass(t *testing.T) {
	//Set a correct time format
	//it should not retrun error
	time := "2023-02-12T15:04:05Z"
	_, err := TimeFormater(time)
	if err != nil {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_Fail_in_Day(t *testing.T) {
	//Set a wrong time format (loss day)
	//it should return error
	time := "2023-01-T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_Fail_in_Month(t *testing.T) {
	//Set a wrong time format (loss month)
	//it should return error
	time := "2023-13-01T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_Fail_in_Year(t *testing.T) {
	//Set a wrong time format (loss year)
	//it should return error
	time := "1-13-01T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}

func Test_TimeFormater_is_Expired(t *testing.T) {
	//Set a wrong time format (time expired)
	//it should return error
	time := "2022-02-10T15:04:05Z"
	_, err := TimeFormater(time)
	if err == nil {
		t.Error("Wrong result")
	}
}
