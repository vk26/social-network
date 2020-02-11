query='SHOW MASTER STATUS;'
docker-compose exec mysql_master sh -c "export MYSQL_PWD=password; mysql -u root -e '$query'"

query='SHOW SLAVE STATUS\G'
docker-compose exec mysql_slave1 sh -c "export MYSQL_PWD=password; mysql -u root -e '$query'"
