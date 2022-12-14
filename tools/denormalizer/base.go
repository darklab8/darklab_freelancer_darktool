package denormalizer

import (
	"darktool/tools/parser"
	"darktool/tools/parser/freelancer/data/universe"
	"darktool/tools/parser/freelancer/data/universe/systems"
	"darktool/tools/parser/freelancer/infocard"
	"darktool/tools/utils/err"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"darktool/tools/parser/freelancer/data/equipment/market"
	_ "darktool/tools/parser/freelancer/data/universe"
)

type DenormalizedBaseGood struct {
	name             string
	recycleCandidate string
}

type BaseDenormalizer struct {
	baseGoods map[string]*DenormalizedBaseGood
}

func (d *BaseDenormalizer) Read(parsed *parser.Parsed) *BaseDenormalizer {
	d.baseGoods = make(map[string]*DenormalizedBaseGood)

	for _, base := range parsed.Universe_config.Bases {
		d.baseGoods[base.Nickname.Get()] = &DenormalizedBaseGood{}
	}

	d.ReadBaseNames(parsed.Market_ships_config, parsed.Universe_config, parsed.Info_config)
	d.ReadRecycle(parsed.Market_ships_config, parsed.Universe_config, parsed.Systems)
	return d
}

func (d *BaseDenormalizer) Write(parsed *parser.Parsed) {
	d.MarketWrite(parsed.Market_commodities)
	d.MarketWrite(parsed.Market_ships_config)
	d.MarketWrite(parsed.Market_misc)
	d.UniverseWrite(parsed.Universe_config)
}

func (denormalizer *BaseDenormalizer) MarketWrite(marketconfig *market.Config) {

	for _, base_good := range marketconfig.BaseGoods {
		base_good.Name.Set(denormalizer.baseGoods[base_good.Base.Get()].name)

		if denormalizer.baseGoods[base_good.Base.Get()].recycleCandidate != "" {
			base_good.RecycleCandidate.Set(denormalizer.baseGoods[base_good.Base.Get()].recycleCandidate)
		} else {
			base_good.RecycleCandidate.Delete()
		}
	}
}

func (denormalizer *BaseDenormalizer) ReadBaseNames(marketconfig *market.Config, universeConfig *universe.Config, infocards *infocard.Config) {
	for _, base := range universeConfig.Bases {

		base_strid_name := int(base.StridName.Get())
		record, ok := infocards.RecordsMap[base_strid_name]
		if !ok {
			log.Fatal("failed to get record from infocardts.txt, with id: ", base_strid_name, " for base: ", base, " record = ", record)
		}

		if (*record).Type() != "NAME" {
			log.Fatal("incorrect type in infocards, expected NAME for base_strid_name: ", base_strid_name, " record: ", *record, " base: ", base)
		}

		denormalizer.baseGoods[base.Nickname.Get()].name = (*record).Content()
	}
}

var system_for_recycled_bases = [...]string{"ga13", "fp7"}

func (denormalizer *BaseDenormalizer) ReadRecycle(marketconfig *market.Config, universeConfig *universe.Config, systems *systems.Config) {
	for _, base := range universeConfig.Bases {

		system, ok := systems.SystemsMap[base.System.Get()]
		if !ok {
			log.Fatal("base ", base.Nickname, "is leading to non existent system", base.System.Get())
		}

		var recycle_builder strings.Builder

		universe_system, ok := universeConfig.SystemMap[universe.SystemNickname(system.Nickname)]
		_, ok = system.BasesByBase[base.Nickname.Get()]
		if !ok {
			recycle_builder.WriteString(fmt.Sprintf("%s base_good.base=%s not in universe.ini->Base.system->System.file->%s | ", err.Base_is_not_present_in_system, base.Nickname.Get(), universe_system.File.Get()))
			// log.Fatal("base ", base_good.Base, " is not present in system ", system.Nickname, " potential crash situation")
		}

		for _, recycle_system := range system_for_recycled_bases {
			if base.System.Get() == recycle_system {
				recycle_builder.WriteString(fmt.Sprintf("universe.ini->Base.system=%s in [%v]", recycle_system, system_for_recycled_bases))
			}
		}

		recycleCandidate := recycle_builder.String()
		denormalizer.baseGoods[base.Nickname.Get()].recycleCandidate = recycleCandidate

	}
}

func (denormalizer *BaseDenormalizer) UniverseWrite(universeConfig *universe.Config) {

	for base_nickname, base := range universeConfig.BasesMap {
		base.Name.Set(denormalizer.baseGoods[string(base_nickname)].name)

		if denormalizer.baseGoods[string(base_nickname)].recycleCandidate != "" {
			base.RecycleCandidate.Set(denormalizer.baseGoods[string(base_nickname)].recycleCandidate)
		} else {
			base.RecycleCandidate.Delete()
		}
	}
}
