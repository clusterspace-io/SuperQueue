# FORMAT: num partitions, num RRs, num 50VU test pods

All tests take the worst one, but make note on their similarity

## 1, 1, 1

data_received..................: 44 MB  1.1 MB/s
data_sent......................: 52 MB  1.3 MB/s
http_req_blocked...............: avg=2.05µs  min=674ns   med=1.63µs  max=2.17ms  p(90)=2.26µs  p(95)=2.69µs
http_req_connecting............: avg=90ns    min=0s      med=0s      max=1.73ms  p(90)=0s      p(95)=0s
http_req_duration..............: avg=5.88ms  min=3.77ms  med=4.89ms  max=32.78ms p(90)=8.26ms  p(95)=8.82ms
  { expected_response:true }...: avg=5.88ms  min=3.77ms  med=4.89ms  max=32.78ms p(90)=8.26ms  p(95)=8.82ms
http_req_failed................: 0.00%  ✓ 0           ✗ 293709
http_req_receiving.............: avg=26.34µs min=6.94µs  med=19.16µs max=11.95ms p(90)=30.87µs p(95)=41.31µs
http_req_sending...............: avg=17.78µs min=3.7µs   med=11.05µs max=11.79ms p(90)=21.07µs p(95)=31.57µs
http_req_tls_handshaking.......: avg=0s      min=0s      med=0s      max=0s      p(90)=0s      p(95)=0s
http_req_waiting...............: avg=5.84ms  min=3.73ms  med=4.85ms  max=32.75ms p(90)=8.22ms  p(95)=8.78ms
http_reqs......................: 293709 7341.560194/s
iteration_duration.............: avg=17.89ms min=15.05ms med=17.31ms max=43.7ms  p(90)=20.07ms p(95)=22.45ms
iterations.....................: 97903  2447.186731/s
vus............................: 1      min=1         max=50
vus_max........................: 50     min=50        max=50

## 1, 1, 2

**teststack_test-services.1.w2l9vkm2zvz2@swarm-04**
data_received..................: 38 MB  959 kB/s
data_sent......................: 45 MB  1.1 MB/s
http_req_blocked...............: avg=1.89µs  min=697ns   med=1.48µs  max=1.93ms  p(90)=2.13µs  p(95)=2.51µs
http_req_connecting............: avg=106ns   min=0s      med=0s      max=1.25ms  p(90)=0s      p(95)=0s
http_req_duration..............: avg=6.76ms  min=3.83ms  med=5.57ms  max=77.91ms p(90)=9.39ms  p(95)=11.81ms
  { expected_response:true }...: avg=6.76ms  min=3.83ms  med=5.57ms  max=77.91ms p(90)=9.39ms  p(95)=11.81ms
http_req_failed................: 0.00%  ✓ 0           ✗ 256419
http_req_receiving.............: avg=24.66µs min=7µs     med=18.98µs max=12.01ms p(90)=29.37µs p(95)=36.5µs
http_req_sending...............: avg=16.45µs min=3.82µs  med=10.57µs max=15.51ms p(90)=20.1µs  p(95)=27.11µs
http_req_tls_handshaking.......: avg=0s      min=0s      med=0s      max=0s      p(90)=0s      p(95)=0s
http_req_waiting...............: avg=6.72ms  min=3.8ms   med=5.52ms  max=77.87ms p(90)=9.35ms  p(95)=11.77ms
http_reqs......................: 256419 6407.723079/s
iteration_duration.............: avg=20.5ms  min=15.36ms med=18.39ms max=91.46ms p(90)=27.57ms p(95)=33.1ms
iterations.....................: 85473  2135.907693/s
vus............................: 1      min=1         max=50
vus_max........................: 50     min=50        max=50

**teststack_test-services.2.x8lbpcoi3md4@swarm-06**
data_received..................: 38 MB  958 kB/s
data_sent......................: 45 MB  1.1 MB/s
http_req_blocked...............: avg=2.03µs  min=711ns   med=1.59µs  max=1.46ms   p(90)=2.25µs  p(95)=2.66µs
http_req_connecting............: avg=106ns   min=0s      med=0s      max=766.88µs p(90)=0s      p(95)=0s
http_req_duration..............: avg=6.76ms  min=3.84ms  med=5.56ms  max=77.41ms  p(90)=9.44ms  p(95)=11.97ms
  { expected_response:true }...: avg=6.76ms  min=3.84ms  med=5.56ms  max=77.41ms  p(90)=9.44ms  p(95)=11.97ms
