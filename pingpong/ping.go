package main

import (
	"context"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration

	select {
	case <-ctx.Done():
		return
	case interval = <-reset: //대기시간 설정
	default:
	}
	if interval <= 0 {
		interval = defaultPingInterval
	}

	timer := time.NewTimer(interval) //timer를 interval로 초기화

	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {

		case <-ctx.Done(): //context 종료 시 같이 종료
			return
		case newInterval := <-reset: //새로운 대기시간 입력 시 timer 초기화
			if !timer.Stop() {
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C: //timer 종료 시 ping Write
			if _, err := w.Write([]byte("ping")); err != nil {
				return
			}
		}
		_ = timer.Reset(interval) //timer 리셋

	}

}
