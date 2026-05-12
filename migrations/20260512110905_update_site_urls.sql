
-- +goose Up
UPDATE requests
SET site_url        = REPLACE(site_url,        'api.russia-heroes.ru/sites/', 'russia-heroes.ru/sites/'),
    screenshot_url  = REPLACE(screenshot_url,  'api.russia-heroes.ru/sites/', 'russia-heroes.ru/sites/')
WHERE site_url LIKE '%api.russia-heroes.ru/sites/%';

UPDATE objects
SET site_url          = REPLACE(site_url,          'api.russia-heroes.ru/sites/', 'russia-heroes.ru/sites/'),
    preview_image_url = REPLACE(preview_image_url, 'api.russia-heroes.ru/sites/', 'russia-heroes.ru/sites/')
WHERE site_url LIKE '%api.russia-heroes.ru/sites/%';

-- +goose Down
UPDATE requests
SET site_url        = REPLACE(site_url,        'russia-heroes.ru/sites/', 'api.russia-heroes.ru/sites/'),
    screenshot_url  = REPLACE(screenshot_url,  'russia-heroes.ru/sites/', 'api.russia-heroes.ru/sites/')
WHERE site_url LIKE '%russia-heroes.ru/sites/%';

UPDATE objects
SET site_url          = REPLACE(site_url,          'russia-heroes.ru/sites/', 'api.russia-heroes.ru/sites/'),
    preview_image_url = REPLACE(preview_image_url, 'russia-heroes.ru/sites/', 'api.russia-heroes.ru/sites/')
WHERE site_url LIKE '%russia-heroes.ru/sites/%';
