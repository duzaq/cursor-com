docker run -it \
  --privileged \
  -e APP_USER_ID=$(id -u) \
  -e APP_GROUP_ID=$(id -g) \
  -p 5800:5800 \
  -p 5900:5900 \
  recluzegeek/cursor-container
