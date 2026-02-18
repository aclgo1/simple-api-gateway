CREATE TABLE grpc_orders(
    order_id UUID PRIMARY KEY NOT NULL,
    account_id UUID NOT NULL,
    products_ids UUID[],
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);