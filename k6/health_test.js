import http from 'k6/http';
import {check, sleep} from 'k6';

export default function () {
  const urlme = 'http://localhost:5000/health';
  const paramsme = {
    headers: {
      'Content-Type': 'application/json'
    },
  };

  const resme = http.get(urlme, paramsme);
  check(resme, {'status was 200': (r) => r.status == 200});
  sleep(1);
}
