import http from 'k6/http';
import { sleep, check } from 'k6';

const hostPorts = [9090]
export const options = {
  stages: [
    { duration: '5s', target: 100 },
    { duration: '30s', target: 100 },
    { duration: '5s', target: 0 },
  ],
  teardownTimeout: '10s',
  insecureSkipTLSVerify: true
};

function getRandomPort() {
  const index = Math.round(Math.random() * (hostPorts.length - 1))
  return hostPorts[index]
}

export default function () {

  const resp = http.get(`https://localhost:${getRandomPort()}/hc`, JSON.stringify({
    payload: 'this is a test payload'
  }), {
    headers: {
      'content-type': 'application/json',
      'sq-queue': 'test-ns'
    }
  });
  // check(resp, {
  //   'protocol is HTTP/2': (r) => r.proto === 'HTTP/2.0',
  // })
  if (resp.status > 299 || resp.status < 200) {
    console.log('Got state code', resp.status, 'with test', resp.status_text, 'post')
  }
}
