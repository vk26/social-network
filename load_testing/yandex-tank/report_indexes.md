#Without indexes
```sql 
EXPLAIN FORMAT=json Select * from users WHERE name LIKE 'jo%' OR surname LIKE 'jo%';
```
```json
{
  "query_block": {
    "select_id": 1,
    "cost_info": {
      "query_cost": "214601.60"
    },
    "table": {
      "table_name": "users",
      "access_type": "ALL",
      "rows_examined_per_scan": 984768,
      "rows_produced_per_join": 206660,
      "filtered": "20.99",
      "cost_info": {
        "read_cost": "173269.55",
        "eval_cost": "41332.05",
        "prefix_cost": "214601.60",
        "data_read_per_join": "18M"
      },
      "used_columns": [
        "id",
        "name",
        "surname",
        "birthday",
        "city",
        "about",
        "avatar",
        "email",
        "password_hash",
        "created_at",
        "updated_at"
      ],
      "attached_condition": "((`social_dev`.`users`.`name` like 'jo%') or (`social_dev`.`users`.`surname` like 'jo%'))"
    }
  }
}
```
Load testing result by yandex-tank:
https://clck.ru/M8F4c

I faced with problem "too many files open" and "too many connections"
I did this tuning:
In Mysql side set:
max_connections = 1000
open_files_limit = 65000
In OS Linux:
change /etc/security/limits.conf
```bash
* soft     nproc          65535
* hard     nproc          65535
* soft     nofile         65535
* hard     nofile         65535
```
In application side:
DB.SetMaxOpenConns(900)
DB.SetMaxIdleConns(100)
srv.ReadTimeout:  time.Second * 15

#With BTREE indexes
Try add index BTREE for name, surname fields. And watch explain again. We use btree because it support LIKE clauses.
```sql
CREATE INDEX name_idx ON users(name(15));
CREATE INDEX surname_idx ON users(surname(15));
EXPLAIN FORMAT=json Select * from users WHERE name LIKE 'jo%' OR surname LIKE 'jo%';
```
```json
{
  "query_block": {
    "select_id": 1,
    "cost_info": {
      "query_cost": "122099.48"
    },
    "table": {
      "table_name": "users",
      "access_type": "index_merge",
      "possible_keys": [
        "name_idx",
        "surname_idx"
      ],
      "key": "sort_union(name_idx,surname_idx)",
      "key_length": "17,17",
      "rows_examined_per_scan": 46759,
      "rows_produced_per_join": 46759,
      "filtered": "100.00",
      "cost_info": {
        "read_cost": "112747.69",
        "eval_cost": "9351.80",
        "prefix_cost": "122099.49",
        "data_read_per_join": "4M"
      },
      "used_columns": [
        "id",
        "name",
        "surname",
        "birthday",
        "city",
        "about",
        "avatar",
        "email",
        "password_hash",
        "created_at",
        "updated_at"
      ],
      "attached_condition": "((`social_dev`.`users`.`name` like 'jo%') or (`social_dev`.`users`.`surname` like 'jo%'))"
    }
  }
} 
```
Load testing with BTREE by yandex-tank:
https://clck.ru/M8Ec9

#With FULLTEXT indexes
And let's try to use FULLTEXT index.
```sql
drop index name_idx on users;
drop index surname_idx on users;
CREATE FULLTEXT INDEX fulltext_name_idx ON users(name, surname);
```
Now we have to change our search query to use FULLTEXT index:
```sql
EXPLAIN FORMAT=JSON Select * from users WHERE MATCH(name, surname) AGAINST ('+jo*' IN BOOLEAN MODE);
```
```json
{
  "query_block": {
    "select_id": 1,
    "cost_info": {
      "query_cost": "1.20"
    },
    "table": {
      "table_name": "users",
      "access_type": "fulltext",
      "possible_keys": [
        "fulltext_name_idx"
      ],
      "key": "fulltext_name_idx",
      "used_key_parts": [
        "name"
      ],
      "key_length": "0",
      "ref": [
        "const"
      ],
      "rows_examined_per_scan": 1,
      "rows_produced_per_join": 1,
      "filtered": "100.00",
      "ft_hints": "no_ranking",
      "cost_info": {
        "read_cost": "1.00",
        "eval_cost": "0.20",
        "prefix_cost": "1.20",
        "data_read_per_join": "96"
      },
      "used_columns": [
        "id",
        "name",
        "surname",
        "birthday",
        "city",
        "about",
        "avatar",
        "email",
        "password_hash",
        "created_at",
        "updated_at"
      ],
      "attached_condition": "(match `social_dev`.`users`.`name`,`social_dev`.`users`.`surname` against ('+jo*' in boolean mode))"
    }
  }
}
```
We increase search query perform rapidly! FULLTEXT index query_cost: 1.20 agianst BTREE query_cost: 122099.48.

Load testing with FULLTEXT indexes:
https://clck.ru/M8FBi

