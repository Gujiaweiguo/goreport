#!/bin/bash

set -e

echo "等待 MySQL 启动..."
until mysqladmin ping -h mysql -uroot -proot --silent; do
  echo "等待 MySQL..."
  sleep 2
done

echo "MySQL 已就绪"
