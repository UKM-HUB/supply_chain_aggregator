DELETE FROM categories
WHERE id IN (
    'food',
    'packaging',
    'textile',
    'raw-material',
    'logistics',
    'manufacturing-support'
);
