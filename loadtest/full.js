import http from 'k6/http';
import { sleep } from 'k6';
export const options = {
  stages: [
    { duration: '5s', target: 230 },
    { duration: '30s', target: 230 },
    { duration: '5s', target: 0 },
  ],
  teardownTimeout: '10s'
};
export default function () {
  const resp = http.post('http://localhost:8080/record', JSON.stringify({
    payload: 'this is a test payload'
  }), {
    headers: {
      'content-type': 'application/json'
    }
  });
  if (resp.status > 299 || resp.status < 200) {
    console.log('Got state code', resp.status, 'with test', resp.status_text, 'post')
  }
  // sleep(0.5);
  const resp2 = http.get('http://localhost:8080/record')
  if (resp2.status > 299 || resp2.status < 200) {
    console.log('Got state code', resp2.status, 'with test', resp2.status_text, 'get')
  }
  try {
    if (resp2.status !== 204) {
      const recordID = JSON.parse(resp2.body).id
      // sleep(0.1);
      // ack
      const resp3 = http.post(`http://localhost:8080/ack/${recordID}`)
      if (resp3.status > 299 || resp3.status < 200) {
        console.log('Got state code', resp3.status, 'with test', resp3.status_text, 'ack')
      }
    }
  } catch (error) {
    console.error("Failed to read body:", resp2.status, resp2.body)
  }
  // if (Math.random() > 0.5) {
  // } else {
  //   // nack
  //   const resp3 = http.post(`http://localhost:8080/nack/${recordID}`)
  //   if (resp3.status > 299) {
  //     console.log('Got state code', resp3.status, 'with test', resp3.status_text, 'nack')
  //   }
  //   sleep(0.5);
  //   // ack
  //   const resp4 = http.post(`http://localhost:8080/ack/${recordID}`)
  //   if (resp4.status > 299) {
  //     console.log('Got state code', resp4.status, 'with test', resp4.status_text, 'nack-ack')
  //   }
  // }
}
