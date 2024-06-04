#! /bin/bash

zipf=$1

ADDRESS=/home/Concurrency/test_data/test/redis_test

go run ${ADDRESS}/cleardata.go
go run ${ADDRESS}/inputdata.go ${zipf}

echo "Reset redis"