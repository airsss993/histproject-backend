
-- +goose Up
UPDATE requests
SET site_url = REPLACE(site_url, 'api.russia-heroes.ru/sites/', 'russia-heroes.ru/sites/')
WHERE site_url LIKE '%api.russia-heroes.ru/sites/%';

UPDATE objects
SET site_url = REPLACE(site_url, 'api.russia-heroes.ru/sites/', 'russia-heroes.ru/sites/')
WHERE site_url LIKE '%api.russia-heroes.ru/sites/%';

-- +goose Down
UPDATE requests
SET site_url = REPLACE(site_url, 'russia-heroes.ru/sites/', 'api.russia-heroes.ru/sites/')
WHERE site_url LIKE '%russia-heroes.ru/sites/%';

UPDATE objects
SET site_url = REPLACE(site_url, 'russia-heroes.ru/sites/', 'api.russia-heroes.ru/sites/')
WHERE site_url LIKE '%russia-heroes.ru/sites/%';
