package trainingTimeZone

import (
	"fmt"
	"time"
)

const location = "Asia/Tokyo"

// MainTimeZone タイムゾーンの勉強
func MainTimeZone() {
	// タイムゾーンが存在するか確認
	loc, err := time.LoadLocation(location)
	// もし Asia/Tokyo が無ければ自分で設定する
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}

	// タイムゾーン変更前の時間
	now01 := time.Now()
	fmt.Println(now01)

	// タイムゾーン分の調整をする書き方
	now02 := time.Now().In(loc)
	fmt.Println(now02)

	// タイムゾーンを変更する
	time.Local = loc
	// タイムゾーン変更後の時間
	now03 := time.Now()
	fmt.Println(now03)
}
