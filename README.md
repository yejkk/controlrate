# control-rate

Sometimes we need to execute a num at once but have some limit in a second. `control-rate` can help you do this by set totalLimit and onceLimit.

## Import
`go get github.com/yejkk/controlrate`

## Usage
```go
conLimiter := utils.NewConcurrentLimit(100, 20)
lastLoadNum := 0
limitNum, _ := conLimiter.FetchConcurrentNumNow("testApp1", lastLoadNum)
lastLoadNum = limitNum
```
