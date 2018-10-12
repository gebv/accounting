CREATE SCHEMA IF NOT EXISTS acca;

CREATE EXTENSION ltree;

-- money in the numeric(69, 0)
-- for example, to store balances for WEI (ETH)
-- change manually if you need less accuracy to
-- NOTE: cannot be justified if used small number (<1^10*12), recomented used bigint

-- currencies
-- currencies any format
CREATE TABLE acca.currencies (
    curr_id bigserial PRIMARY KEY,
    key ltree NOT NULL,
    meta jsonb NOT NULL DEFAULT '{}'
);
CREATE INDEX currencies_key_gist_idx ON acca.currencies USING GIST (key);

COMMENT ON COLUMN acca.currencies.curr_id IS 'Currency ID.';
COMMENT ON COLUMN acca.currencies.key IS 'Currency key (it is not primary key).';
COMMENT ON COLUMN acca.currencies.meta IS 'Container with meta information.';

-- accounts
-- account and related meta information and current balance
CREATE TABLE acca.accounts (
    acc_id bigserial PRIMARY KEY,
    curr_id bigint REFERENCES acca.currencies(curr_id),
    key ltree NOT NULL,
    balance numeric(69, 00) NOT NULL DEFAULT 0 CHECK(balance >= 0),
    meta jsonb NOT NULL DEFAULT '{}'
);

COMMENT ON COLUMN acca.accounts.acc_id IS 'Account ID.';
COMMENT ON COLUMN acca.accounts.curr_id IS 'Currency of account.';
COMMENT ON COLUMN acca.accounts.key IS 'Account key (it is not primary key).';
COMMENT ON COLUMN acca.accounts.balance IS 'Current balance.';
COMMENT ON COLUMN acca.accounts.meta IS 'Container with meta information.';


-- transaction status
-- diagram of the transition of the statuses (see TODO: link to WIKI)
CREATE TYPE acca.transaction_status AS enum (
    'unknown',
    'draft',
    'auth',
    'accepted',
    'rejected',

    'failed'
);

-- ALTER TYPE  acca.transaction_status ADD VALUE 'failed';

-- transactions
-- transactions and related meta information and current status
CREATE TABLE acca.transactions (
    tx_id bigserial PRIMARY KEY,
    reason ltree NOT NULL,
    meta jsonb NOT NULL DEFAULT '{}',
    status acca.transaction_status NOT NULL DEFAULT 'unknown' CHECK (status <> 'unknown'),
    errm text,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone
);
CREATE INDEX transactions_reason_gist_idx ON acca.transactions USING GIST (reason);

COMMENT ON COLUMN acca.transactions.tx_id IS 'Transaction ID.';
COMMENT ON COLUMN acca.transactions.reason IS 'The reason for the transfer.';
COMMENT ON COLUMN acca.transactions.meta IS 'Container with meta information.';
COMMENT ON COLUMN acca.transactions.status IS 'Transaction status.';

-- type of operation
-- - internal - transfer between accounts
-- - recharge - entering funds into the system
-- - withdraw - output funds from the system
CREATE TYPE acca.operation_type AS enum (
    'unknown',
    'internal',
    'recharge',
    'withdraw'
);

-- status of operation
-- diagram of the transition of the statuses (see TODO: link to WIKI)
CREATE TYPE acca.operation_status AS enum (
    'unknown',
    'draft',
    'hold',
    'accepted',
    'rejected'
);

-- operation included in the transaction
CREATE TABLE acca.operations (
    oper_id bigserial PRIMARY KEY,
    tx_id bigint NOT NULL REFERENCES acca.transactions (tx_id),
    src_acc_id bigint NOT NULL REFERENCES acca.accounts(acc_id),
    dst_acc_id bigint NOT NULL REFERENCES acca.accounts(acc_id),
    type acca.operation_type NOT NULL DEFAULT 'unknown' CHECK (type <> 'unknown'),
    amount numeric(69, 00) NOT NULL,
    reason ltree NOT NULL DEFAULT '',
    meta jsonb NOT NULL DEFAULT '{}',
    hold boolean NOT NULL DEFAULT false,
    hold_acc_id bigint REFERENCES acca.accounts(acc_id),
    status acca.operation_status NOT NULL DEFAULT 'unknown' CHECK (type <> 'unknown'),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone
);
CREATE INDEX operations_reason_gist_idx ON acca.operations USING GIST (reason);
CREATE INDEX operations_tx_idx ON acca.operations USING GIST (tx_id);

COMMENT ON COLUMN acca.operations.src_acc_id IS 'Withdrawal account.';
COMMENT ON COLUMN acca.operations.dst_acc_id IS 'Deposit account.';
COMMENT ON COLUMN acca.operations.type IS 'Type of operation.';
COMMENT ON COLUMN acca.operations.amount IS 'Transaction amount.';
COMMENT ON COLUMN acca.operations.reason IS 'The reason for the operation.';
COMMENT ON COLUMN acca.operations.meta IS 'Container with meta information.';
COMMENT ON COLUMN acca.operations.hold IS 'If true, the translation is two-step.';
COMMENT ON COLUMN acca.operations.hold_acc_id IS 'Suspense account. Only for two-step transaction. May be NULL.';
COMMENT ON COLUMN acca.operations.status IS 'Operation status.';

CREATE TYPE acca.request_type AS enum (
    'unknown',
    'auth',
    'accept',
    'reject',

    'rollback'
);

-- request for action
CREATE TABLE acca.requests_queue (
    tx_id bigserial REFERENCES acca.transactions (tx_id),
    type request_type NOT NULL DEFAULT 'unknown' CHECK (type <> 'unknown'),
    created_at timestamp without time zone NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX uniq_request_type_for_tx_idx ON acca.requests_queue (tx_id, type);

-- history requests
CREATE TABLE acca.requests_history (
    tx_id bigserial REFERENCES acca.transactions (tx_id),
    type request_type NOT NULL DEFAULT 'unknown' CHECK (type <> 'unknown'),
    created_at timestamp without time zone NOT NULL,
    executed_at timestamp without time zone NOT NULL DEFAULT now()
);

-- TODO: balance_changes

