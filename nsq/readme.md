
https://nsq.io/deployment/docker.html

docker run --name lookupd -p 4160:4160 -p 4161:4161 -d nsqio/nsq /nsqlookupd
docker run --name nsqadmin  -p 4171:4171   -d nsqio/nsq /nsqadmin  --lookupd-http-address=172.17.0.1:4161

//--nsqd-http-address=172.17.0.1:4151

docker run --name nsqd -p 4150:4150 -p 4151:4151 -d nsqio/nsq /nsqd --broadcast-address=172.17.0.1  --lookupd-tcp-address=172.17.0.1:4160