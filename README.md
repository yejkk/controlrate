# control-rate

## Import
`go get github.com/yejkk/controlrate`

## Usage
```go
conLimiter := utils.NewConcurrentLimit(maxConcurrentNum, onceAppMaxConcurrentNum)
lastLoadNum := 0
limitNum, err := conLimiter.FetchConcurrentNumNow(config.App.GetAppID(), lastLoadNum)
if err != nil {
    return 
}
lastLoadNum = limitNum
```
