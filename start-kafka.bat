@echo off
REM Navigate to the Kafka bin directory
cd /d C:\kafka_2.13\bin\windows

REM Start Kafka in a new window
start "Kafka" cmd /k "kafka-server-start.bat ..\..\config\server.properties"

echo  Kafka have been started.
