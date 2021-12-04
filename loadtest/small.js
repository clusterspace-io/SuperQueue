import http from 'k6/http';
import { sleep } from 'k6';
export const options = {
  vus: 10,
  stages: [
    { duration: '10s', target: 500 },
    { duration: '60s', target: 500 },
    { duration: '10s', target: 0 },
  ],
};
export default function () {
  const resp = http.post('http://localhost:8080/record', JSON.stringify({
    payload: 'this is a test payload'
  }), {
    headers: {
      'content-type': 'application/json'
    }
  });
  sleep(0.5);
  const resp2 = http.get('http://localhost:8080/record')
  const recordID = JSON.parse(resp2.body).id
  sleep(0.5);
  if (Math.random() > 0.5) {
    // ack
    const resp3 = http.post(`http://localhost:8080/ack/${recordID}`)
  } else {
    // nack
    const resp3 = http.post(`http://localhost:8080/nack/${recordID}`)
    sleep(0.5);
    // ack
    const resp4 = http.post(`http://localhost:8080/ack/${recordID}`)
  }
}
