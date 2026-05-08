package ui

import (
	"os"
	"sync"
)

type SoundEvent string

const (
	SoundEventMessageReceived SoundEvent = "message_received"
	SoundEventMessageSent     SoundEvent = "message_sent"
	SoundEventError           SoundEvent = "error"
	SoundEventSuccess         SoundEvent = "success"
	SoundEventNotification    SoundEvent = "notification"
)

type SoundManager struct {
	mu           sync.RWMutex
	enabled      bool
	sounds       map[SoundEvent]string
	defaultSound string
}

func NewSoundManager() *SoundManager {
	sm := &SoundManager{
		enabled:      true,
		sounds:       make(map[SoundEvent]string),
		defaultSound: "\a",
	}

	if envSound := os.Getenv("FREECODE_SOUND_MESSAGE_RECEIVED"); envSound != "" {
		sm.sounds[SoundEventMessageReceived] = envSound
	} else {
		sm.sounds[SoundEventMessageReceived] = sm.defaultSound
	}

	if envSound := os.Getenv("FREECODE_SOUND_MESSAGE_SENT"); envSound != "" {
		sm.sounds[SoundEventMessageSent] = envSound
	} else {
		sm.sounds[SoundEventMessageSent] = sm.defaultSound
	}

	if envSound := os.Getenv("FREECODE_SOUND_ERROR"); envSound != "" {
		sm.sounds[SoundEventError] = envSound
	} else {
		sm.sounds[SoundEventError] = sm.defaultSound
	}

	if envSound := os.Getenv("FREECODE_SOUND_SUCCESS"); envSound != "" {
		sm.sounds[SoundEventSuccess] = envSound
	} else {
		sm.sounds[SoundEventSuccess] = sm.defaultSound
	}

	if envSound := os.Getenv("FREECODE_SOUND_NOTIFICATION"); envSound != "" {
		sm.sounds[SoundEventNotification] = envSound
	} else {
		sm.sounds[SoundEventNotification] = sm.defaultSound
	}

	if os.Getenv("FREECODE_SOUND_ENABLED") == "false" {
		sm.enabled = false
	}

	return sm
}

func (sm *SoundManager) Play(event SoundEvent) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if !sm.enabled {
		return
	}

	sound, ok := sm.sounds[event]
	if !ok {
		sound = sm.defaultSound
	}

	os.Stderr.WriteString(sound)
}

func (sm *SoundManager) SetEnabled(enabled bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.enabled = enabled
}

func (sm *SoundManager) IsEnabled() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.enabled
}

func (sm *SoundManager) SetSound(event SoundEvent, sound string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sounds[event] = sound
}

func (sm *SoundManager) GetSound(event SoundEvent) string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sound, ok := sm.sounds[event]; ok {
		return sound
	}
	return sm.defaultSound
}