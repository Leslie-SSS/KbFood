#!/bin/bash
cd /home/leslie/keepbuild/projects/food-go
export FOOD_SERVER_PORT=9000
export FOOD_DB_HOST=127.0.0.1
export FOOD_DB_PORT=5433
export FOOD_DB_USER=food
export FOOD_DB_PASSWORD=123456
export FOOD_DB_NAME=food
export FOOD_LOG_LEVEL=debug
export FOOD_PLATFORM_TANTANTANG_TOKEN=test
export FOOD_PLATFORM_TANTANTANG_SECRET_KEY=test
export FOOD_PLATFORM_TANTANTANG_BASE_URL=https://test.com
export FOOD_PLATFORM_DT_TOKEN=test
export FOOD_PLATFORM_XIAOCAN_X_VAYNE=test
export FOOD_PLATFORM_XIAOCAN_X_TEEMO=test
export FOOD_PLATFORM_XIAOCAN_X_ASHE=test
export FOOD_PLATFORM_XIAOCAN_X_NAMI=test
export FOOD_PLATFORM_XIAOCAN_X_SIVIR=test
export FOOD_PLATFORM_XIAOCAN_USER_ID=test
export FOOD_PLATFORM_XIAOCAN_SILK_ID=test

./bin/main 2>&1 | tee server-current.log
