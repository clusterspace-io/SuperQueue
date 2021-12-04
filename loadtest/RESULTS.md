# Test Results

## Dan's MBP, ScyllaDB in docker backed

This is not representative of real world performance, but stress testing locally. Real world performance should be far better when using real hardware for both the server and DB.

I think that the lower latency on the slower requests has to do with 2 things:

1. Less lock contention
2. My laptop's ability to handle this many requests

### full.js, with 0.5s delay between get and ack, 230 VUs

```
 dangoodman: ~/clusterSpaceCode/SuperQueue/loadtest git:(master) k6 run full.js                        6:18PM

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: full.js
     output: -

  scenarios: (100.00%) 1 scenario, 230 max VUs, 1m10s max duration (incl. graceful stop):
           * default: Up to 230 looping VUs for 40s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


running (0m40.4s), 000/230 VUs, 15884 complete and 0 interrupted iterations
default ✓ [======================================] 000/230 VUs  40s

     data_received..............: 7.8 MB 192 kB/s
     data_sent..................: 6.2 MB 154 kB/s
     http_req_blocked...........: avg=4.06µs   min=1µs      med=2µs      max=1.91ms   p(90)=4µs      p(95)=6µs
     http_req_connecting........: avg=958ns    min=0s       med=0s       max=716µs    p(90)=0s       p(95)=0s
     http_req_duration..........: avg=3.47ms   min=999µs    med=2.87ms   max=35.55ms  p(90)=6.15ms   p(95)=7.66ms
     http_req_receiving.........: avg=30.37µs  min=9µs      med=25µs     max=914µs    p(90)=49µs     p(95)=61µs
     http_req_sending...........: avg=15.18µs  min=5µs      med=13µs     max=610µs    p(90)=24µs     p(95)=32µs
     http_req_tls_handshaking...: avg=0s       min=0s       med=0s       max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...........: avg=3.43ms   min=970µs    med=2.83ms   max=35.43ms  p(90)=6.09ms   p(95)=7.6ms
     http_reqs..................: 47652  1179.257861/s
     iteration_duration.........: avg=510.92ms min=505.03ms med=510.34ms max=558.79ms p(90)=516.18ms p(95)=518.39ms
     iterations.................: 15884  393.085954/s
     vus........................: 19     min=19  max=230
     vus_max....................: 230    min=230 max=230
```

### full.js no delay, 230 VUs

```
 dangoodman: ~/clusterSpaceCode/SuperQueue/loadtest git:(master) k6 run full.js                        6:20PM

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: full.js
     output: -

  scenarios: (100.00%) 1 scenario, 230 max VUs, 1m10s max duration (incl. graceful stop):
           * default: Up to 230 looping VUs for 40s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


running (0m40.0s), 000/230 VUs, 178303 complete and 0 interrupted iterations
default ✓ [======================================] 000/230 VUs  40s

     data_received..............: 87 MB  2.2 MB/s
     data_sent..................: 70 MB  1.7 MB/s
     http_req_blocked...........: avg=2.44µs  min=0s     med=2µs     max=21.02ms  p(90)=3µs     p(95)=3µs
     http_req_connecting........: avg=112ns   min=0s     med=0s      max=970µs    p(90)=0s      p(95)=0s
     http_req_duration..........: avg=14.98ms min=1.03ms med=12.22ms max=234.55ms p(90)=26.8ms  p(95)=33.12ms
     http_req_receiving.........: avg=28.5µs  min=8µs    med=21µs    max=14.66ms  p(90)=37µs    p(95)=51µs
     http_req_sending...........: avg=12.05µs min=4µs    med=10µs    max=15.52ms  p(90)=14µs    p(95)=22µs
     http_req_tls_handshaking...: avg=0s      min=0s     med=0s      max=0s       p(90)=0s      p(95)=0s
     http_req_waiting...........: avg=14.94ms min=996µs  med=12.18ms max=234.51ms p(90)=26.76ms p(95)=33.08ms
     http_reqs..................: 534909 13361.181607/s
     iteration_duration.........: avg=45.17ms min=5.11ms med=42.24ms max=281.21ms p(90)=66.44ms p(95)=79.97ms
     iterations.................: 178303 4453.727202/s
     vus........................: 2      min=2   max=230
     vus_max....................: 230    min=230 max=230
```

### full.js no delay, 100 VUs

