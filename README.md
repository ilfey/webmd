# WEBmd

Сервак, который рендерит md и отдает клиенту.

# Запуск

## Зависимости:
- golang ^1.21.0
- nodejs ^14.17.0 + yarn

## Eсли есть task:

```bash
task build
webmd.exe
```

## Если нет (команды под линуху)

Если нет линухи юзаните `git bash`

```bash
mkdir -p .dist
mkdir -p .dist/css
mkdir -p .dist/font
mkdir -p .dist/img
mkdir -p .dist/js

# билдим ассеты
cd assets
yarn install

# билдим css
npx tailwindcss -m -c tailwind.config.js -o ../.dist/css/tailwind.min.css
npx esbuild --minify --external:*.woff2 --bundle ./src/index.css --outfile=../.dist/css/bundle.min.css

# копируем шрифты
cp -r ./src/font ../.dist/font

# копируем изображения
cp -r ./src/img ../.dist/img

# билдим жс
cat js/*.js | npx esbuild --minify > ../.dist/js/bundle.min.js
```