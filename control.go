package controlrate

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	defaultStartNum     = 1
	defaultLimitTimeout = 30 * time.Second

	defaultCapacity float64 = 0.75

	defaultFastIncreaseMaxNum = 31

	defaultSlowIncrease = 1
	defaultFastIncrease = 2
)

type rateLimiter interface {
	WaitN(ctx context.Context, n int) (err error)
}

// ConcurrentLimit is a concurrent limit by total and once limit num in one time.
type ConcurrentLimit struct {
	totalNum int

	onceNum int

	lastInfo sync.Map

	limiter rateLimiter
}

type limitInfo struct {
	num   int
	times int //连续的次数，防止间隔时间误判

	updateTime time.Time
}

func (i *limitInfo) increaseTimes(count int) {
	i.times += count
}

func (i *limitInfo) decreaseTimes(count int) {
	if i.times < count {
		i.times = 0
		return
	}
	i.times -= count
}

func (i *limitInfo) increaseNum(count int) {
	i.num += count
}

func (i *limitInfo) decreaseNum(count int) {
	if i.num-count < 1 {
		i.num = 1
		return
	}
	i.num -= count
}

func (i *limitInfo) zeroTimes() {
	i.times = 0
}

func (i *limitInfo) shouldIncreaseNum(lastLoadNum int) bool {
	if lastLoadNum <= 0 {
		return false
	}
	return float64(i.num)*defaultCapacity <= float64(lastLoadNum)
}

func (i *limitInfo) reset() {
	i.updateTime = time.Now()
	i.times = 1
	i.num = defaultStartNum
}

// NewConcurrentLimit returns a concurrent limit using the provided total and once limit num.
func NewConcurrentLimit(totalNum, onceNum int) *ConcurrentLimit {
	limiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(totalNum)), onceNum)
	limit := &ConcurrentLimit{
		totalNum: totalNum,
		limiter:  limiter,
		onceNum:  onceNum,
	}
	return limit
}

// FetchConcurrentNumNow returns num which can be used now.
func (s *ConcurrentLimit) FetchConcurrentNumNow(name string, lastLoadNum int) (int, error) {
	infoContent, ok := s.lastInfo.Load(name)
	currentTime := time.Now()
	var info limitInfo
	if ok {
		info, ok = infoContent.(limitInfo)
		if !ok {
			return 0, errors.New("map store error")
		}
	}
	if !ok || (currentTime.Sub(info.updateTime) > defaultLimitTimeout && info.times <= 1) {
		info.reset()
		s.lastInfo.Store(name, info)
		err := s.limiter.WaitN(context.TODO(), info.num)
		if err != nil {
			return 0, err
		}
		return info.num, nil
	}

	shouldIncrease := info.shouldIncreaseNum(lastLoadNum)
	if info.num < s.onceNum-1 && shouldIncrease {
		if info.num > defaultFastIncreaseMaxNum {
			if ((info.times - defaultFastIncreaseMaxNum) >> 2) > info.num-defaultFastIncreaseMaxNum {
				info.increaseNum(defaultSlowIncrease)
			}
		} else if info.times>>1 > info.num {
			info.increaseNum(defaultFastIncrease)
		}
	} else if !shouldIncrease && (info.times>>1) < info.num && info.num > 1 {
		info.decreaseNum(1)
	}
	err := s.limiter.WaitN(context.TODO(), info.num)
	if err != nil {
		return 0, err
	}

	if currentTime.Sub(info.updateTime) > defaultLimitTimeout {
		info.zeroTimes()
	} else if shouldIncrease {
		info.increaseTimes(1)
	} else {
		info.decreaseTimes(1)
	}
	info.updateTime = currentTime
	s.lastInfo.Store(name, info)
	return info.num, nil
}