http_req_failed................: 0.00%  ✓ 0           ✗ 256023
http_req_receiving.............: avg=26.87µs min=7.46µs  med=20.31µs max=14.36ms  p(90)=31.47µs p(95)=39.18µs
http_req_sending...............: avg=17.71µs min=3.91µs  med=11.26µs max=13.21ms  p(90)=21.45µs p(95)=28.98µs
http_req_tls_handshaking.......: avg=0s      min=0s      med=0s      max=0s       p(90)=0s      p(95)=0s
http_req_waiting...............: avg=6.72ms  min=3.79ms  med=5.51ms  max=77.38ms  p(90)=9.39ms  p(95)=11.91ms
http_reqs......................: 256023 6398.406949/s
iteration_duration.............: avg=20.53ms min=15.32ms med=18.37ms max=91.19ms  p(90)=27.88ms p(95)=33.39ms
iterations.....................: 85341  2132.802316/s
vus............................: 1      min=1         max=50
vus_max........................: 50     min=50        max=50

## 1, 1, 3

data_received..................: 30 MB  740 kB/s
data_sent......................: 35 MB  875 kB/s
http_req_blocked...............: avg=1.89µs  min=689ns   med=1.44µs  max=1.9ms   p(90)=2.12µs  p(95)=2.48µs
http_req_connecting............: avg=148ns   min=0s      med=0s      max=1.48ms  p(90)=0s      p(95)=0s
http_req_duration..............: avg=8.78ms  min=3.84ms  med=7.68ms  max=71.18ms p(90)=13.11ms p(95)=17.85ms
  { expected_response:true }...: avg=8.78ms  min=3.84ms  med=7.68ms  max=71.18ms p(90)=13.11ms p(95)=17.85ms
http_req_failed................: 0.00%  ✓ 0           ✗ 197826
http_req_receiving.............: avg=24.56µs min=7.09µs  med=19.47µs max=9.76ms  p(90)=29.48µs p(95)=35.76µs
http_req_sending...............: avg=16.28µs min=3.8µs   med=10.46µs max=13.25ms p(90)=19.91µs p(95)=26.04µs
http_req_tls_handshaking.......: avg=0s      min=0s      med=0s      max=0s      p(90)=0s      p(95)=0s
http_req_waiting...............: avg=8.74ms  min=3.78ms  med=7.64ms  max=71.15ms p(90)=13.06ms p(95)=17.8ms
http_reqs......................: 197826 4944.053894/s
iteration_duration.............: avg=26.57ms min=15.42ms med=23.58ms max=90.94ms p(90)=40.06ms p(95)=48.84ms
iterations.....................: 65942  1648.017965/s
vus............................: 1      min=1         max=50
vus_max........................: 50     min=50        max=50

## 3, 1, 1

data_received..................: 45 MB  1.1 MB/s
data_sent......................: 53 MB  1.3 MB/s
http_req_blocked...............: avg=2.09µs  min=703ns    med=1.62µs  max=11.88ms  p(90)=2.25µs  p(95)=2.7µs
http_req_connecting............: avg=71ns    min=0s       med=0s      max=723.17µs p(90)=0s      p(95)=0s
http_req_duration..............: avg=5.8ms   min=624.41µs med=4.78ms  max=28.89ms  p(90)=8.23ms  p(95)=8.71ms
  { expected_response:true }...: avg=5.8ms   min=624.41µs med=4.78ms  max=28.89ms  p(90)=8.23ms  p(95)=8.71ms
http_req_failed................: 0.00%  ✓ 0           ✗ 298178
http_req_receiving.............: avg=25.72µs min=7.19µs   med=18.84µs max=11.63ms  p(90)=30.34µs p(95)=39.87µs
http_req_sending...............: avg=17.49µs min=3.63µs   med=11.05µs max=10.77ms  p(90)=21.07µs p(95)=30.11µs
http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s       p(90)=0s      p(95)=0s
http_req_waiting...............: avg=5.75ms  min=595.09µs med=4.73ms  max=27.65ms  p(90)=8.19ms  p(95)=8.67ms
http_reqs......................: 298178 7452.623161/s
iteration_duration.............: avg=17.59ms min=7.77ms   med=17.07ms max=41.96ms  p(90)=19.71ms p(95)=22.27ms
iterations.....................: 99574  2488.739943/s
vus............................: 1      min=1         max=50
vus_max........................: 50     min=50        max=50

## 3, 1, 3

All very similar

data_received..................: 32 MB  799 kB/s
data_sent......................: 38 MB  945 kB/s
http_req_blocked...............: avg=4.08µs  min=793ns    med=2.15µs  max=11.85ms p(90)=3.31µs  p(95)=4.4µs
http_req_connecting............: avg=132ns   min=0s       med=0s      max=1.04ms  p(90)=0s      p(95)=0s
http_req_duration..............: avg=8.05ms  min=733.59µs med=7ms     max=62.91ms p(90)=12.3ms  p(95)=16.69ms
  { expected_response:true }...: avg=8.05ms  min=733.59µs med=7ms     max=62.91ms p(90)=12.3ms  p(95)=16.69ms
