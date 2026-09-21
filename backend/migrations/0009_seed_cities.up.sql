INSERT INTO cities (name, label) VALUES
    ('almaty',    'Almaty'),
    ('astana',    'Astana'),
    ('shymkent',  'Shymkent'),
    ('karaganda', 'Karaganda'),
    ('aktobe',    'Aktobe'),
    ('taraz',     'Taraz'),
    ('pavlodar',  'Pavlodar'),
    ('oskemen',   'Oskemen'),
    ('semey',     'Semey'),
    ('atyrau',    'Atyrau')
ON CONFLICT (name) DO NOTHING;