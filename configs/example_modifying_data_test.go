/*
Such code is primiarily used for fl-darklint. You could check its code for more examples
https://github.com/darklab8/fl-darklint
*/
package configs

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/configs_fixtures"
	"github.com/darklab8/fl-configs/configs/settings/logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
)

// ExampleModifyingData demononstrating how to change configs values
func Example_modifyingConfigs() {
	freelancer_folder := configs_fixtures.FixtureGameLocation()
	configs := configs_mapped.NewMappedConfigs()
	logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))

	// Reading ini reading universal format
	// and mapping to ORM objects
	configs.Read(freelancer_folder)

	// Modifying files
	for _, base := range configs.Universe_config.Bases {
		base.Nickname.Set(base.Nickname.Get())
		base.System.Set(base.System.Get())
		base.File.Set(base.File.Get())
	}

	for _, system := range configs.Universe_config.Systems {
		system.Nickname.Set(system.Nickname.Get())
		system.Msg_id_prefix.Set(system.Msg_id_prefix.Get())

		if system.File.Get() != "" {
			system.File.Set(system.File.Get())
		}
	}

	// Write without Dry Run for writing to files modified values back!
	configs.Write(configs_mapped.IsDruRun(true))
}
