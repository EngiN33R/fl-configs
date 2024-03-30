package equipment_mapped

import (
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/iniload"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/semantic"

	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type Commodity struct {
	semantic.Model
	Nickname  *semantic.String
	Equipment *semantic.String
	Category  *semantic.String

	Price         *semantic.Int
	Combinable    *semantic.Bool
	GoodSellPrice *semantic.Float
	BadBuyPrice   *semantic.Float
	BadSellPrice  *semantic.Float
	GoodBuyPrice  *semantic.Float
	ShopArchetype *semantic.Path
	ItemIcon      *semantic.Path
	JumpDist      *semantic.Int
}

type Ship struct {
	semantic.Model
	Category *semantic.String
	Nickname *semantic.String
	Hull     *semantic.String
}
type ShipHull struct {
	semantic.Model
	Nickname *semantic.String
	Category *semantic.String
	Ship     *semantic.String
	Price    *semantic.Int
	IdsName  *semantic.Int
}

type Good struct {
	semantic.Model
	Category *semantic.String
	Nickname *semantic.String
	Price    *semantic.Int
}

type Config struct {
	Files []*iniload.IniLoader

	Goods    []*Good
	GoodsMap map[string]*Good

	Commodities    []*Commodity
	CommoditiesMap map[string]*Commodity
	Ships          []*Ship
	ShipsMap       map[string]*Ship
	ShipHulls      []*ShipHull
	ShipHullsMap   map[string]*ShipHull
}

const (
	FILENAME utils_types.FilePath = "goods.ini"
)

func Read(configs []*iniload.IniLoader) *Config {
	frelconfig := &Config{Files: configs}
	frelconfig.Commodities = make([]*Commodity, 0, 100)
	frelconfig.CommoditiesMap = make(map[string]*Commodity)
	frelconfig.Ships = make([]*Ship, 0, 100)
	frelconfig.ShipsMap = make(map[string]*Ship)
	frelconfig.ShipHulls = make([]*ShipHull, 0, 100)
	frelconfig.ShipHullsMap = make(map[string]*ShipHull)

	frelconfig.Goods = make([]*Good, 0, 100)
	frelconfig.GoodsMap = make(map[string]*Good)

	for _, config := range configs {
		for _, section := range config.SectionMap["[Good]"] {
			good := &Good{}
			good.Map(section)
			good.Nickname = semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			good.Category = semantic.NewString(section, "category", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
			good.Price = semantic.NewInt(section, "price", semantic.Optional())
			frelconfig.Goods = append(frelconfig.Goods, good)
			frelconfig.GoodsMap[good.Nickname.Get()] = good

			category := good.Category.Get()
			switch category {
			case "commodity":
				commodity := &Commodity{}
				commodity.Map(section)
				commodity.Category = semantic.NewString(section, "category", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				commodity.Nickname = semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				commodity.Equipment = semantic.NewString(section, "equipment", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				commodity.Price = semantic.NewInt(section, "price")
				commodity.Combinable = semantic.NewBool(section, "combinable", semantic.StrBool)
				commodity.GoodSellPrice = semantic.NewFloat(section, "good_sell_price", semantic.Precision(2))
				commodity.BadBuyPrice = semantic.NewFloat(section, "bad_buy_price", semantic.Precision(2))
				commodity.BadSellPrice = semantic.NewFloat(section, "bad_sell_price", semantic.Precision(2))
				commodity.GoodBuyPrice = semantic.NewFloat(section, "good_buy_price", semantic.Precision(2))
				commodity.ShopArchetype = semantic.NewPath(section, "shop_archetype")
				commodity.ItemIcon = semantic.NewPath(section, "item_icon")
				commodity.JumpDist = semantic.NewInt(section, "jump_dist")

				frelconfig.Commodities = append(frelconfig.Commodities, commodity)
				frelconfig.CommoditiesMap[commodity.Nickname.Get()] = commodity
			case "ship":
				ship := &Ship{}
				ship.Map(section)
				ship.Category = semantic.NewString(section, "category")
				ship.Nickname = semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				ship.Hull = semantic.NewString(section, "hull")

				frelconfig.Ships = append(frelconfig.Ships, ship)
				frelconfig.ShipsMap[ship.Nickname.Get()] = ship
			case "shiphull":
				shiphull := &ShipHull{}
				shiphull.Map(section)
				shiphull.Category = semantic.NewString(section, "category", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				shiphull.Nickname = semantic.NewString(section, "nickname", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				shiphull.Ship = semantic.NewString(section, "ship", semantic.WithLowercaseS(), semantic.WithoutSpacesS())
				shiphull.Price = semantic.NewInt(section, "price")
				shiphull.IdsName = semantic.NewInt(section, "ids_name")

				frelconfig.ShipHulls = append(frelconfig.ShipHulls, shiphull)
				frelconfig.ShipHullsMap[shiphull.Nickname.Get()] = shiphull
			}

		}
	}

	return frelconfig
}

func (frelconfig *Config) Write() []*file.File {
	var files []*file.File
	for _, file := range frelconfig.Files {
		inifile := file.Render()
		inifile.Write(inifile.File)
		files = append(files, inifile.File)
	}
	return files
}
