INSERT INTO warehouses (name, availability) VALUES
                                                ('Warehouse A', true),
                                                ('Warehouse B', false),
                                                ('Warehouse C', true),
                                                ('Warehouse D', true);

INSERT INTO products (name, size, code) VALUES
                                            ('Product 1', 'Small', '123'),
                                            ('Product 2', 'Medium', '456'),
                                            ('Product 3', 'Large', '789'),
                                            ('Product 4', 'Small', '444');

INSERT INTO warehouse_product (warehouse_id, product_id, quantity, reserved_quantity) VALUES
    (1, 1, 100, 20),
    (1, 2, 50, 10),
    (2, 3, 75, 15);
