-- +goose Up
-- +goose StatementBegin

-- Добавляем типы событий
INSERT INTO event_types (name, description)
VALUES ('Расстрел/Расправа', 'Массовые расстрелы и расправы над народом'),
       ('Военная операция', 'Штурмы, осады и военные операции'),
       ('Пограничный конфликт', 'Военные конфликты на границе'),
       ('Восстание/Бунт', 'Народные восстания и бунты');

-- Добавляем тестовую заявку (id сгенерируется автоматически = 1)
INSERT INTO requests (title, description, event_date, event_type_id, email, telegram_username, archive_url, site_url, status)
VALUES ('Тестовая заявка для seed данных', 'Заявка для загрузки тестовых исторических объектов', '2024-01-01', 1,
        'test@example.com', 'testuser', 'https://example.com/archive.zip', 'https://example.com/test-site', 'Опубликована');

-- Добавляем тестовые объекты
INSERT INTO objects (request_id, title, description, coordinates, event_date, event_type_id, site_url, preview_image_url)
VALUES (1, 'Бородинское сражение',
        'Крупнейшая битва Отечественной войны 1812 года между русской и французской армиями.',
        ST_GeogFromText('POINT(35.8167 55.5167)'), '1812-09-07', 1, 'https://example.com/borodino', 'https://example.com/images/borodino.jpg'),
       (1, 'Взятие Бастилии', 'Начало Великой французской революции, штурм крепости-тюрьмы.',
        ST_GeogFromText('POINT(2.3698 48.8531)'), '1789-07-14', 4, 'https://example.com/bastille', 'https://example.com/images/bastille.jpg'),
       (1, 'Потсдамская конференция', 'Встреча лидеров трех держав для определения послевоенного устройства Европы.',
        ST_GeogFromText('POINT(13.0475 52.4)'), '1945-07-17', 1, 'https://example.com/potsdam', 'https://example.com/images/potsdam.jpg'),
       (1, 'Кровавое воскресенье',
        'Расстрел мирной демонстрации рабочих в Петербурге, шедших с петицией к царю. Стало началом первой русской революции.',
        ST_GeogFromText('POINT(30.3141 59.9398)'), '1905-01-09', 2, 'https://example.com/bloody-sunday', 'https://example.com/bloody-sunday.jpg'),
       (1, 'Ходынская трагедия',
        'Массовая давка на Ходынском поле в Москве во время раздачи царских подарков на коронации Николая II.',
        ST_GeogFromText('POINT(37.5405 55.7897)'), '1896-05-30', 1, 'https://example.com/khodynka', 'https://example.com/khodynka.jpg'),
       (1, 'Новочеркасский расстрел',
        'Расстрел демонстрации рабочих Новочеркасского электровозостроительного завода в ответ на повышение цен.',
        ST_GeogFromText('POINT(40.0933 47.4214)'), '1962-06-02', 3, 'https://example.com/novocherkassk', 'https://example.com/novocherkassk.jpg'),
       (1, 'Восстание декабристов',
        'Вооружённое восстание офицеров на Сенатской площади против самодержавия и крепостного права.',
        ST_GeogFromText('POINT(30.3063 59.9341)'), '1825-12-26', 2, 'https://example.com/decembrists', 'https://example.com/decembrists.jpg'),
       (1, 'Соляной бунт', 'Народное восстание в Москве против повышения налога на соль и произвола бояр.',
        ST_GeogFromText('POINT(37.6173 55.7558)'), '1648-06-01', 1, 'https://example.com/salt-riot', 'https://example.com/salt-riot.jpg'),
       (1, 'Медный бунт', 'Восстание в Москве из-за обесценивания медных денег и финансового кризиса.',
        ST_GeogFromText('POINT(37.6173 55.7558)'), '1662-07-25', 4, 'https://example.com/copper-riot', 'https://example.com/copper-riot.jpg'),
       (1, 'Астраханское восстание', 'Народное восстание стрельцов и посадских людей против воеводы и приказных.',
        ST_GeogFromText('POINT(48.04 46.3497)'), '1705-07-30', 4, 'https://example.com/astrakhan-uprising', 'https://example.com/astrakhan-uprising.jpg'),
       (1, 'Чумной бунт', 'Восстание в Москве против карантинных мер во время эпидемии чумы.',
        ST_GeogFromText('POINT(37.6173 55.7558)'), '1771-09-15', 4, 'https://example.com/plague-riot', 'https://example.com/plague-riot.jpg'),
       (1, 'Кронштадтское восстание',
        'Антибольшевистское восстание матросов и красноармейцев в Кронштадте под лозунгом "Советы без коммунистов".',
        ST_GeogFromText('POINT(29.77 60.008)'), '1921-03-01', 4, 'https://example.com/kronstadt', 'https://example.com/kronstadt.jpg'),
       (1, 'Тамбовское восстание',
        'Крупнейшее крестьянское восстание против политики военного коммунизма и продразвёрстки.',
        ST_GeogFromText('POINT(41.4446 52.7213)'), '1920-08-19', 1, 'https://example.com/tambov', 'https://example.com/tambov.jpg'),
       (1, 'Взятие Казани', 'Штурм и взятие Казани войсками Ивана Грозного, присоединение Казанского ханства к России.',
        ST_GeogFromText('POINT(49.1055 55.7985)'), '1552-10-02', 2, 'https://example.com/kazan-siege', 'https://example.com/kazan-siege.jpg'),
       (1, 'Штурм Измаила', 'Взятие турецкой крепости Измаил войсками Суворова во время русско-турецкой войны.',
        ST_GeogFromText('POINT(28.84 45.3567)'), '1790-12-22', 2, 'https://example.com/izmail', 'https://example.com/izmail.jpg'),
       (1, 'Оборона Севастополя',
        'Героическая 349-дневная оборона Севастополя во время Крымской войны против англо-франко-турецких войск.',
        ST_GeogFromText('POINT(33.5244 44.6178)'), '1854-10-13', 2, 'https://example.com/sevastopol-defense', 'https://example.com/sevastopol-defense.jpg'),
       (1, 'Блокада Ленинграда', 'Начало блокады Ленинграда немецкими войсками, длившейся 872 дня.',
        ST_GeogFromText('POINT(30.3141 59.9398)'), '1941-09-08', 2, 'https://example.com/leningrad-blockade', 'https://example.com/leningrad-blockade.jpg'),
       (1, 'Штурм Берлина', 'Финальная операция Великой Отечественной войны, взятие Берлина советскими войсками.',
        ST_GeogFromText('POINT(13.405 52.52)'), '1945-04-16', 2, 'https://example.com/berlin-assault', 'https://example.com/berlin-assault.jpg'),
       (1, 'Бой у озера Хасан', 'Вооруженный конфликт между СССР и Японией у озера Хасан на границе с Маньчжурией.',
        ST_GeogFromText('POINT(131.6 42.6833)'), '1938-07-29', 3, 'https://example.com/khasan', 'https://example.com/khasan.jpg'),
       (1, 'Бои на Халхин-Голе', 'Военный конфликт между СССР и Японией на монгольско-маньчжурской границе.',
        ST_GeogFromText('POINT(118.25 47.7333)'), '1939-05-11', 3, 'https://example.com/khalkhin-gol', 'https://example.com/khalkhin-gol.jpg'),
       (1, 'Даманский конфликт',
        'Пограничный вооруженный конфликт между СССР и КНР на острове Даманский на реке Уссури.',
        ST_GeogFromText('POINT(133.7833 46.4833)'), '1969-03-02', 3, 'https://example.com/damansky', 'https://example.com/damansky.jpg'),
       (1, 'Инцидент у острова Жаланашколь',
        'Пограничное столкновение между советскими и китайскими войсками в районе озера Жаланашколь.',
        ST_GeogFromText('POINT(79.0833 45.1167)'), '1969-08-13', 3, 'https://example.com/zhalanashkol', 'https://example.com/zhalanashkol.jpg'),
       (1, 'Бой у острова Даманский', 'Второе столкновение советских пограничников с китайскими военными на Даманском.',
        ST_GeogFromText('POINT(133.7833 46.4833)'), '1969-03-15', 3, 'https://example.com/damansky-2', 'https://example.com/damansky-2.jpg');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE objects CASCADE;
TRUNCATE TABLE requests CASCADE;
TRUNCATE TABLE event_types CASCADE;
-- +goose StatementEnd