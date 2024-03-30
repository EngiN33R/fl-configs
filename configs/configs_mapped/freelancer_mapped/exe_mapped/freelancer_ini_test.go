package exe_mapped

import (
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"
	"github.com/darklab8/go-utils/goutils/utils"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	config := FixtureFLINIConfig()
	assert.Greater(t, len(config.Dlls), 0)

	equips := utils.CompL(config.Equips, func(x *semantic.Path) utils_types.FilePath { return x.Get() })
	assert.Greater(t, len(equips), 0)
}
