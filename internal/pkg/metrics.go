package requests

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

type Timing struct {
	Count    int
	Duration time.Duration
}

type ctxTimings struct {
	mu   sync.Mutex
	Data map[string]*Timing
}

func TrackContextTimings(ctx context.Context, metricName string, start time.Time) {
	timings, ok := ctx.Value(TimingsKey).(*ctxTimings)
	if !ok {
		return
	}
	elapsed := time.Since(start)

	timings.mu.Lock()
	defer timings.mu.Unlock()

	if metric, metricExist := timings.Data[metricName]; !metricExist {
		timings.Data[metricName] = &Timing{
			Count:    1,
			Duration: elapsed,
		}
	} else {
		metric.Count++
		metric.Duration += elapsed
	}
}

func LogContextTimings(ctx context.Context, path string, start time.Time) {
	// получаем тайминги из контекста
	// поскольку там пустой интерфейс, то нам надо преобразовать к нужному типу
	timings, ok := ctx.Value(TimingsKey).(*ctxTimings)
	if !ok {
		return
	}
	//totalReal := time.Since(start)

	path = strings.Replace(path, "/", "-", -1)

	prefix := "requests." + path + "."

	// buf := bytes.NewBufferString(path)
	var total time.Duration
	for _, value := range timings.Data {
		//	metric := prefix + "timings." + timing
		//	tm.StatsReciever.Increment(metric)
		//	tm.StatsReciever.Timing(metric+"_time", uint64(value.Duration/time.Millisecond))
		total += value.Duration
	}
	//
	//tm.StatsReciever.Increment(prefix + "hits")
	//tm.StatsReciever.Timing(prefix+"tracked", uint64(totalReal/time.Millisecond))
	//tm.StatsReciever.Timing(prefix+"real_time", uint64(total/time.Millisecond))

	log.Println(prefix+"real_time", uint64(total/time.Millisecond))
}
