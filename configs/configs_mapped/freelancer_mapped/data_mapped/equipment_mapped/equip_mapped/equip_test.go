package equip_mapped

import (
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/configfile"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/go-utils/goutils/utils"
	"github.com/darklab8/go-utils/goutils/utils/utils_filepath"
	"github.com/stretchr/testify/assert"
)

func TestReadSelectEquip(t *testing.T) {
	test_directory := utils.GetCurrrentTestFolder()
	fileref := file.NewFile(utils_filepath.Join(test_directory, FILENAME_SELECT_EQUIP))

	config := Read([]*configfile.ConfigFile{configfile.NewConfigFile(fileref).Scan()})

	assert.Greater(t, len(config.Commodities), 0, "expected finding items")

	for _, commodity := range config.Commodities {
		commodity.IdsName.Get()
	}
}
