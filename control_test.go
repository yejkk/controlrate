package controlrate

import (
	"fmt"
	"testing"
	"time"
)

func TestFetchConcurrentNumNow(t *testing.T) {
	limit := NewConcurrentLimit(10000, 80)

	name1 := "test_1"

	var lastNum int
	var allSize int
	for i := 0; i < 100; i++ {
		// lastNum = 1
		num, err := limit.FetchConcurrentNumNow("name1", lastNum)
		fmt.Println("name:", name1, " num:", num, " err", err, " time:", time.Now())
		lastNum = num
		// if i == 50 || i == 54 {
		// 	time.Sleep(3 * time.Second)
		// }
		time.Sleep(time.Second / 100)
		allSize += num
	}
	fmt.Println(" allSize:", allSize)

	// wg := sync.WaitGroup{}
	// var allSize int64
	// for i := 0; i < 30; i++ {
	// 	index := fmt.Sprintf("name%d", i)
	// 	wg.Add(1)
	// 	go func() {
	// 		var lastNum int
	// 		var totalSize int
	// 		for i := 0; i < 100; i++ {
	// 			num, err := limit.FetchConcurrentNumNow(index, lastNum)
	// 			fmt.Println("name:", name1, " num:", num, " err", err, " time:", time.Now())
	// 			lastNum = num
	// 			totalSize += num
	// 		}
	// 		fmt.Println(index+" total_num:", totalSize)
	// 		atomic.AddInt64(&allSize, int64(totalSize))
	// 		wg.Done()
	// 	}()
	// }
	// wg.Wait()
	// fmt.Println(" allSize:", allSize)
}
