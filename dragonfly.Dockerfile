FROM docker.dragonflydb.io/dragonflydb/dragonfly

RUN apt update && apt install -y redis-tools