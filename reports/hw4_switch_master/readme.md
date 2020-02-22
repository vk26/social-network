# Switch master to slave replica without loss of transactions
In this experiment we try to add one more slave, setup GTID, semi-sync replication, run db client for insertion data, stop force our master and promote one of slave to master. Check our success transactions in replication nodes.

## Add another slave-replica
In our docker-compose file add another replica (mysql_slave2):
```yaml
version: '3'
services:
  mysql_master:
    build: ./mysql_master
    container_name: mysql_master
    env_file: ./mysql_master/master.env
    volumes:
      - ./mysql_master/cnf/my.cnf:/etc/mysql/my.cnf
      - mysql_master_data:/var/lib/mysql
    ports:
      - 4406:3306
    networks: 
      - app_network  
  mysql_slave1:
    build: ./mysql_slave1
    container_name: mysql_slave1
    ports: 
      - 5506:3306
    env_file: ./mysql_slave1/slave.env
    volumes: 
      - ./mysql_slave1/cnf/my.cnf:/etc/mysql/my.cnf
      - mysql_slave1:/var/lib/mysql  
    networks: 
      - app_network
    depends_on:
      - mysql_master  
  mysql_slave2:
    build: ./mysql_slave2
    container_name: mysql_slave2
    ports: 
      - 6606:3306
    env_file: ./mysql_slave1/slave.env
    volumes: 
      - ./mysql_slave2/cnf/my.cnf:/etc/mysql/my.cnf
      - mysql_slave2:/var/lib/mysql  
    networks: 
      - app_network
    depends_on:
      - mysql_master  
volumes:
  mysql_master_data:
  mysql_slave1:
  mysql_slave2:
networks:
  app_network:
```

Find out binlog file and position of replication in mysql master:
```
mysql> show master status;
+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000001 |      155 |              |                  |                   |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)
```

Go to our slave2 container and enter mysql console:
```bash
docker exec -it mysql_slave2 bash
mysql -uroot -p
```
Setup replication:
```sql
CHANGE MASTER TO MASTER_HOST='172.21.0.2', MASTER_USER='mysql_slave_user', MASTER_PASSWORD='password', MASTER_LOG_FILE='mysql-bin.000001', MASTER_LOG_POS=155;
```

Start slave replication:
```sql
START SLAVE;
```

Check our slave status:
```
mysql> SHOW SLAVE STATUS\G
*************************** 1. row ***************************
               Slave_IO_State: Waiting for master to send event
                  Master_Host: 172.21.0.2
                  Master_User: mysql_slave_user
                  Master_Port: 3306
                Connect_Retry: 60
              Master_Log_File: mysql-bin.000001
          Read_Master_Log_Pos: 155
               Relay_Log_File: mysql-relay-bin.000002
                Relay_Log_Pos: 322
        Relay_Master_Log_File: mysql-bin.000001
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
              Replicate_Do_DB: 
          Replicate_Ignore_DB: 
           Replicate_Do_Table: 
       Replicate_Ignore_Table: 
      Replicate_Wild_Do_Table: 
  Replicate_Wild_Ignore_Table: 
                   Last_Errno: 0
                   Last_Error: 
                 Skip_Counter: 0
          Exec_Master_Log_Pos: 155
              Relay_Log_Space: 530
              Until_Condition: None
               Until_Log_File: 
                Until_Log_Pos: 0
           Master_SSL_Allowed: No
           Master_SSL_CA_File: 
           Master_SSL_CA_Path: 
              Master_SSL_Cert: 
            Master_SSL_Cipher: 
               Master_SSL_Key: 
        Seconds_Behind_Master: 0
Master_SSL_Verify_Server_Cert: No
                Last_IO_Errno: 0
                Last_IO_Error: 
               Last_SQL_Errno: 0
               Last_SQL_Error: 
  Replicate_Ignore_Server_Ids: 
             Master_Server_Id: 1
                  Master_UUID: 13900c38-4b34-11ea-bc02-0242c0a84002
             Master_Info_File: mysql.slave_master_info
                    SQL_Delay: 0
          SQL_Remaining_Delay: NULL
      Slave_SQL_Running_State: Slave has read all relay log; waiting for more updates
           Master_Retry_Count: 86400
                  Master_Bind: 
      Last_IO_Error_Timestamp: 
     Last_SQL_Error_Timestamp: 
               Master_SSL_Crl: 
           Master_SSL_Crlpath: 
           Retrieved_Gtid_Set: 
            Executed_Gtid_Set: 
                Auto_Position: 0
         Replicate_Rewrite_DB: 
                 Channel_Name: 
           Master_TLS_Version: 
       Master_public_key_path: 
        Get_master_public_key: 0
            Network_Namespace: 
1 row in set (0.00 sec)
```
Pay attention to Master_Log_File: mysql-bin.000001, Read_Master_Log_Pos: 155. This values match values of master. That's OK!

## Using row-based replication
In master side into my.cnf file add:
```
binlog_format = ROW
```

## Turn on GTID replication
Setup gtid-mode in my.cnf in both sides(master and slave):
```
gtid_mode=ON
enforce-gtid-consistency=ON
```
Check master status:
```
mysql> show master status;
+------------------+----------+--------------+------------------+--------------------------------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set                          |
+------------------+----------+--------------+------------------+--------------------------------------------+
| mysql-bin.000001 |   106205 |              |                  | 13900c38-4b34-11ea-bc02-0242c0a84002:1-200 |
+------------------+----------+--------------+------------------+--------------------------------------------+
1 row in set (0.00 sec)
```

Stop slave replication:
```sql
STOP SLAVE;
```

