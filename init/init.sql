CREATE TABLE IF NOT EXISTS orders (
                                      order_id        BIGSERIAL PRIMARY KEY,
                                      order_uid       VARCHAR(64) UNIQUE NOT NULL,
    track_number    VARCHAR(64) NOT NULL,
    entry           VARCHAR(16) NOT NULL,
    locale          VARCHAR(8) NOT NULL,
    internal_signature VARCHAR(255),
    customer_id     VARCHAR(64) NOT NULL,
    delivery_service VARCHAR(64) NOT NULL,
    shardkey        VARCHAR(16),
    sm_id           INT NOT NULL,
    date_created    TIMESTAMPTZ NOT NULL DEFAULT now(),
    oof_shard       VARCHAR(16)
    );

CREATE TABLE IF NOT EXISTS deliveries (
                                          delivery_id BIGSERIAL PRIMARY KEY,
                                          order_id    BIGINT NOT NULL UNIQUE REFERENCES orders(order_id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    phone       VARCHAR(32) NOT NULL,
    zip         VARCHAR(16) NOT NULL,
    city        VARCHAR(128) NOT NULL,
    address     VARCHAR(255) NOT NULL,
    region      VARCHAR(128) NOT NULL,
    email       VARCHAR(128) NOT NULL
    );

CREATE TABLE IF NOT EXISTS payments (
                                        payment_id  BIGSERIAL PRIMARY KEY,
                                        order_id    BIGINT NOT NULL UNIQUE REFERENCES orders(order_id) ON DELETE CASCADE,
    transaction VARCHAR(64) NOT NULL,
    request_id  VARCHAR(64),
    currency    VARCHAR(8) NOT NULL,
    provider    VARCHAR(64) NOT NULL,
    amount      DECIMAL(12,2) NOT NULL,
    payment_dt  BIGINT NOT NULL,
    bank        VARCHAR(64) NOT NULL,
    delivery_cost DECIMAL(12,2) NOT NULL,
    goods_total DECIMAL(12,2) NOT NULL,
    custom_fee  DECIMAL(12,2) DEFAULT 0
    );

CREATE TABLE IF NOT EXISTS items (
                                     item_id     BIGSERIAL PRIMARY KEY,
                                     order_id    BIGINT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    chrt_id     BIGINT NOT NULL,
    track_number VARCHAR(64) NOT NULL,
    price       DECIMAL(12,2) NOT NULL,
    rid         VARCHAR(64),
    name        VARCHAR(255) NOT NULL,
    sale        INT,
    size        VARCHAR(32),
    total_price DECIMAL(12,2) NOT NULL,
    nm_id       BIGINT NOT NULL,
    brand       VARCHAR(128),
    status      INT NOT NULL
    );
