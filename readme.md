# Zarg

Бот сообщества ВКонтакте для игры в беседе.

## Поддерживаемые команды

- `в поход [<N>]` - начинается новый поход в данж. Максимальное количество игроков ограничивается числом `N` (`4` по умолчанию).
- `статы [<имя>]` - показывает статистику конкретного пользователя или того, кто вызвал команду.

---

## Процесс игры

Игра начинается с набора игроков в группу командой `в поход`. Кто хочет участвовать в походе, пишет `я`. Если по окончанию времени набора никто не захотел участвовать, игра отменяется.

Далее игрокам предлагается выбрать стартовое оружие. Два игрока не могут выбрать одно и то же оружие, поэтому кто успел... Если вообще не успел ничего взять, дерешься собственными кулаками.

После игрокам предлагается определить очередность ходов в боёвках. Для этого первый игрок должен написать `я`, а остальные - `потом я`. Тут игроки должны договориться между собой сами, ограничений по времени нет (но 5 минут максимум - потом отмена игры).

*И происходит спуск в подземелье!*

Первая локация игры - подземелье. Количество комнат от 5 до 8.

---

## Локации

Локация состоит из комнат, которые игроки проходят последовательно. В каждой комнате могут быть либо враги, либо сокровища, либо ловушки, либо комната может быть удобной для привала.

### Привал

Восстанавливает все здоровье игрокам и позволяет сменить очередность в боёвках.

### Сокровища

Игроки находят гору сокровищ, которые разбираются по принципу кто успел...

### Ловушки

Случайным образом выбирается игрок, которому наносится 10-20 урона.

### Враги

Боёвка начинается с хода первого игрока. После этого ходит враг. Потом второй игрок, и так далее. Частоты ходов игроков и врагов подбираются таким образом, чтобы игроки и враги успели походить как минимум по 1 разу.

Варианты, из которых может выбрать игрок:

1. Атаковать своим оружием.
2. Поставить блок (следующая атака врага ослабляется но идет по этому игроку).
3. Использовать предмет из инвентаря.

Если по игроку наносится урон больший его текущего здоровья, он умирает. После завершения боёвки живым игрокам разрешается забрать вещи павших товарищей. (TODO: но за каждую вещь повышается уровень проклятья и играть становится сложнее).

---

