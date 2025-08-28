#!/bin/bash

# Проверяем наличие ffmpeg
if ! command -v ffmpeg &>/dev/null; then
	echo "Ошибка: ffmpeg не установлен. Установите его сначала."
	echo "Для Ubuntu/Debian: sudo apt install ffmpeg"
	echo "Для CentOS/RHEL: sudo yum install ffmpeg"
	exit 1
fi

# Создаем временный файл со списком
list_file=$(mktemp)

# Находим все .ogg файлы в папке audio, сортируем по имени и записываем в список
find audio -name "*.ogg" -type f | sort | while read -r file; do
	echo "file '$file'" >>"$list_file"
done

# Проверяем, есть ли файлы для обработки
if [ ! -s "$list_file" ]; then
	echo "Не найдено .ogg файлов в папке audio"
	rm "$list_file"
	exit 1
fi

# Объединяем и конвертируем в mp3
ffmpeg -f concat -safe 0 -i "$list_file" -c:a libmp3lame -q:a 2 -y output.mp3

# Удаляем временный файл
rm "$list_file"

echo "Готово! Все файлы объединены в output.mp3"
