package exe_mapped

import (
	"fmt"
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped/configs_fixtures"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind"
	"github.com/stretchr/testify/assert"
)

func TestReadInfocards(t *testing.T) {
	game_location := configs_fixtures.FixtureGameLocation()
	config := FixtureFLINIConfig()
	ids := GetAllInfocards(filefind.FindConfigs(game_location), config.Resources.Dll)

	assert.Greater(t, len(ids), 0)

	for id, text := range ids {
		fmt.Println(id)
		fmt.Println(text)
		break
	}

	assert.Contains(t, ids[132903], "We just brought a load of Fertilizers")

	fmt.Println(ids[196624])
	fmt.Println("second:", ids[66089])

	fmt.Println("Abandoned Depot infocard\n",
		ids[465639],
		ids[465639+1], // value from infocardmap.txt mapped
		ids[500904],   // faction infocard id
	)
}
