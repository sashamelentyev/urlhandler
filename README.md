# urlhandler

Необходимо реализовать CLI-утилиту, которая реализует асинхронную обработку входящих URL из файла, переданного в качестве аргумента данной утилите.
Формат входного файла: на каждой строке – один URL. URL может быть очень много! Но могут быть и невалидные URL.

Пример входного файла:
https://myoffice.ru
https://yandex.ru

По каждому URL получить контент и вывести в консоль его размер и время обработки. Предусмотреть обработку ошибок.


## Устновка
```bash
go install github.com/sashamelentyev/urlhandler@latest
```

## Запуск
```bash
urlhandler -request-timeout [sec] -filename [file]
```

## Пример использования

Файл urls.txt
```txt title='urls.txt'
https://myoffice.ru
https://yandex.ru
test://
```

Запуск проверки URL
```bash
urlhandler -request-timeout 2s -filename urls.txt
```

Результат
```
error: do request to "test:": Get "test:": unsupported protocol scheme "test"
url: https://myoffice.ru, Content length: -1, Response time: 117.588208ms
url: https://yandex.ru, Content length: 15672, Response time: 214.930167ms
```