```
 dangoodman: ~/clusterSpaceCode/SuperQueue/loadtest git:(master) k6 run full.js                        6:21PM

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: full.js
     output: -

  scenarios: (100.00%) 1 scenario, 100 max VUs, 1m10s max duration (incl. graceful stop):
           * default: Up to 100 looping VUs for 40s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


running (0m40.0s), 000/100 VUs, 134170 complete and 0 interrupted iterations
default ✓ [======================================] 000/100 VUs  40s

     data_received..............: 66 MB  1.6 MB/s
     data_sent..................: 53 MB  1.3 MB/s
     http_req_blocked...........: avg=2.44µs  min=0s     med=2µs     max=2.34ms   p(90)=3µs     p(95)=3µs
     http_req_connecting........: avg=73ns    min=0s     med=0s      max=1.39ms   p(90)=0s      p(95)=0s
     http_req_duration..........: avg=8.62ms  min=1.12ms med=7.2ms   max=165.4ms  p(90)=14.79ms p(95)=17.86ms
     http_req_receiving.........: avg=28.99µs min=8µs    med=22µs    max=6.3ms    p(90)=41µs    p(95)=55µs
     http_req_sending...........: avg=12.44µs min=4µs    med=10µs    max=12.93ms  p(90)=16µs    p(95)=24µs
     http_req_tls_handshaking...: avg=0s      min=0s     med=0s      max=0s       p(90)=0s      p(95)=0s
     http_req_waiting...........: avg=8.58ms  min=1.08ms med=7.16ms  max=165.37ms p(90)=14.74ms p(95)=17.82ms
     http_reqs..................: 402510 10057.122286/s
     iteration_duration.........: avg=26.1ms  min=5.18ms med=24.43ms max=186.74ms p(90)=36.45ms p(95)=42.96ms
     iterations.................: 134170 3352.374095/s
     vus........................: 1      min=1   max=100
     vus_max....................: 100    min=100 max=100
```

### full.js 0.5s delay between get and ack, 100 VUs

```
 dangoodman: ~/clusterSpaceCode/SuperQueue/loadtest git:(master) ✗ k6 run full.js                      6:23PM

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: full.js
     output: -

  scenarios: (100.00%) 1 scenario, 100 max VUs, 1m10s max duration (incl. graceful stop):
           * default: Up to 100 looping VUs for 40s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


running (0m40.4s), 000/100 VUs, 6888 complete and 0 interrupted iterations
default ✓ [======================================] 000/100 VUs  40s

     data_received..............: 3.4 MB 83 kB/s
     data_sent..................: 2.7 MB 67 kB/s
     http_req_blocked...........: avg=5.2µs    min=1µs      med=3µs      max=2.92ms   p(90)=5µs      p(95)=7µs
     http_req_connecting........: avg=1.33µs   min=0s       med=0s       max=467µs    p(90)=0s       p(95)=0s
     http_req_duration..........: avg=3.56ms   min=1.02ms   med=3.16ms   max=48.11ms  p(90)=5.47ms   p(95)=6.38ms
     http_req_receiving.........: avg=33.04µs  min=10µs     med=28µs     max=228µs    p(90)=52µs     p(95)=64µs
     http_req_sending...........: avg=17.33µs  min=5µs      med=15µs     max=759µs    p(90)=27µs     p(95)=35µs
     http_req_tls_handshaking...: avg=0s       min=0s       med=0s       max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...........: avg=3.51ms   min=999µs    med=3.11ms   max=48.07ms  p(90)=5.41ms   p(95)=6.31ms
     http_reqs..................: 20664  511.489757/s
     iteration_duration.........: avg=512.06ms min=505.59ms med=511.77ms max=556.99ms p(90)=515.42ms p(95)=516.48ms
     iterations.................: 6888   170.496586/s
     vus........................: 7      min=7   max=100
     vus_max....................: 100    min=100 max=100
```

### full.js 0.1s delay between get and ack, 100 VUs

_This run did have scylla timeout errors, I find this somewhat rare when running the docker on my laptop_

```
 dangoodman: ~/clusterSpaceCode/SuperQueue/loadtest git:(master) ✗ k6 run full.js                      6:24PM

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: full.js
     output: -

  scenarios: (100.00%) 1 scenario, 100 max VUs, 1m10s max duration (incl. graceful stop):
           * default: Up to 100 looping VUs for 40s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


running (0m40.1s), 000/100 VUs, 30732 complete and 0 interrupted iterations
default ✓ [======================================] 000/100 VUs  40s

     data_received..............: 15 MB 375 kB/s
     data_sent..................: 12 MB 300 kB/s
     http_req_blocked...........: avg=3.05µs   min=1µs     med=2µs      max=2.73ms   p(90)=4µs      p(95)=5µs
     http_req_connecting........: avg=220ns    min=0s      med=0s       max=436µs    p(90)=0s       p(95)=0s
     http_req_duration..........: avg=4.62ms   min=111µs   med=2.82ms   max=605.65ms p(90)=5.66ms   p(95)=7.19ms
     http_req_receiving.........: avg=31.17µs  min=10µs    med=26µs     max=1.83ms   p(90)=50µs     p(95)=62µs
     http_req_sending...........: avg=14.82µs  min=5µs     med=13µs     max=582µs    p(90)=23µs     p(95)=32µs
     http_req_tls_handshaking...: avg=0s       min=0s      med=0s       max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...........: avg=4.57ms   min=90µs    med=2.78ms   max=605.61ms p(90)=5.6ms    p(95)=7.14ms
     http_reqs..................: 92165 2300.022967/s
     iteration_duration.........: avg=114.14ms min=105.1ms med=109.85ms max=1.01s    p(90)=115.08ms p(95)=119.13ms
     iterations.................: 30732 766.932196/s
     vus........................: 2     min=2   max=100
     vus_max....................: 100   min=100 max=100
```

