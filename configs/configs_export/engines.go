package configs_export

type Engine struct {
	Name  string
	Price int

	CruiseSpeed      int
	CruiseChargeTime int
	LinearDrag       int
	MaxForce         int
	ReverseFraction  float64
	ImpulseSpeed     float64

	HpType      string
	FlameEffect string
	TrailEffect string

	Nickname string
	NameID   int
	InfoID   int

	Bases []GoodAtBase
}

func (e *Exporter) GetEngines() []Engine {
	var engines []Engine

	for _, engine_info := range e.configs.Equip.Engines {
		engine := Engine{}
		engine.Nickname = engine_info.Nickname.Get()
		if cruise_speed, ok := engine_info.CruiseSpeed.GetValue(); ok {
			engine.CruiseSpeed = cruise_speed
		} else {
			engine.CruiseSpeed = e.configs.Consts.EngineEquipConsts.CRUISING_SPEED.Get()
		}
		engine.CruiseChargeTime, _ = engine_info.CruiseChargeTime.GetValue()
		engine.LinearDrag = engine_info.LinearDrag.Get()
		engine.MaxForce = engine_info.MaxForce.Get()
		engine.ReverseFraction = engine_info.ReverseFraction.Get()
		engine.ImpulseSpeed = float64(engine.MaxForce) / float64(engine.LinearDrag)

		engine.HpType, _ = engine_info.HpType.GetValue()
		engine.FlameEffect, _ = engine_info.FlameEffect.GetValue()
		engine.TrailEffect, _ = engine_info.TrailEffect.GetValue()

		engine.NameID = engine_info.IdsName.Get()
		engine.InfoID = engine_info.IdsInfo.Get()

		if good_info, ok := e.configs.Goods.GoodsMap[engine.Nickname]; ok {
			if price, ok := good_info.Price.GetValue(); ok {
				engine.Price = price
				engine.Bases = e.GetAtBasesSold(GetAtBasesInput{
					Nickname:       good_info.Nickname.Get(),
					Price:          price,
					PricePerVolume: -1,
				})
			}
		}

		if name, ok := e.configs.Infocards.Infonames[engine.NameID]; ok {
			engine.Name = string(name)
		}

		e.infocards_parser.Set(InfocardKey(engine.Nickname), engine.InfoID)

		if engine.HpType == "" {
			continue
		}

		engines = append(engines, engine)
	}
	return engines
}