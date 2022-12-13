Build the container:
```
docker build --tag docker-kv-server .
```

Run the built container:
```
docker run -p 8000:8000 docker-kv-server
```
