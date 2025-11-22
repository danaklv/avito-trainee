import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
  vus: 10,               // число виртуальных пользователей
  duration: '10s',       // время теста
};

export default function () {
  let res = http.get('http://localhost:8080/team/get?team_name=backend');

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 300ms': (r) => r.timings.duration < 300,
  });

  sleep(0.2); 
}
