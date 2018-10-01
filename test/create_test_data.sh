#!/bin/bash

PROJ_NAME=project_$RANDOM
TEST_NAME=test_$RANDOM

http POST localhost:3000/api/projects name=$PROJ_NAME
echo "Created project $PROJ_NAME"

http POST localhost:3000/api/projects/$PROJ_NAME/tests name=$TEST_NAME
echo "Created test $TEST_NAME"

./fix_test_data.js

FILES=./run*.json
for f in $FILES
do
	echo "Processing $f $PROJ_NAME $TEST_NAME"
    http POST localhost:3000/api/projects/$PROJ_NAME/tests/$TEST_NAME/raw/ @$f
done
