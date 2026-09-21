INSERT INTO product_categories (name, label) VALUES
    ('surprise_box', 'Surprise Box'),
    ('burgers',      'Burgers'),
    ('pizza',        'Pizza'),
    ('sushi',        'Sushi'),
    ('drinks',       'Drinks'),
    ('salads',       'Salads'),
    ('soups',        'Soups'),
    ('desserts',     'Desserts'),
    ('snacks',       'Snacks'),
    ('breakfast',    'Breakfast'),
    ('sets',         'Sets')
ON CONFLICT (name) DO NOTHING;