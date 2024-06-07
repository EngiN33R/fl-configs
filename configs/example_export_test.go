/*
Such code is primiarily used for fl-darkstat. You could check its code for more examples
https://github.com/darklab8/fl-darkstat
*/
package configs

import (
	"fmt"

	"github.com/darklab8/fl-configs/configs/configs_export"
	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/darklab8/fl-configs/configs/configs_settings"
	"github.com/darklab8/fl-configs/configs/configs_settings/logus"

	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
)

// ExampleExportingData demonstrating exporting freelancer folder data for comfortable usage
func Example_exportingData() {
	freelancer_folder := configs_settings.Env.FreelancerFolder
	configs := configs_mapped.NewMappedConfigs()
	logus.Log.Debug("scanning freelancer folder", utils_logus.FilePath(freelancer_folder))

	// Reading to ini universal custom format and mapping to ORM objects
	// which have both reading and writing back capabilities
	configs.Read(freelancer_folder)

	// For elegantly exporting enriched data objects with better type safety for just reading access
	// it is already combined with multiple configs sources for flstat view
	exported := configs_export.Export(configs)
	for _, base := range exported.Bases {
		// do smth with exported bases
		fmt.Println(base.Name)
		fmt.Println(base.Infocard)
		fmt.Println(base.System)
		fmt.Println(base.SystemNickname)
		fmt.Printf("%d\n", base.InfocardID)
		break
	}
}
