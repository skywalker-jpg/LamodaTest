BEGIN;

CREATE TABLE IF NOT EXISTS warehouses (
                                          id SERIAL PRIMARY KEY,
                                          name VARCHAR(255) NOT NULL,
                                          availability BOOLEAN NOT NULL
    );

CREATE TABLE IF NOT EXISTS products (
                                        id SERIAL PRIMARY KEY,
                                        name VARCHAR(255) NOT NULL,
                                        size VARCHAR(50) NOT NULL,
                                        code VARCHAR(50) UNIQUE NOT NULL
    );

CREATE TABLE IF NOT EXISTS warehouse_product (
                                                 id SERIAL PRIMARY KEY,
                                                 warehouse_id INT REFERENCES warehouses(id) ON DELETE CASCADE,
                                                 product_id INT REFERENCES products(id) ON DELETE CASCADE,
                                                 UNIQUE (warehouse_id, product_id),

                                                 quantity INT NOT NULL DEFAULT 0,
                                                 reserved_quantity INT NOT NULL DEFAULT 0
    );


COMMIT;