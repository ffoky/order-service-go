-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS deliveries
(
    delivery_id BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    phone       VARCHAR(32)  NOT NULL,
    zip         VARCHAR(16)  NOT NULL,
    city        VARCHAR(128) NOT NULL,
    address     VARCHAR(255) NOT NULL,
    region      VARCHAR(128) NOT NULL,
    email       VARCHAR(128) NOT NULL
    );

CREATE TABLE IF NOT EXISTS payments
(
    transaction   VARCHAR(64) PRIMARY KEY,
    request_id    VARCHAR(64),
    currency      VARCHAR(8)     NOT NULL,
    provider      VARCHAR(64)    NOT NULL,
    amount        DECIMAL(12, 2) NOT NULL,
    payment_dt    BIGINT         NOT NULL,
    bank          VARCHAR(64)    NOT NULL,
    delivery_cost DECIMAL(12, 2) NOT NULL,
    goods_total   DECIMAL(12, 2) NOT NULL,
    custom_fee    DECIMAL(12, 2) DEFAULT 0
    );

CREATE TABLE IF NOT EXISTS items
(
    chrt_id      INT PRIMARY KEY,
    track_number VARCHAR(64)    NOT NULL,
    price        DECIMAL(12, 2) NOT NULL,
    rid          VARCHAR(64),
    name         VARCHAR(255)   NOT NULL,
    sale         INT,
    size         VARCHAR(32)    NOT NULL,
    total_price  DECIMAL(12, 2) NOT NULL,
    nm_id        BIGINT         NOT NULL,
    brand        VARCHAR(128)   NOT NULL,
    status       INT            NOT NULL
    );

CREATE TABLE IF NOT EXISTS orders
(
    order_uid          VARCHAR(64) PRIMARY KEY,
    delivery_id        BIGINT REFERENCES deliveries (delivery_id)       NOT NULL,
    payment_id         VARCHAR(64) REFERENCES payments (transaction) NOT NULL,
    track_number       VARCHAR(64)                                   NOT NULL,
    entry              VARCHAR(16)                                   NOT NULL,
    locale             VARCHAR(8)                                    NOT NULL,
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(64)                                   NOT NULL,
    delivery_service   VARCHAR(64)                                   NOT NULL,
    shardkey           VARCHAR(16),
    sm_id              INT                                           NOT NULL,
    date_created       TIMESTAMPTZ                                   NOT NULL,
    oof_shard          VARCHAR(16)
    );

CREATE TABLE IF NOT EXISTS order_items
(
    order_uid    VARCHAR(64) REFERENCES orders (order_uid) NOT NULL,
    item_chrt_id INTEGER REFERENCES items (chrt_id)         NOT NULL,
    PRIMARY KEY (order_uid, item_chrt_id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;
-- +goose StatementEnd