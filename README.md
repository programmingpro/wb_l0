контейнеры запускаются при помощи docker compose up --build
go run main.go - запуск сервиса
веб приложение доступно по http://localhost:7000/

в папке messages находится test.json - массив сообщений

при запуске сервиса producer закинет в канал массив сообщений
