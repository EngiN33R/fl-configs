package infocard_mapped

import (
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/exe_mapped"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-configs/configs/tests"
	"github.com/darklab8/go-utils/goutils/utils"
	"github.com/stretchr/testify/assert"
)

// Not used any longer?
func TestReader(t *testing.T) {
	one_file_filesystem := filefind.FindConfigs(utils.GetCurrrentTestFolder())

	filesystem := tests.FixtureFileFind()

	freelancer_ini := exe_mapped.FixtureFLINIConfig()
	config, _ := Read(filesystem, freelancer_ini, one_file_filesystem.GetFile("temp.disco.infocards.txt"))

	assert.Greater(t, len(config.Infocards), 0)
}
