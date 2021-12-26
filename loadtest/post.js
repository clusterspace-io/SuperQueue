import http from 'k6/http';
import { sleep } from 'k6';
export const options = {
  stages: [
    { duration: '5s', target: 100 },
    { duration: '30s', target: 100 },
    { duration: '5s', target: 0 },
  ],
  teardownTimeout: '10s'
};
export default function () {
  const resp = http.post('http://localhost:9090/record', JSON.stringify({
    payload: 'this is a test payload'
  }), {
    headers: {
      'content-type': 'application/json',
      'sq-queue': 'test-ns'
    }
  });
  if (resp.status > 299 || resp.status < 200) {
    console.log('Got state code', resp.status, 'with test', resp.status_text, 'post')
  }
}
