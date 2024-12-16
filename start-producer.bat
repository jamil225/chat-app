@echo off
REM Navigate to the Kafka bin directory
cd /d C:\kafka_2.13\bin\windows

REM Start Zookeeper in a new window
start "Zookeeper" cmd /k "zookeeper-server-start.bat ..\..\config\zookeeper.properties"

echo Zookeeper  have been started.
