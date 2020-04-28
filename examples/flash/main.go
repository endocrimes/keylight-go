package main

import (
	"context"
	"log"
	"time"

	"github.com/endocrimes/keylight-go"
)

func mapLights(lights []*keylight.KeyLightLight, f func(*keylight.KeyLightLight)) {
	for _, light := range lights {
		f(light)
	}
}

func main() {
	discovery, err := keylight.NewDiscovery()
	if err != nil {
		log.Fatalf("failed to initialize keylight discovery, err: %v", err)
	}

	discoveryCtx, discoveryShutdownFn := context.WithCancel(context.Background())
	go func() {
		err := discovery.Run(discoveryCtx)
		if err != nil {
			log.Fatalf("discovery failed, err: %v", err)
		}
	}()

	firstLight := <-discovery.ResultsCh()
	discoveryShutdownFn()

	info, err := firstLight.FetchAccessoryInfo(context.TODO())
	if err != nil {
		log.Fatalf("failed to retrieve light info, err: %v", err)
	}

	log.Printf("Detected light: %s", info.DisplayName)

	log.Printf("Flashing light: %s", info.DisplayName)

	currentOpts, err := firstLight.FetchLightOptions(context.TODO())
	if err != nil {
		log.Fatalf("failed to retrieve light options, err: %v", err)
	}

	newOpts := currentOpts.Copy()
	mapLights(newOpts.Lights, func(l *keylight.KeyLightLight) {
		l.On = 1
		l.Brightness = 40
	})

	_, err = firstLight.UpdateLightOptions(context.TODO(), newOpts)
	if err != nil {
		log.Fatalf("failed to update light options, err: %v", err)
	}

	time.Sleep(3 * time.Second)

	_, err = firstLight.UpdateLightOptions(context.TODO(), currentOpts)
	if err != nil {
		log.Fatalf("failed to reset light options, err: %v", err)
	}
}