package controlrate

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchConcurrentNumNow(t *testing.T) {
	limit := NewConcurrentLimit(10000, 80)

	var lastNum int
	var allSize int
	var maxNum int
	for i := 0; i < 1000; i++ {
		num, _ := limit.FetchConcurrentNumNow("name_test1", lastNum)
		lastNum = num
		allSize += num
		if maxNum < num {
			maxNum = num
		}
	}
	assert.Equal(t, 79, maxNum)
	assert.Greater(t, allSize, 70000)
	assert.Less(t, allSize, 80000)

	maxNum = 0
	wg := sync.WaitGroup{}
	var size int64 = 0
	for i := 0; i < 30; i++ {
		index := fmt.Sprintf("name%d", i)
		wg.Add(1)
		go func() {
			var lastNum int
			var totalSize int
			for i := 0; i < 100; i++ {
				num, _ := limit.FetchConcurrentNumNow(index, lastNum)
				lastNum = num
				totalSize += num
				if maxNum < num {
					maxNum = num
				}
			}
			atomic.AddInt64(&size, int64(totalSize))
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Less(t, maxNum, 79)
	assert.Greater(t, size, int64(75000))
	assert.Less(t, size, int64(80000))
}

func TestLimitInfo(t *testing.T) {
	info := limitInfo{
		updateTime: time.Now(),
		times:      1,
		num:        1,
	}

	info.increaseTimes(1)
	assert.Equal(t, 2, info.times)

	info.decreaseTimes(1)
	assert.Equal(t, 1, info.times)
	info.decreaseTimes(10)
	assert.Equal(t, 0, info.times)

	info.increaseNum(1)
	assert.Equal(t, 2, info.num)

	info.decreaseNum(1)
	assert.Equal(t, 1, info.num)
	info.decreaseNum(10)
	assert.Equal(t, 1, info.num)

	info.increaseTimes(10)
	info.zeroTimes()
	assert.Equal(t, 0, info.times)

	info.reset()
	assert.Equal(t, 1, info.times)
	assert.Equal(t, defaultStartNum, info.num)

	info.increaseNum(12)
	assert.False(t, info.shouldIncreaseNum(1))
	assert.True(t, info.shouldIncreaseNum(10))
}

func BenchmarkFetchConcurrentNumNowLarge(b *testing.B) {
	limit := NewConcurrentLimit(1000000, 8000)
	var lastNum int
	for i := 0; i < b.N; i++ {
		num, _ := limit.FetchConcurrentNumNow("name_test1", lastNum)
		lastNum = rand.Intn(num)
	}
}

func BenchmarkFetchConcurrentNumNowSmall(b *testing.B) {
	limit := NewConcurrentLimit(1000, 50)
	var lastNum int
	for i := 0; i < b.N; i++ {
		num, _ := limit.FetchConcurrentNumNow("name_test1", lastNum)
		lastNum = rand.Intn(num)
	}
}
