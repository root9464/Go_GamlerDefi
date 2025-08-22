package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/go-faker/faker/v4"
)

// Player представляет состояние игрока
type Player struct {
	Name      string
	Level     int
	Guild     string
	IsOnline  bool
	IPAddress string
}

var guilds = []string{"Shadow Hunters", "Dragon Knights", "Mystic Order", "Iron Legion", "Crimson Guard"}
var locations = []string{"Лес Теней", "Горы Дракона", "Столица Королевства", "Подземелья Ужаса", "Пустоши", "Эльфийские Руины", "Порт Черного Паруса", "Вулкан Рока"}
var monsters = []string{"Зеленый дракон", "Скелет-воин", "Гоблин-разбойник", "Древний лич", "Гигантский паук", "Огненный элементаль", "Темный рыцарь", "Вампир-лорд", "Морской змей", "Каменный голем"}
var quests = []string{"Уничтожение орков", "Спасение деревни", "Поиск древнего артефакта", "Охота на дракона", "Тайны подземелья", "Королевский заказ", "Проклятие некроманта", "Затерянные сокровища"}
var itemRarity = []string{"Обычный", "Редкий", "Эпический", "Легендарный"}

// Используем Mutex для безопасной записи в консоль из разных горутин
var logMutex = &sync.Mutex{}

func main() {
	// Инициализация генератора случайных чисел
	rand.New(rand.NewSource(time.Now().UnixNano()))

	logMessage("SYSTEM", "Сервер 'DragonRealm Online' запущен (версия 2.0.1)")
	numPlayersOnline := rand.Intn(1500) + 500
	logMessage("SYSTEM", fmt.Sprintf("Подключено %d игрока", numPlayersOnline))
	fmt.Println("----------------------------------------")

	var wg sync.WaitGroup
	numSimulatedPlayers := 3 // Количество игроков, чьи действия мы будем симулировать

	for i := 0; i < numSimulatedPlayers; i++ {
		wg.Add(1)
		go simulatePlayerSession(&wg)
	}

	// Добавим случайное системное событие во время симуляции
	time.Sleep(time.Second * time.Duration(rand.Intn(10)+5))
	logMessage("SYSTEM", "Начался ивент 'Вторжение драконов'!")

	wg.Wait()

	fmt.Println("----------------------------------------")
	logMessage("SYSTEM", "Завершение работы сервера...")
}

// simulatePlayerSession имитирует полную игровую сессию одного игрока
func simulatePlayerSession(wg *sync.WaitGroup) {
	defer wg.Done()

	player := &Player{
		Name:      generateFantasyName(),
		Level:     rand.Intn(20) + 1,
		Guild:     guilds[rand.Intn(len(guilds))],
		IsOnline:  true,
		IPAddress: faker.IPv4(),
	}

	logMessage("AUTH", fmt.Sprintf("Игрок '%s' вошел в игру (IP: %s)", player.Name, player.IPAddress))
	logMessage("GUILD", fmt.Sprintf("Игрок '%s' представляет гильдию '%s'", player.Name, player.Guild))

	sessionDuration := rand.Intn(15) + 5 // Игрок будет в игре от 5 до 20 секунд
	for i := 0; i < sessionDuration; i++ {
		printRandomGameEvent(player)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)+200))
	}

	logMessage("AUTH", fmt.Sprintf("Игрок '%s' вышел из игры", player.Name))
}

// printRandomGameEvent выбирает и печатает случайное игровое событие
func printRandomGameEvent(p *Player) {
	actions := []func(*Player){
		logCombat,
		logLoot,
		logLevelUp,
		logQuest,
		logMove,
		logTrade,
		logPlayerInteraction,
		logError,
	}
	actions[rand.Intn(len(actions))](p)
}

// logMessage - потокобезопасная функция для вывода логов
func logMessage(tag, message string) {
	logMutex.Lock()
	fmt.Printf("[%s] %s\n", tag, message)
	logMutex.Unlock()
}

// --- Генераторы логов ---

func logCombat(p *Player) {
	monster := monsters[rand.Intn(len(monsters))]
	dmg := rand.Intn(p.Level*150) + 500
	attackType := "физической"
	if rand.Intn(100) > 60 {
		attackType = "магической"
	}

	crit := ""
	if rand.Intn(100) > 85 { // 15% шанс крита
		dmg *= 2
		crit = " (Критический удар!)"
	}

	logMessage("COMBAT", fmt.Sprintf("%s (Ур. %d) атаковал %s с помощью %s атаки и нанес %d урона!%s", p.Name, p.Level, monster, attackType, dmg, crit))

	if rand.Intn(100) > 70 { // 30% шанс убить монстра
		logMessage("COMBAT", fmt.Sprintf("%s побежден игроком %s!", monster, p.Name))
	}
}

func logLoot(p *Player) {
	rarity := itemRarity[rand.Intn(len(itemRarity))]
	// ИСПРАВЛЕНО: Используем прямой вызов faker.LastName() для получения строки.
	item := faker.LastName()
	logMessage("LOOT", fmt.Sprintf("%s получил предмет: [%s] %s", p.Name, rarity, "Амулет "+item))
}

func logLevelUp(p *Player) {
	exp := rand.Intn(3000) + 1000
	p.Level++
	logMessage("LEVEL", fmt.Sprintf("%s получил %d опыта! (Уровень повышен до %d)", p.Name, exp, p.Level))
}

func logQuest(p *Player) {
	quest := quests[rand.Intn(len(quests))]
	goldReward := (rand.Intn(5) + 1) * p.Level * 10
	logMessage("QUEST", fmt.Sprintf("%s завершил квест '%s' и получил %d золота!", p.Name, quest, goldReward))
}

func logMove(p *Player) {
	location := locations[rand.Intn(len(locations))]
	logMessage("MOVE", fmt.Sprintf("%s переместился в локацию '%s'", p.Name, location))
}

func logTrade(p *Player) {
	gold := rand.Intn(800) + 100
	logMessage("TRADE", fmt.Sprintf("%s продал предмет на аукционе и получил %d золота", p.Name, gold))
}

func logPlayerInteraction(p *Player) {
	otherPlayer := generateFantasyName()
	action := "вызвал на дуэль"
	if rand.Intn(100) > 50 {
		// ИСПРАВЛЕНО: Используем прямой вызов faker.FirstName() для получения строки.
		item := faker.FirstName()
		action = fmt.Sprintf("начал торговлю с игроком %s (предмет: %s)", otherPlayer, "Свиток "+item)
	} else {
		action = fmt.Sprintf("%s игрока %s", action, otherPlayer)
	}
	logMessage("PLAYER", fmt.Sprintf("%s %s", p.Name, action))
}

func logError(p *Player) {
	errors := []string{"Не удалось подключиться к чату гильдии", "Ошибка синхронизации инвентаря", "Задержка сети: 350ms"}
	logMessage("ERROR", fmt.Sprintf("Для игрока %s: %s", p.Name, errors[rand.Intn(len(errors))]))
}

// --- Вспомогательные функции ---

func generateFantasyName() string {
	// Простая генерация фэнтезийных имен
	part1 := []string{"Myst", "Shadow", "Iron", "Dragon", "Silver", "Night", "Fire", "Holy", "Dark", "Forest"}
	part2 := []string{"Blade", "Slayer", "Wizard", "Tank", "Arrow", "Hunter", "Mage", "Paladin", "Rogue", "Ranger"}
	name := fmt.Sprintf("%s%s%d", part1[rand.Intn(len(part1))], part2[rand.Intn(len(part2))], rand.Intn(99))
	return strings.Title(name)
}
