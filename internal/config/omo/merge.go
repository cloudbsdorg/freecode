package omo

import (
	"github.com/freecode/freecode/internal/config"
)

func Merge(cfg *config.Config, omo *OMOConfig) {
	if omo.SkillsDir != "" {
		cfg.Platform.CacheDir = omo.SkillsDir
	}

	if omo.SlopRemove {
		cfg.Hooks.Session = append(cfg.Hooks.Session, "slop_remove")
	}
}

func MergeInto(cfg *config.Config, omoPath string) error {
	omo, err := Read(omoPath)
	if err != nil {
		return err
	}
	Merge(cfg, omo)
	return nil
}
