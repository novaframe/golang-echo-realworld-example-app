#!/bin/sh
# wait-for-postgres.sh

set -e
  
host="$1"
shift
cmd="$@"
  
until mysql -u$MYSQL_USER -p$MYSQL_PASSWORD -h "$host" -P 3306 -D $MYSQL_DATABASE -e '\q'; do
  >&2 echo "MySql is unavailable - sleeping"
  sleep 1
done
  
>&2 echo "MySql is up - executing command"
exec $cmd
