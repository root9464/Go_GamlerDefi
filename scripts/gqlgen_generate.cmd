@echo off
cd /d %~dp0..

if "%1"=="" (
    echo Использование: generate.cmd [модуль1] [модуль2] ...
    echo Пример: generate.cmd test auth user
    exit /b 1
)

:loop
if "%1"=="" goto end
echo Генерация для модуля %1...

powershell -Command "(Get-Content gqlgen.yml) -replace '{module}', '%1' | Set-Content gqlgen_%1.yml"

go run github.com/99designs/gqlgen generate --config gqlgen_%1.yml
if %errorlevel% neq 0 (
    echo Ошибка при генерации кода для модуля %1
    del gqlgen_%1.yml
    exit /b 1
)

del gqlgen_%1.yml
echo Генерация для модуля %1 завершена успешно

shift
goto loop

:end
echo Генерация кода завершена для всех указанных модулей