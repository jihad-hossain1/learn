- enter exiting container or run container

```bash
docker exec -it <container-id> bash
```

- copy local file to an container

```bash
docker cp <local-folder-file-path> <container-id>:<which-path-or-folder-you-want-give-here>
```

- docker image commit, that's mean you open a container and install some tools,essential software now you want to save it an image this command follow

```bash
docker commit <running-container-id> <image-name>
```

- build an `Dockerfile`

```bash
 docker build -t <image-name> .
```