http_req_failed................: 0.00%  ✓ 0           ✗ 213668
http_req_receiving.............: avg=52.03µs min=8.55µs   med=26.95µs max=17.77ms p(90)=76.99µs p(95)=143.82µs
http_req_sending...............: avg=34.38µs min=4.37µs   med=15.64µs max=20.22ms p(90)=48.21µs p(95)=100.37µs
http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s      p(90)=0s      p(95)=0s
http_req_waiting...............: avg=7.96ms  min=709µs    med=6.9ms   max=62.88ms p(90)=12.18ms p(95)=16.56ms
http_reqs......................: 213668 5340.997149/s
iteration_duration.............: avg=24.57ms min=8.08ms   med=21.41ms max=83.95ms p(90)=37.15ms p(95)=45.92ms
iterations.....................: 71290  1782.015495/s
vus............................: 1      min=1         max=50
vus_max........................: 50     min=50        max=50

## 1, 3, 1

data_received..................: 45 MB  1.1 MB/s
data_sent......................: 54 MB  1.3 MB/s
http_req_blocked...............: avg=3.7µs   min=727ns   med=2.06µs  max=8.78ms  p(90)=3.15µs  p(95)=4.05µs
http_req_connecting............: avg=86ns    min=0s      med=0s      max=1.92ms  p(90)=0s      p(95)=0s
http_req_duration..............: avg=5.65ms  min=3.65ms  med=4.69ms  max=26.01ms p(90)=7.96ms  p(95)=8.45ms
  { expected_response:true }...: avg=5.65ms  min=3.65ms  med=4.69ms  max=26.01ms p(90)=7.96ms  p(95)=8.45ms
http_req_failed................: 0.00%  ✓ 0           ✗ 302382
http_req_receiving.............: avg=48.98µs min=8.98µs  med=24.96µs max=15.99ms p(90)=61.91µs p(95)=142.73µs
http_req_sending...............: avg=31.35µs min=4.39µs  med=14.79µs max=12.93ms p(90)=37.92µs p(95)=82.97µs
http_req_tls_handshaking.......: avg=0s      min=0s      med=0s      max=0s      p(90)=0s      p(95)=0s
http_req_waiting...............: avg=5.57ms  min=3.59ms  med=4.61ms  max=25.97ms p(90)=7.88ms  p(95)=8.35ms
http_reqs......................: 302382 7557.150739/s
iteration_duration.............: avg=17.37ms min=14.69ms med=16.95ms max=37.33ms p(90)=19.16ms p(95)=20.77ms
iterations.....................: 100794 2519.050246/s
vus............................: 1      min=1         max=50
vus_max........................: 50     min=50        max=50

## 1, 3, 3

Weird hanging at the end, see the num iterations and the latency to show performance (ignore the throughput)

data_received..................: 33 MB  472 kB/s
data_sent......................: 39 MB  558 kB/s
http_req_blocked...............: avg=2.52µs  min=697ns   med=1.7µs   max=11.14ms p(90)=2.6µs   p(95)=3.14µs
http_req_connecting............: avg=102ns   min=0s      med=0s      max=1.71ms  p(90)=0s      p(95)=0s
http_req_duration..............: avg=5.88ms  min=3.55ms  med=4.81ms  max=39.85ms p(90)=8.11ms  p(95)=9.17ms
  { expected_response:true }...: avg=5.88ms  min=3.55ms  med=4.81ms  max=39.85ms p(90)=8.11ms  p(95)=9.17ms
http_req_failed................: 0.00%  ✓ 0           ✗ 220030
http_req_receiving.............: avg=30.32µs min=7.33µs  med=21.68µs max=13.75ms p(90)=36.52µs p(95)=52.98µs
http_req_sending...............: avg=19.71µs min=3.95µs  med=12.14µs max=11.36ms p(90)=23.88µs p(95)=35.58µs
http_req_tls_handshaking.......: avg=0s      min=0s      med=0s      max=0s      p(90)=0s      p(95)=0s
http_req_waiting...............: avg=5.83ms  min=3.51ms  med=4.75ms  max=39.81ms p(90)=8.06ms  p(95)=9.1ms
http_reqs......................: 220030 3152.220807/s
iteration_duration.............: avg=17.91ms min=14.54ms med=16.88ms max=49.32ms p(90)=21.87ms p(95)=24.7ms
iterations.....................: 73338  1050.663862/s
vus............................: 3      min=3         max=50
vus_max........................: 50     min=50        max=50

## 3, 3, 3

All very similar

