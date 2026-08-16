// Package main is the entrypoint for the campus-lms API.
//
// -----------------------------------------------------------------------------
// TASK BRIEF — Week 1
// Agent: read AGENTS.md and agent/rules/10-go-api.md before implementing.
// Human: this file is prime interview material. After the agent implements it,
// use agent/prompts/teach.md to walk through it line by line.
// -----------------------------------------------------------------------------
//
// REQUIREMENTS
//
//  1. CONFIGURATION (12-factor)
//     - Read all configuration from environment variables (see .env.example).
//     - Fail fast: if a required variable is missing, log and exit(1).
//       Never start with partial configuration.
//     - No default value for security-sensitive variables (e.g. JWT_SECRET).
//     - Lives in internal/config.
//
//  2. STRUCTURED LOGGING
//     - log/slog with a JSON handler. Level from APP_LOG_LEVEL.
//     - Every line carries: service, env, version. trace_id is added in Week 6.
//     - No fmt.Println anywhere.
//
//  3. HTTP SERVER
//     - net/http with chi (or stdlib router).
//     - MANDATORY timeouts: ReadHeaderTimeout, ReadTimeout, WriteTimeout,
//       IdleTimeout. A server without them is vulnerable to slowloris.
//
//  4. ENDPOINTS
//     - GET /healthz — liveness: 200 while the process is alive.
//     - GET /readyz  — readiness: 200 only when dependencies are reachable
//       (DB check added in Week 3).
//     - These behave differently under Kubernetes (Week 11): a failed liveness
//       probe kills the pod; a failed readiness probe only removes it from the
//       load balancer.
//
//  5. GRACEFUL SHUTDOWN
//     - Catch SIGTERM and SIGINT via signal.NotifyContext.
//     - Stop accepting new connections, drain in-flight requests within
//       APP_SHUTDOWN_TIMEOUT, close resources, exit 0.
//     - Log "shutting down gracefully" when shutdown begins.
//     - Docker and Kubernetes send SIGTERM before SIGKILL. Without this,
//       every deploy severs live user requests.
//
// DEFINITION OF DONE (proposed by agent, ticked by human)
//   - curl -v http://localhost:8080/healthz returns 200 with a JSON log line
//   - SIGTERM produces "shutting down gracefully" and exit code 0
//   - no fmt.Println remains
//   - all four server timeouts are set
//   - human can explain liveness vs readiness out loud
//
// EVIDENCE REQUIRED (agent, see agent/evidence-protocol.md)
//   docs/progress/evidence/week-01/healthz-curl.txt
//   docs/progress/evidence/week-01/graceful-shutdown.txt
package main

func main() {
	// TODO Week 1: implement per the brief above.
	panic("not implemented yet")
}
