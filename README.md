# control-rate

Sometimes we need to execute some job at once but have limit in a second. `control-rate` can help you do this by set totalLimit and onceLimit and return num which will increase by time but never over the limit.

## Import
`go get github.com/yejkk/controlrate`

## Usage
```go
conLimiter := controlrate.NewConcurrentLimit(100, 20)
lastLoadNum := 0
for {
  limitNum, _ := conLimiter.FetchConcurrentNumNow("testApp1", lastLoadNum)
  lastLoadNum = limitNum
  ...
}
```
limitNum will never 
