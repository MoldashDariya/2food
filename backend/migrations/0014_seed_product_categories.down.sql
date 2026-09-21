DELETE FROM product_categories
WHERE name IN (
    'surprise_box', 'burgers', 'pizza', 'sushi', 'drinks',
    'salads', 'soups', 'desserts', 'snacks',
    'breakfast', 'sets'
);