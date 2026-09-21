INSERT INTO cuisine_types (name, label) VALUES
    ('bakery',     'Bakery'),
    ('cafe',       'Cafe'),
    ('restaurant', 'Restaurant'),
    ('sushi',      'Sushi'),
    ('pizza',      'Pizza'),
    ('italian',    'Italian'),
    ('asian',      'Asian'),
    ('fastfood',   'Fast Food'),
    ('vegetarian', 'Vegetarian'),
    ('desserts',   'Desserts')
ON CONFLICT (name) DO NOTHING;