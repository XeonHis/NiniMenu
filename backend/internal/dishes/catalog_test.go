package dishes

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDefaultDishesAreComplete(t *testing.T) {
	categorySet := toSet(DefaultCategories())
	tasteSet := toSet(DefaultTastes())
	seenNames := map[string]bool{}

	packs := DefaultPacks()
	if len(packs) != len(DefaultCategories()) {
		t.Fatalf("DefaultPacks() length = %d, want %d", len(packs), len(DefaultCategories()))
	}
	for _, pack := range packs {
		if len(pack.Recipes) < 3 {
			t.Fatalf("%s recipe count = %d, want at least 3", pack.Name, len(pack.Recipes))
		}
		if !categorySet[pack.Category] {
			t.Fatalf("%s category %q is not in defaults", pack.Name, pack.Category)
		}
	}

	dishes := DefaultDishes()
	expectedDishCount := recipeCount(packs)
	if len(dishes) != expectedDishCount {
		t.Fatalf("DefaultDishes() length = %d, want %d", len(dishes), expectedDishCount)
	}

	for _, dish := range dishes {
		if strings.TrimSpace(dish.Name) == "" {
			t.Fatalf("dish has empty name: %+v", dish)
		}
		if seenNames[dish.Name] {
			t.Fatalf("duplicate dish name: %s", dish.Name)
		}
		seenNames[dish.Name] = true

		if !categorySet[dish.Category] {
			t.Fatalf("%s category %q is not in defaults", dish.Name, dish.Category)
		}
		for _, taste := range splitTaste(dish.Taste) {
			if !tasteSet[taste] {
				t.Fatalf("%s taste %q is not in defaults", dish.Name, taste)
			}
		}
		if dish.MealType != "all" && dish.MealType != "lunch" && dish.MealType != "dinner" {
			t.Fatalf("%s meal_type = %q", dish.Name, dish.MealType)
		}
		if dish.Difficulty != "easy" && dish.Difficulty != "medium" && dish.Difficulty != "hard" {
			t.Fatalf("%s difficulty = %q", dish.Name, dish.Difficulty)
		}
		if dish.CookTime <= 0 {
			t.Fatalf("%s cook_time should be positive", dish.Name)
		}
		if jsonArrayLen(dish.Ingredients) == 0 {
			t.Fatalf("%s has no ingredients", dish.Name)
		}
		if jsonArrayLen(dish.Seasonings) == 0 {
			t.Fatalf("%s has no seasonings", dish.Name)
		}
		if jsonArrayLen(dish.Steps) == 0 {
			t.Fatalf("%s has no steps", dish.Name)
		}
		if jsonArrayLen(dish.Tags) == 0 {
			t.Fatalf("%s has no tags", dish.Name)
		}
	}
}

func toSet(values []string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}

func splitTaste(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	for _, sep := range []string{"，", "、", "/", "|", ";", "；", " "} {
		raw = strings.ReplaceAll(raw, sep, ",")
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func jsonArrayLen(raw string) int {
	var values []any
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return 0
	}
	return len(values)
}
