# go-timings

Go package implementing interface and methods for background timing monitors

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-timings.svg)](https://pkg.go.dev/github.com/sfomuseum/go-timings)

## Exmaple

```
import (
       "context"
       "github.com/sfomuseum/go-timings"
       "os"       
       "time"
)

func main() {

	ctx := context.Background()
	
	d := time.Second * 60
	monitor, _ := timings.NewCounterMonitor(ctx, d)

	monitor.Start(ctx, os.Stdout)
	defer monitor.Stop(ctx)

	monitor.Signal()	// increments by 1
}
```