The banking tests need the following container running:
docker run -p 8000:8000 -d lendoab/interview-service:stable

The rabbitmq tests need the following docker container running
docker run -d --hostname my-rabbit6 --name some-rabbit6 --network host rabbitmq:3.8.14-management

Rabbit MQ Management Console
http://localhost:15672/