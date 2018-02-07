# gomypartition
Time-series partitioning tool for MySQL/MariaDB databases

## Docker
Docker (https://www.docker.com/get-docker) is used for testing. 
This tool can be compiled and used without docker. However, all examples below are executed in docker environment.
It is also strongly recommended to test your setup before you will run it on production tables.

## Building app in docker
Application is built automatically on go container start:
```bash
docker-compose up -d
```

## Building app manually
* install go lang (https://golang.org/dl/)
* install glide package manager (https://github.com/Masterminds/glide)
* download dependencies:
```bash
cd <project_root>/src
glide install
```
* set GOPATH and compile app:
```bash
DIR="$( pwd )"
export GOPATH=$DIR

go install gomypartition
```

## Commands
```bash
gomypartition --help #or gomypartition h
```
* **docker-test**, dt  - Create partitioned table in docker (see .env settings) and insert random records there
* **info**             - Read and display partition info for specified table
   
### Docker test command
Command "gomypartition docker-test" creates partitioned table in docker container and it inserts random records there.
Thus it is possible to test this tool in isolated environment if you need it. For more details please see .env configuration file 
in project root.

Usage:
```bash
docker-compose up -d
docker-compose exec app bin/gomypartition docker-test
```

### Info command
Command "gomypartition info [command options]" gathers and displays partition info for --host=<host>, --database=<database>
and --table=<table>

#### Info command options
```bash
gomypartition h info
```
 * --host=value                 - [REQUIRED] database server hostname/ip address
 * --port=value                 - database server port (default: 3306)
 * --user=value                 - [REQUIRED] database username
 * --password=value             - [REQUIRED] database password
 * --database=value             - [REQUIRED] database schema name
 * --table=value                - [REQUIRED] database table name
 * --orderfld=value, -o value   - column for sorting partitions info records (default: "PARTITION_ORDINAL_POSITION")

Usage example:
```bash
gomypartition info --host=db --user=root --password=root \
--database=partition_test \
--table=test_table_partitioned \
--orderfld=TABLE_ROWS
```  

### Maintenance command
Command "gomypartition info [command options]" loads partition info for given table 
and it executes table partition maintenance - automatically creates new time-series partitions 
and drops the old one. To see what queries are execute, pls use --dry-run option.

#### Maintenance command options
```bash
gomypartition h maintenance
``` 
 * --host=value                 - [REQUIRED] database server hostname/ip address
 * --port=value                 - database server port (default: 3306)
 * --user=value                 - [REQUIRED] database username
 * --password=value             - [REQUIRED] database password
 * --database=value             - [REQUIRED] database schema name
 * --table=value                - [REQUIRED] database table name
 * --column=value               - [REQUIRED] column for partitioning
 * --max-partitions=value       - max. partition count (default: 50)
 * --range=value                - partition range in days, e.g. 1, 7, 30 ... (default: 30)
 * --retention=value            - partition retention in days. Partitions older than NOW() + retention will be removed (default: 90)
 * --prefix=value               - partition name prefix, e.g. prefix to will be used in partition name as to_20180801 (default: "to")
 * --dry-run                    - output queries only (no execution)

Usage example:
```bash
gomypartition maintenance --host=db --user=root --password=root \
--database=partition_test --table=test_table_partitioned \
--column=event_date \
--max-partitions=32 \
--dry-run
```  
