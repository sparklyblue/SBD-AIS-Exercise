# todo note commands

docker swarm init

docker stack deploy -c docker-compose.yml orderservice

docker secret create postgres_user docker/postgres_user_secret
docker secret create postgres_password docker/postgres_password_secret
docker secret create s3_user docker/s3_user_secret
docker secret create s3_password docker/s3_password_secret

docker service ls

docker stack ps orderservice

docker stack rm orderservice