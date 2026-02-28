# Deadlines and Errors

Performance isn't just about speed; it's about reliability under pressure.

## Deadline Exceeded Patterns

During the `sweep.txt` runs, we observed `DEADLINE_EXCEEDED` errors specifically when starting the Gateway tests at low concurrency.

### Error Example:
```
DEADLINE_EXCEEDED: context deadline exceeded
```

### Analysis:
1. **Warmup Importance**: These errors often occur during the first few requests as connections are established and the system reaches a steady state.
2. **Cascading Failure**: In a gateway model, if the backend slows down, in-flight requests can time out. Even if the backend eventually recovers, requests might already be "stale" by the time the gateway tries to return them.
3. **Configuration**: The lab uses a 20s deadline (`DEADLINE_MS=20000`), which is extremely generous. Errors here indicate either extreme CPU starvation during cold start or a misconfiguration in the timing logic of the load generator.

## Recovery
After the initial "hiccups," the system stabilizes and reports 0 errors for the remainder of the 50,000 request run, indicating that the forwarding logic is robust once "warm."
