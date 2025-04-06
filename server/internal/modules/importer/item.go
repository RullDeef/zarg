package importer

import (
	"context"
	"errors"
	"os"
	"server/domain"
	"server/internal/modules/logger"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrUnknownItemID    = errors.New("unknown item id")
	ErrUnknownItemKind  = errors.New("unknown item kind")
	ErrInvalidUseCaseID = errors.New("invalid use case id")
)

var predefinedActions = map[string]func(context.Context, *domain.PickableItem) error{
	"attack/melee": attackMelee,
}

type ItemImporter struct {
	path  string
	items map[int]*domain.PickableItem
}

type ItemYaml struct {
	ID                int
	Title             string
	Description       string
	Kind              string
	Weight            float64
	Cost              int
	IsStoryline       bool `yaml:",omitempty"`
	IsWeapon          bool `yaml:",omitempty"`
	IsEquipable       bool `yaml:",omitempty"`
	KeepOnEscape      bool `yaml:",omitempty"`
	KeepOnDeath       bool `yaml:",omitempty"`
	DropAfterCompaign bool `yaml:",omitempty"`
	Rarity            float64
	UseCases          []UseCaseYaml
}

type UseCaseYaml struct {
	Title            string
	Description      string
	CanBeUsedOnRest  bool `yaml:",omitempty"`
	CanBeUsedInFight bool `yaml:",omitempty"`
	UseSpendsMove    bool `yaml:",omitempty"`
	IsDestructive    bool `yaml:",omitempty"`
	UsesLeft         int
	ActionID         string
}

func NewItemImporter(path string) *ItemImporter {
	return &ItemImporter{
		path:  path,
		items: nil,
	}
}

func (ii *ItemImporter) Import(ctx context.Context, id int) (*domain.PickableItem, error) {
	log := logger.GetSugaredLogger(ctx)
	if ii.items == nil {
		items, err := loadItems(ii.path)
		if err != nil {
			log.Errorw("failed to load items", "path", ii.path, "err", err)
			return nil, err
		}
		log.Infow("items loaded", "path", ii.path, "items_count", len(items))
		ii.items = items
	}

	if item, ok := ii.items[id]; !ok {
		return nil, ErrUnknownItemID
	} else {
		return item, nil
	}
}

func loadItems(path string) (map[int]*domain.PickableItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var itemsYaml []ItemYaml
	err = yaml.NewDecoder(f).Decode(&itemsYaml)
	if err != nil {
		return nil, err
	}

	items := make(map[int]*domain.PickableItem)
	for _, itemYaml := range itemsYaml {
		item, id, err := toItem(itemYaml)
		if err != nil {
			return nil, err
		}
		items[id] = item
	}
	return items, nil
}

func toItem(yml ItemYaml) (*domain.PickableItem, int, error) {
	kind, err := toItemKind(yml.Kind)
	if err != nil {
		return nil, 0, err
	}

	useCases := make([]*domain.ItemUseCase, 0)
	for _, uc := range yml.UseCases {
		if action, ok := predefinedActions[uc.ActionID]; !ok {
			return nil, 0, ErrInvalidUseCaseID
		} else {
			useCases = append(useCases, &domain.ItemUseCase{
				Title:            uc.Title,
				Description:      uc.Description,
				CanBeUsedOnRest:  uc.CanBeUsedOnRest,
				CanBeUsedInFight: uc.CanBeUsedInFight,
				UseSpendsMove:    uc.UseSpendsMove,
				IsDestructive:    uc.IsDestructive,
				UsesLeft:         uc.UsesLeft,
				Action:           action,
			})
		}
	}

	return &domain.PickableItem{
		Title:             yml.Title,
		Description:       yml.Description,
		Kind:              kind,
		Weight:            yml.Weight,
		Cost:              yml.Cost,
		IsStoryline:       yml.IsStoryline,
		IsWeapon:          yml.IsWeapon,
		IsEquipable:       yml.IsEquipable,
		KeepOnEscape:      yml.KeepOnEscape,
		KeepOnDeath:       yml.KeepOnDeath,
		DropAfterCompaign: yml.DropAfterCompaign,
		Rarity:            yml.Rarity,
		UseCases:          useCases,
	}, yml.ID, nil
}

func toItemKind(kind string) (domain.ItemKind, error) {
	kindMap := map[string]domain.ItemKind{
		"weapon/melee":           domain.ItemKindWeaponMelee,
		"weapon/ranged":          domain.ItemKindWeaponRanged,
		"weapon/magic":           domain.ItemKindWeaponMagic,
		"weapon/throwable":       domain.ItemKindWeaponThrowable,
		"armor/cuirass":          domain.ItemKindArmorCuirass,
		"armor/tassets":          domain.ItemKindArmorTassets,
		"armor/helmet":           domain.ItemKindArmorHelmet,
		"armor/gloves":           domain.ItemKindArmorGloves,
		"armor/boots":            domain.ItemKindArmorBoots,
		"potion/healing":         domain.ItemKindPotionHealing,
		"potion/effect":          domain.ItemKindPotionEffect,
		"potion/boost":           domain.ItemKindPotionBoost,
		"special/carpenterTool":  domain.ItemKindSpecialCarpenterTool,
		"special/jewel":          domain.ItemKindSpecialJewel,
		"special/enchantedStone": domain.ItemKindSpecialEnchantedStone,
		"charm":                  domain.ItemKindCharm,
	}

	val, ok := kindMap[strings.ToLower(kind)]
	if !ok {
		return 0, ErrUnknownItemKind
	}
	return val, nil
}

func attackMelee(ctx context.Context, item *domain.PickableItem) error {
	log := logger.GetSugaredLogger(ctx)
	log.Infow("attackMelee", "item", item)
	return nil
}
