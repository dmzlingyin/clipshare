//go:build client
// +build client

package utils

import (
	"os"
	"time"

	"github.com/dmzlingyin/clipshare/pkg/log"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func Play(sound string) {
	f, err := os.Open(sound)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.ErrorLogger.Println(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
}
