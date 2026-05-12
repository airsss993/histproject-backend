
-- +goose Up
-- Опубликованные заявки и объекты на карте → публичный /sites/
UPDATE requests
SET site_url = REPLACE(site_url, 'api.russia-heroes.ru/sites/', 'russia-heroes.ru/sites/')
WHERE site_url LIKE '%api.russia-heroes.ru/sites/%'
  AND status = 'Опубликована';

UPDATE objects
SET site_url = REPLACE(site_url, 'api.russia-heroes.ru/sites/', 'russia-heroes.ru/sites/')
WHERE site_url LIKE '%api.russia-heroes.ru/sites/%';

-- Неопубликованные заявки (Новая, На проверке, Отклонена) → приватный /preview/
UPDATE requests
SET site_url = REPLACE(site_url, 'api.russia-heroes.ru/sites/', 'russia-heroes.ru/preview/')
WHERE site_url LIKE '%api.russia-heroes.ru/sites/%'
  AND status != 'Опубликована';

-- +goose Down
UPDATE requests
SET site_url = REPLACE(site_url, 'russia-heroes.ru/sites/', 'api.russia-heroes.ru/sites/')
WHERE site_url LIKE '%russia-heroes.ru/sites/%';

UPDATE requests
SET site_url = REPLACE(site_url, 'russia-heroes.ru/preview/', 'api.russia-heroes.ru/sites/')
WHERE site_url LIKE '%russia-heroes.ru/preview/%';

UPDATE objects
SET site_url = REPLACE(site_url, 'russia-heroes.ru/sites/', 'api.russia-heroes.ru/sites/')
WHERE site_url LIKE '%russia-heroes.ru/sites/%';
