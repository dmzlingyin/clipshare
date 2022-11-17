//go:build client
// +build client

package utils

import (
	"os"
	"time"

	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func Play(sound string) {
	f, err := os.Open(sound)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}