### put.js no delay, 230 VUs

_This run did have scylla timeout errors, I find this somewhat rare when running the docker on my laptop_

```
 dangoodman: ~/clusterSpaceCode/SuperQueue/loadtest git:(master) ✗ k6 run put.js                       6:26PM

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: put.js
     output: -

  scenarios: (100.00%) 1 scenario, 230 max VUs, 1m10s max duration (incl. graceful stop):
           * default: Up to 230 looping VUs for 40s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


running (0m40.0s), 000/230 VUs, 356651 complete and 0 interrupted iterations
default ✓ [======================================] 000/230 VUs  40s

     data_received..............: 64 MB  1.6 MB/s
     data_sent..................: 62 MB  1.6 MB/s
     http_req_blocked...........: avg=2.61µs  min=0s     med=2µs     max=2.44ms   p(90)=3µs     p(95)=3µs
     http_req_connecting........: avg=186ns   min=0s     med=0s      max=1.5ms    p(90)=0s      p(95)=0s
     http_req_duration..........: avg=22.48ms min=2.09ms med=19.56ms max=625.32ms p(90)=34.29ms p(95)=42.57ms
     http_req_receiving.........: avg=33.35µs min=11µs   med=24µs    max=36.62ms  p(90)=45µs    p(95)=62µs
     http_req_sending...........: avg=14.39µs min=5µs    med=12µs    max=5.96ms   p(90)=18µs    p(95)=27µs
     http_req_tls_handshaking...: avg=0s      min=0s     med=0s      max=0s       p(90)=0s      p(95)=0s
     http_req_waiting...........: avg=22.44ms min=2.06ms med=19.51ms max=625.28ms p(90)=34.24ms p(95)=42.51ms
     http_reqs..................: 356651 8910.15595/s
     iteration_duration.........: avg=22.58ms min=2.15ms med=19.66ms max=625.38ms p(90)=34.4ms  p(95)=42.68ms
     iterations.................: 356651 8910.15595/s
     vus........................: 1      min=1   max=230
     vus_max....................: 230    min=230 max=230
```

### put.js 0.5s delay, 230 VUs

```
 dangoodman: ~/clusterSpaceCode/SuperQueue/loadtest git:(master) ✗ k6 run put.js                       6:27PM

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: put.js
     output: -

  scenarios: (100.00%) 1 scenario, 230 max VUs, 1m10s max duration (incl. graceful stop):
           * default: Up to 230 looping VUs for 40s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)


running (0m40.2s), 000/230 VUs, 16010 complete and 0 interrupted iterations
default ✓ [======================================] 000/230 VUs  40s

     data_received..............: 2.8 MB 71 kB/s
     data_sent..................: 2.8 MB 70 kB/s
     http_req_blocked...........: avg=7.93µs   min=1µs      med=3µs      max=1.89ms   p(90)=6µs      p(95)=9µs
     http_req_connecting........: avg=2.85µs   min=0s       med=0s       max=434µs    p(90)=0s       p(95)=0s
     http_req_duration..........: avg=5.63ms   min=2.37ms   med=5.46ms   max=30.77ms  p(90)=7.74ms   p(95)=8.51ms
     http_req_receiving.........: avg=35.16µs  min=12µs     med=30µs     max=217µs    p(90)=56µs     p(95)=70µs
     http_req_sending...........: avg=24.43µs  min=8µs      med=19µs     max=463µs    p(90)=39µs     p(95)=52µs
     http_req_tls_handshaking...: avg=0s       min=0s       med=0s       max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...........: avg=5.57ms   min=2.33ms   med=5.4ms    max=30.7ms   p(90)=7.67ms   p(95)=8.44ms
     http_reqs..................: 16010  397.778473/s
     iteration_duration.........: avg=506.72ms min=502.51ms med=506.49ms max=531.29ms p(90)=509.63ms p(95)=510.33ms
     iterations.................: 16010  397.778473/s
     vus........................: 12     min=12  max=230
     vus_max....................: 230    min=230 max=230
```