Configure the slave to use GTID-based auto-positioning:
```sql
CHANGE MASTER TO MASTER_HOST='172.21.0.2', MASTER_USER='mysql_slave_user', MASTER_PASSWORD='password', MASTER_AUTO_POSITION=1;
```

Check slave status:
```
START SLAVE;
SHOW SLAVE STATUS\G
*************************** 1. row ***************************
               Slave_IO_State: Waiting for master to send event
                  Master_Host: 172.21.0.2
                  Master_User: mysql_slave_user
                  Master_Port: 3306
                Connect_Retry: 60
              Master_Log_File: mysql-bin.000001
          Read_Master_Log_Pos: 106205
               Relay_Log_File: mysql-relay-bin.000002
                Relay_Log_Pos: 106419
        Relay_Master_Log_File: mysql-bin.000001
             Slave_IO_Running: Yes
            Slave_SQL_Running: Yes
              Replicate_Do_DB: 
          Replicate_Ignore_DB: 
           Replicate_Do_Table: 
       Replicate_Ignore_Table: 
      Replicate_Wild_Do_Table: 
  Replicate_Wild_Ignore_Table: 
                   Last_Errno: 0
                   Last_Error: 
                 Skip_Counter: 0
          Exec_Master_Log_Pos: 106205
              Relay_Log_Space: 106627
              Until_Condition: None
               Until_Log_File: 
                Until_Log_Pos: 0
           Master_SSL_Allowed: No
           Master_SSL_CA_File: 
           Master_SSL_CA_Path: 
              Master_SSL_Cert: 
            Master_SSL_Cipher: 
               Master_SSL_Key: 
        Seconds_Behind_Master: 0
Master_SSL_Verify_Server_Cert: No
                Last_IO_Errno: 0
                Last_IO_Error: 
               Last_SQL_Errno: 0
               Last_SQL_Error: 
  Replicate_Ignore_Server_Ids: 
             Master_Server_Id: 1
                  Master_UUID: 13900c38-4b34-11ea-bc02-0242c0a84002
             Master_Info_File: mysql.slave_master_info
                    SQL_Delay: 0
          SQL_Remaining_Delay: NULL
      Slave_SQL_Running_State: Slave has read all relay log; waiting for more updates
           Master_Retry_Count: 86400
                  Master_Bind: 
      Last_IO_Error_Timestamp: 
     Last_SQL_Error_Timestamp: 
               Master_SSL_Crl: 
           Master_SSL_Crlpath: 
           Retrieved_Gtid_Set: 13900c38-4b34-11ea-bc02-0242c0a84002:1-200
            Executed_Gtid_Set: 13900c38-4b34-11ea-bc02-0242c0a84002:1-200
                Auto_Position: 1
         Replicate_Rewrite_DB: 
                 Channel_Name: 
           Master_TLS_Version: 
       Master_public_key_path: 
        Get_master_public_key: 0
            Network_Namespace: 
1 row in set (0.00 sec)
```

## Setup semisync replication mode

We need setup plugin for semi-sync replication and turn on semi-sync replication in both sides(master and all slaves). Add to master my.cnf file:
```
[mysqld]
plugin-load=rpl_semi_sync_master=semisync_master.so
rpl_semi_sync_master_enabled=1
rpl_semi_sync_master_timeout=10000
```
And for slaves:
```
[mysqld]
plugin-load=rpl_semi_sync_slave=semisync_slave.so
rpl_semi_sync_slave_enabled=1
```
After restart mysql check installed variables and status of semi-sync replication. For example in master:
```
mysql> SHOW VARIABLES LIKE 'rpl_semi_sync%';
+-------------------------------------------+------------+
| Variable_name                             | Value      |
+-------------------------------------------+------------+
| rpl_semi_sync_master_enabled              | ON         |
| rpl_semi_sync_master_timeout              | 10000      |
| rpl_semi_sync_master_trace_level          | 32         |
| rpl_semi_sync_master_wait_for_slave_count | 1          |
| rpl_semi_sync_master_wait_no_slave        | ON         |
| rpl_semi_sync_master_wait_point           | AFTER_SYNC |
+-------------------------------------------+------------+
6 rows in set (0.00 sec)

mysql> SHOW STATUS LIKE 'Rpl_semi_sync%';
+--------------------------------------------+-------+
| Variable_name                              | Value |
+--------------------------------------------+-------+
| Rpl_semi_sync_master_clients               | 2     |
| Rpl_semi_sync_master_net_avg_wait_time     | 0     |
| Rpl_semi_sync_master_net_wait_time         | 0     |
| Rpl_semi_sync_master_net_waits             | 6     |
| Rpl_semi_sync_master_no_times              | 1     |
| Rpl_semi_sync_master_no_tx                 | 27    |
| Rpl_semi_sync_master_status                | ON    |
| Rpl_semi_sync_master_timefunc_failures     | 0     |
| Rpl_semi_sync_master_tx_avg_wait_time      | 10755 |
| Rpl_semi_sync_master_tx_wait_time          | 32266 |
| Rpl_semi_sync_master_tx_waits              | 3     |
| Rpl_semi_sync_master_wait_pos_backtraverse | 0     |
| Rpl_semi_sync_master_wait_sessions         | 0     |
| Rpl_semi_sync_master_yes_tx                | 3     |
+--------------------------------------------+-------+
```

And show slave semi-sync status:
```
mysql> SHOW STATUS LIKE 'Rpl_semi_sync%';
+----------------------------+-------+
| Variable_name              | Value |
+----------------------------+-------+
| Rpl_semi_sync_slave_status | ON    |
+----------------------------+-------+
```



