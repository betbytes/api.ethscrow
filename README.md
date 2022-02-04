# ethscrow

## `/broker`

- `POST /create`
  - Body:
    ```json
    {
      "caller_username": "string",
      "mediator_username": "string - optional",
      "reason": "string"
    }
    ```
    
- `GET /connect/{poolId}` - if all conditions met, it'll be upgraded to a websocket

## `/user`

- 

## SQL for table creation

```sql
CREATE TABLE users(
    username varchar(20) not null primary key,
    public_key char(182) not null unique,
    enc_public_key char(182) not null,
    email varchar(50),
    created_at timestamp default now()
);

CREATE TABLE pools(
    id char(32) not null primary key,
    address varchar(42),
    mediator_username varchar(20) references users(username),
    bettor_username varchar(20) references users(username),
    caller_username varchar(20) references users(username),
    bettor_state smallint default 0 check ( bettor_state >= -1 AND bettor_state <= 1 ),
    caller_state smallint default 0 check ( caller_state >= -1 AND caller_state <= 1 ),
    threshold_key varchar(300),
    created_at timestamp default now(),
    reason varchar(200) not null,
    balance decimal default 0 not null,
    balance_last_updated timestamp,
    accepted bool default false,
    initialized bool default false not null
);

CREATE TABLE chats(
    id char(32) not null primary key,
    pool_id char(32) not null,
    message varchar(150) not null,
    timestamp timestamp default now(),
    from_username varchar(20) not null references users(username)
);
```