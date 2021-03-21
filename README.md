The banking tests need the following container running:
docker run -p 8000:8000 lendoab/interview-service:stable

This rabbitmq tests need the following docker container running
run -d --hostname my-rabbit6 --name some-rabbit6 --network host rabbitmq:3.8.14-management
