package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type SleepParam struct {
	sleepMin   int
	sleepNoise int
	pause      int
	pauseNoise int
}

func (s *SleepParam) sleepFromEnv() error {
	sleepMin, err := strconv.Atoi(os.Getenv("SLEEP_MIN"))
	if err != nil {
		return fmt.Errorf("func sleepFromEnv: Failed to get sleep_min (%v): %w", sleepMin, err)
	}
	sleepNoise, err := strconv.Atoi(os.Getenv("SLEEP_NOISE"))
	if err != nil {
		return fmt.Errorf("func sleepFromEnv: Failed to get sleep_noise (%v): %w", sleepNoise, err)
	}
	pause, err := strconv.Atoi(os.Getenv("PAUSE"))
	if err != nil {
		return fmt.Errorf("func sleepFromEnv: Failed to get pause (%v): %w", pause, err)
	}
	pauseNoise, err := strconv.Atoi(os.Getenv("PAUSE_NOISE"))
	if err != nil {
		return fmt.Errorf("func sleepFromEnv: Failed to get pause_noise (%v): %w", pauseNoise, err)
	}

	s.sleepMin = sleepMin
	s.sleepNoise = sleepNoise
	s.pause = pause
	s.pauseNoise = pauseNoise

	return nil
}

func (s *SleepParam) sleep() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	sleep := s.sleepMin + rand.Intn(s.sleepNoise)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
}

func (s *SleepParam) randomPause() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var pause int
	if rand.Intn(15) == 0 {
		pause = s.pause + rand.Intn(s.pauseNoise)
		time.Sleep(time.Duration(pause) * time.Millisecond)
	}
}
