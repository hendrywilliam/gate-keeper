import http from "k6/http";
import { sleep, check } from "k6";

// Simulate 500 users.
export const options = {
    vus: 500,
    iterations: 1000,
};

export default function () {
    const payload = JSON.stringify({
        concert_id: 1,
        ticket_category: 1,
    });
    const headers = { "Content-Type": "application/json" };
    const res = http.post("http://localhost:8080/ticket/", payload, {
        headers,
    });
    check(res, {
        "status was 201": (r) => r.status === 201,
        "status was 429": (r) => r.status === 429,
    });
    sleep(1);
}
