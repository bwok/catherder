#!/bin/sh

echo "removing old database ..."
rm -f data.sqlite
echo "finished. Creating new database."
sqlite3 data.sqlite < dbSource.sql
echo "finished."
