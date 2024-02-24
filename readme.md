# how to start

Firstly, you should create a yaml file: ``conf/deploy.local.yml``

Then you can run in local with the follwing command:

``` shell
sh build.sh
cd output
sh bootstrap.sh
```

Use the following command if you want to run in docker

```shell
sh build.sh
sudo docker build -t web_chat:latest .
sudo docker run -itd --network=host -v web_chat_log:/app/log web_chat:latest
```