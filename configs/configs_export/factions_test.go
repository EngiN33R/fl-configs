package configs_export

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped"
	"github.com/stretchr/testify/assert"
)

func TestFaction(t *testing.T) {
	configs := configs_mapped.TestFixtureConfigs()
	exporter := NewExporter(configs)

	items := exporter.GetFactions()
	assert.Greater(t, len(items), 0)

	infocards := exporter.infocards_parser.Get()
	for _, faction := range items {
		lines := infocards.MapGet(faction.Infocard)
		fmt.Println(lines.Lines)
		break
	}
}
