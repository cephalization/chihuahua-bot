docker build -t chibot .
docker run --publish 8080:8080 --name chibot --rm chibot