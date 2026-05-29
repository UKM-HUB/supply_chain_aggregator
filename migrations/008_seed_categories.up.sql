INSERT INTO categories (id, name, description) VALUES
    ('food',                  'Food Supplier',          'Food and beverage suppliers'),
    ('packaging',             'Packaging',              'Packaging material and services'),
    ('textile',               'Textile',                'Textile and garment suppliers'),
    ('raw-material',          'Raw Material',           'Raw material suppliers'),
    ('logistics',             'Logistics',              'Logistics and delivery support'),
    ('manufacturing-support', 'Manufacturing Support',  'Supporting vendors for manufacturing operations')
ON CONFLICT (id) DO NOTHING;