data_received..................: 41 MB  1.0 MB/s
data_sent......................: 48 MB  1.2 MB/s
http_req_blocked...............: avg=3.88µs  min=767ns    med=2.21µs  max=12.34ms p(90)=3.4µs   p(95)=4.51µs
http_req_connecting............: avg=108ns   min=0s       med=0s      max=1.43ms  p(90)=0s      p(95)=0s
http_req_duration..............: avg=6.27ms  min=735.43µs med=5.21ms  max=34.17ms p(90)=8.76ms  p(95)=10.2ms
  { expected_response:true }...: avg=6.27ms  min=735.43µs med=5.21ms  max=34.17ms p(90)=8.76ms  p(95)=10.2ms
http_req_failed................: 0.00%  ✓ 0           ✗ 272769
http_req_receiving.............: avg=49.39µs min=9.5µs    med=27.08µs max=15.34ms p(90)=72.37µs p(95)=129.02µs
http_req_sending...............: avg=32.62µs min=4.47µs   med=15.91µs max=14.16ms p(90)=45.92µs p(95)=90.08µs
http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s      p(90)=0s      p(95)=0s
http_req_waiting...............: avg=6.19ms  min=653.57µs med=5.12ms  max=34.05ms p(90)=8.66ms  p(95)=10.07ms
http_reqs......................: 272769 6816.828046/s
iteration_duration.............: avg=19.24ms min=7.99ms   med=18.08ms max=50.21ms p(90)=23.86ms p(95)=26.48ms
iterations.....................: 91023  2274.775137/s
vus............................: 1      min=1         max=50
vus_max........................: 50     min=50        max=50

## 3, 3, 5

All pretty similar

data_received..................: 32 MB  789 kB/s
data_sent......................: 37 MB  934 kB/s
http_req_blocked...............: avg=5.08µs  min=770ns    med=2.49µs  max=10.56ms p(90)=3.86µs  p(95)=5.56µs
http_req_connecting............: avg=139ns   min=0s       med=0s      max=4.27ms  p(90)=0s      p(95)=0s
http_req_duration..............: avg=8.11ms  min=862.16µs med=7.31ms  max=56.11ms p(90)=12.83ms p(95)=16.08ms
  { expected_response:true }...: avg=8.11ms  min=862.16µs med=7.31ms  max=56.11ms p(90)=12.83ms p(95)=16.08ms
http_req_failed................: 0.00%  ✓ 0           ✗ 211105
http_req_receiving.............: avg=63.79µs min=8.69µs   med=28.66µs max=26.35ms p(90)=89.88µs p(95)=167.18µs
http_req_sending...............: avg=44.39µs min=4.46µs   med=17.24µs max=23.08ms p(90)=57.58µs p(95)=117.91µs
http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s      p(90)=0s      p(95)=0s
http_req_waiting...............: avg=8ms     min=821.3µs  med=7.22ms  max=56.07ms p(90)=12.68ms p(95)=15.89ms
http_reqs......................: 211105 5275.466524/s
iteration_duration.............: avg=24.86ms min=8.06ms   med=22.48ms max=87.23ms p(90)=35.66ms p(95)=41.29ms
iterations.....................: 70430  1760.029877/s
vus............................: 1      min=1         max=50
vus_max........................: 50     min=50        max=50

## 3, 3, 9

All pretty similar

data_received..................: 20 MB  497 kB/s
data_sent......................: 24 MB  588 kB/s
http_req_blocked...............: avg=5.19µs  min=787ns    med=2.33µs  max=17.5ms   p(90)=3.55µs  p(95)=4.96µs
http_req_connecting............: avg=206ns   min=0s       med=0s      max=2.14ms   p(90)=0s      p(95)=0s
http_req_duration..............: avg=12.97ms min=729.05µs med=9.96ms  max=128.03ms p(90)=24.1ms  p(95)=32.05ms
  { expected_response:true }...: avg=12.97ms min=729.05µs med=9.96ms  max=128.03ms p(90)=24.1ms  p(95)=32.05ms
http_req_failed................: 0.00%  ✓ 0           ✗ 133059
http_req_receiving.............: avg=67.53µs min=8.29µs   med=26.81µs max=22.79ms  p(90)=89.74µs p(95)=172.74µs
http_req_sending...............: avg=47.3µs  min=4.22µs   med=16.13µs max=22.54ms  p(90)=54.24µs p(95)=117.06µs
http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s       p(90)=0s      p(95)=0s
http_req_waiting...............: avg=12.85ms min=688.67µs med=9.86ms  max=127.99ms p(90)=23.93ms p(95)=31.88ms
http_reqs......................: 133059 3324.976543/s
iteration_duration.............: avg=39.46ms min=8.24ms   med=33.9ms  max=201.85ms p(90)=65.45ms p(95)=78.05ms
iterations.....................: 44396  1109.40003/s
vus............................: 1      min=1         max=50
vus_max........................: 50     min=50        max=50
