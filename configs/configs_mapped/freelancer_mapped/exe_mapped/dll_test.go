package exe_mapped

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"testing"

	"github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/infocard_mapped/infocard"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-configs/configs/configs_settings"
	"github.com/darklab8/fl-configs/configs/configs_settings/logus"
	"github.com/darklab8/fl-configs/configs/tests"
	"github.com/darklab8/go-utils/goutils/utils/time_measure"
	"github.com/stretchr/testify/assert"
)

func TestReadInfocards(t *testing.T) {
	game_location := configs_settings.Env.FreelancerFolder
	filesystem := filefind.FindConfigs(game_location)

	fileref := filesystem.GetFile(FILENAME_FL_INI)
	config := Read(iniload.NewLoader(fileref).Scan())

	dlls := config.GetDlls()
	infocards := GetAllInfocards(filesystem, dlls)

	assert.Greater(t, len(infocards.Infocards), 0)
	assert.Greater(t, len(infocards.Infonames), 0)

	for id, text := range infocards.Infonames {
		fmt.Println(id)
		fmt.Println(text)
		break
	}

	// Works only on Discovery dlls
	// assert.Contains(t, infocards.Infocards[132903].Content, "We just brought a load of Fertilizers")

	// fmt.Println(infocards.Infocards[196624])
	// fmt.Println("second:", infocards.Infocards[66089])

	// fmt.Println("Abandoned Depot infocard\n",
	// 	infocards.Infocards[465639],
	// 	infocards.Infocards[465639+1], // value from infocardmap.txt mapped
	// 	infocards.Infocards[500904],   // faction infocard id
	// )
}

func TestReadInfocardsToHtml(t *testing.T) {
	f, err := os.Create("prof.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	result := time_measure.TimeMeasure(func(m *time_measure.TimeMeasurer) {
		filesystem := tests.FixtureFileFind()
		fileref := filesystem.GetFile(FILENAME_FL_INI)
		config := Read(iniload.NewLoader(fileref).Scan())

		infocards := GetAllInfocards(tests.FixtureFileFind(), config.GetDlls())

		// assert.Greater(t, len(ids), 0)

		// 503718 faction BMM
		// 465639 base Bandoned Depot
		// 465640 continuation
		// infocard tail 500904

		// TradeLaneGaName := infocards.Infonames[33389]
		// TraceLaneGaInfocard := infocards.Infocards[33390]
		// _ = TraceLaneGaInfocard
		// _ = TradeLaneGaName
		// assert.Contains(t, TradeLaneGaName, "EFL Gate/Lane Parts")

		xml_stuff := infocards.Infocards[501545]
		if xml_stuff != nil {
			// Only for Discovery
			fmt.Println("xml_stuff=", xml_stuff)

			text, err := xml_stuff.XmlToText()
			logus.Log.CheckPanic(err, "unable convert to text")
			assert.Greater(t, len(text), 0)
			assert.NotEmpty(t, text)
			fmt.Println(text)
		}

	}, time_measure.WithMsg("measure time"))
	logus.Log.CheckPanic(result.ResultErr, "non nil exit")
}

func TestValidateInfocards(t *testing.T) {
	game_location := configs_settings.Env.FreelancerFolder

	filesystem := filefind.FindConfigs(game_location)
	fileref := filesystem.GetFile(FILENAME_FL_INI)
	config := Read(iniload.NewLoader(fileref).Scan())
	infocards := GetAllInfocards(filefind.FindConfigs(game_location), config.GetDlls())

	var parsed []*infocard.Infocard = make([]*infocard.Infocard, 0, 100)
	var parsed_text map[int][]string = make(map[int][]string)
	var failed []*infocard.Infocard = make([]*infocard.Infocard, 0, 100)

	for id, infocard := range infocards.Infocards {
		text, err := infocard.XmlToText()
		parsed_text[id] = text

		if logus.Log.CheckWarn(err, "unable convert to text") {
			failed = append(failed, infocard)
			fmt.Println("failed=", id, infocard.Lines)
		} else {
			parsed = append(parsed, infocard)
		}
	}

	fmt.Println("parsed_count=", len(parsed))
	assert.Equal(t, len(failed), 0, "expected no failed")
}
