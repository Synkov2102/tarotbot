package tarot

// Структура для представления карты Таро
type TarotCard struct {
	Name   string
	Suit   string // "Старшие Арканы" для старших арканов
	Number int    // Номер карты, 0 для старших арканов
}

// Функция для создания колоды Таро
func CreateTarotDeck() []TarotCard {
	var deck []TarotCard

	// Старшие Арканы
	majorArcana := []string{
		"Шут", "Маг", "Верховная Жрица", "Императрица", "Император",
		"Иерофант", "Влюбленные", "Колесница", "Сила", "Отшельник",
		"Колесо Фортуны", "Правосудие", "Повешенный", "Смерть", "Умеренность",
		"Дьявол", "Башня", "Звезда", "Луна", "Солнце",
		"Суд", "Мир",
	}
	for i, name := range majorArcana {
		deck = append(deck, TarotCard{Name: name, Suit: "Старшие Арканы", Number: i})
	}

	// Масти младших арканов
	suits := []string{"Жезлов", "Кубков", "Мечей", "Пентаклей"}
	numbers := []string{
		"Туз", "Двойка", "Тройка", "Четверка", "Пятерка", "Шестерка", "Семерка", "Восьмерка", "Девятка", "Десятка",
		"Паж", "Рыцарь", "Королева", "Король",
	}
	for _, suit := range suits {
		for i, number := range numbers {
			deck = append(deck, TarotCard{Name: number, Suit: suit, Number: i + 1})
		}
	}

	return deck
}
