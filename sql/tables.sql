create table accounts (
  id serial primary key,
  account_id varchar(128) unique,
  currency varchar(32) not null,
  balance numeric not null,
  
  constraint balance_non_negative check (balance >= 0.0)
);

create table payments (
  id uuid primary key,
  source_id integer not null,
  destination_id integer not null,
  amount numeric not null,
  
  constraint payments_source_fk foreign key (source_id)
    references accounts (id) match simple
    on update no action
    on delete no action,
  constraint payments_destination_fk foreign key (destination_id)
    references accounts (id) match simple
    on update no action
    on delete no action
)