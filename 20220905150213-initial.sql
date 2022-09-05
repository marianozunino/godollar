-- +migrate Up

create table if not exists "rates" (
		date timestamp primary key,
		buy_price float,
		sell_price float,
		ebrou_buy_price float,
		ebrou_sell_price float
);

-- add index to date column
create index if not exists "rates_date_idx" on "rates" (date);



-- +migrate Down
drop table if exists "rates";

