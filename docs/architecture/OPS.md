# Operations & Scaling — Persian Music Ecosystem

> **Document**: Deployment, Scaling, Monitoring, Roadmap
> **Status**: v1.0 — Final
> **Target**: 50M users, 10M tracks, 1B streams/month

---

## 1. Deployment Architecture

### 1.1 Kubernetes Cluster Design

```
┌─────────────────────────────────────────────────────────────┐
│                   Production Cluster                        │
├─────────────────────────────────────────────────────────────┤
│ ┌─────────────────────────────────────────────────────────┐ │
│ │                  Control Plane                          │ │
│ │  kube-apiserver  ·  etcd (3-node)  ·  kube-scheduler    │ │
│ │  kube-controller-manager  ·  cloud-controller-manager   │ │
│ └─────────────────────────────────────────────────────────┘ │
│                                                             │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │                   Node Pools                            │ │
│ │                                                         │ │
│ │  System Pool (3x c5.xlarge)                             │ │
│ │  ├── istio-ingressgateway, cert-manager, external-dns   │ │
│ │  ├── prometheus, grafana, loki, tempo                   │ │
│ │  └── keda, cluster-autoscaler, reloader                 │ │
│ │                                                         │ │
│ │  Service Pool (5x c5.2xlarge)                           │ │
│ │  ├── API Gateway (5 pods), User (3), Music (5)         │ │
│ │  ├── Streams (5), Search (3), Social (3)               │ │
│ │  ├── Contribution (3), Moderation (3), Gamification (2)│ │
│ │  ├── Analytics (3), Notification (2)                   │ │
│ │  └── Recommendation (3)                                │ │
│ │                                                         │ │
│ │  Data Pool (3x r5.2xlarge)                              │ │
│ │  ├── PostgreSQL (3 pods, Patroni HA)                    │ │
│ │  ├── Redis (3 pods, Redis Sentinel HA)                  │ │
│ │  ├── OpenSearch (3 pods, cluster)                       │ │
│ │  └── Kafka (3 pods, KRaft mode)                         │ │
│ │                                                         │ │
│ │  Worker Pool (3x c5.xlarge)                             │ │
│ │  ├── Stream processors (2)                              │ │
│ │  ├── Transcoding workers (2)                            │ │
│ │  └── Analytics aggregators (2)                          │ │
│ └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 Resource Specifications

| Service | CPU Request | CPU Limit | Memory Request | Memory Limit | Replicas (min/max) |
|---------|-------------|-----------|----------------|--------------|-------------------|
| API Gateway | 500m | 2000m | 512Mi | 1Gi | 3 / 10 |
| User Service | 250m | 1000m | 256Mi | 512Mi | 2 / 6 |
| Music Service | 500m | 2000m | 512Mi | 1Gi | 3 / 10 |
| Streams Service | 500m | 2000m | 512Mi | 1Gi | 3 / 10 |
| Search Service | 250m | 1000m | 512Mi | 1Gi | 2 / 6 |
| Social Service | 250m | 1000m | 512Mi | 1Gi | 2 / 6 |
| Recommendation | 500m | 2000m | 1Gi | 2Gi | 2 / 5 |
| Contribution | 250m | 1000m | 256Mi | 512Mi | 2 / 5 |
| Moderation | 250m | 1000m | 512Mi | 1Gi | 2 / 5 |
| Gamification | 250m | 500m | 256Mi | 512Mi | 2 / 4 |
| Analytics | 500m | 2000m | 1Gi | 2Gi | 2 / 6 |
| Notification | 250m | 500m | 256Mi | 512Mi | 2 / 4 |
| PostgreSQL | 2000m | 8000m | 8Gi | 32Gi | 3 / 3 (stateful) |
| Redis | 1000m | 4000m | 4Gi | 16Gi | 3 / 3 (stateful) |
| OpenSearch | 2000m | 8000m | 8Gi | 32Gi | 3 / 6 |
| Kafka | 1000m | 4000m | 4Gi | 8Gi | 3 / 6 |

### 1.3 Auto-scaling Policies

**HPA (Horizontal Pod Autoscaler):**
```yaml
# Example: Music Service
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: music-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: music-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
    - type: Pods
      pods:
        metric:
          name: http_requests_per_second
        target:
          type: AverageValue
          averageValue: 500
```

**KEDA (Kubernetes Event-driven Autoscaling):**
```yaml
# Example: Stream processor scaling by Kafka lag
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: stream-worker-scaledobject
spec:
  scaleTargetRef:
    name: stream-worker
  triggers:
    - type: kafka
      metadata:
        topic: stream.completed
        bootstrapServers: kafka-cluster:9092
        consumerGroup: stream-worker-group
        lagThreshold: "1000"
```

---

## 2. Scaling Strategy

### 2.1 Scaling Dimensions

| Dimension | Horizontal Strategy | Vertical Strategy |
|-----------|-------------------|-------------------|
| API traffic | HPA (CPU + RPS) | Larger instance types |
| Database reads | Read replicas (3-5) | Increased memory + IOPS |
| Database writes | Sharding by tenant | Write-optimized instances |
| Search | OpenSearch cluster nodes | Larger instance + faster storage |
| Cache | Redis cluster sharding | Larger memory instances |
| Streaming | CDN + regional edge caches | Higher bandwidth |
| Queue processing | KEDA (Kafka lag) | More partitions |
| Background jobs | Job queue workers | More CPU cores |

### 2.2 Database Sharding (Future)

**When to shard:** > 100M users or > 20M tracks

**Sharding key:** `user_id % shard_count` for user-owned data, `track_id % shard_count` for music data

**Shard count:** Start at 4, double when any shard exceeds 70% capacity

**Routing:** Proxy layer (pgcat/pgpool) with consistent hashing, rebalancing via virtual nodes

### 2.3 CDN Strategy

| Content | CDN | Cache TTL | Edge Locations |
|---------|-----|-----------|----------------|
| Audio files (MP3/AAC) | CloudFront + ArvanCloud | 24h (with invalidation on re-upload) | Iran + MENA + Europe |
| Album art | CloudFront + ArvanCloud | 7d | Global |
| Static assets (JS/CSS) | CloudFront | 1y (content-hashed) | Global |
| HLS segments | CloudFront | 6h | Global |
| API responses | CloudFront (for GET only) | 1min | Regional |

### 2.4 Regional Deployment

```
Phase 1: Iran (Tehran + Mashhad data centers)
  - Primary: Tehran (IRAN)
  - DR: Mashhad (IRAN) — async replication

Phase 2: MENA expansion (Dubai, Istanbul, Riyadh)
  - Regional read replicas in each PoP
  - Local CDN edge for audio delivery
  - Geo-routed DNS (users → nearest region)

Phase 3: Global (Frankfurt, Singapore, Virginia)
  - Full multi-region active-active
  - Global leader election for writes
  - Cross-region Kafka mirroring
```

---

## 3. CI/CD Pipeline

### 3.1 Pipeline Stages

```
Commit → Lint → TypeCheck → Unit Test → Build → Integration Test →
  → Security Scan → Container Build → Push → Deploy Staging →
    → E2E Test → Smoke Test → Deploy Production (Canary) →
      → Health Check → Rollout (100%)
```

### 3.2 Deployment Strategy

- **Canary deployments**: 10% → 50% → 100% over 15 minutes
- **Rollback**: Automatic if error rate > 1% or latency P99 > 500ms
- **Blue-green**: For database migrations and major version upgrades
- **Feature flags**: LaunchDarkly for gradual feature rollout
- **Built-in feature flags**: The platform also supports environment-variable-based `FEATURE_*_ENABLED` flags (defined in `.env.example`) that can disable entire modules at startup. These are exposed at runtime via `GET /api/v1/features`. Flags include `FEATURE_SOCIAL_ENABLED`, `FEATURE_AI_ENABLED`, `FEATURE_GAMIFICATION_ENABLED`, and 10 more — see `docs/architecture/API.md` for the full list.

### 3.3 GitHub Actions Workflow

```yaml
# .github/workflows/deploy.yml (abbreviated)
jobs:
  test:
    - golangci-lint
    - go test ./... -race -cover
    - npx vue-tsc --noEmit
    - npx vitest run

  build:
    - docker build -t music-api:${{ github.sha }}
    - docker build -t music-frontend:${{ github.sha }}

  deploy-staging:
    - helm upgrade --install music-api ./deploy/charts/api
    - run e2e tests against staging

  deploy-production:
    - trigger canary via ArgoCD
    - monitor for 10 min (error rate, latency, traffic)
    - promote to full rollout
```

---

## 4. Monitoring & Observability

### 4.1 Pillars

| Pillar | Tool | What It Answers |
|--------|------|-----------------|
| Metrics | Prometheus + Grafana | What's happening? |
| Logs | Loki + Grafana | Why is it happening? |
| Traces | Tempo + Grafana | Where is it happening? |
| Alerts | Alertmanager + PagerDuty | When should I care? |
| Dashboards | Grafana | How are we doing? |
| APM | OpenTelemetry | What's slow? |

### 4.2 Key Metrics

**Infrastructure:**
- CPU / Memory / Disk usage per pod (avg, P99)
- Network I/O, TCP connections
- Cluster node utilization
- PVC capacity and growth rate

**Application (RED method):**
- **Rate**: Requests per second (per endpoint, per service)
- **Errors**: Error rate (5xx, 4xx) per endpoint
- **Duration**: Latency (P50, P95, P99) per endpoint

**Business:**
- DAU / MAU (daily/monthly active users)
- Streams per second, per hour, per day
- Stream completion rate
- Search queries per second
- New registrations per day
- Contributions submitted/approved/rejected
- Moderation queue depth and age
- Cache hit ratio
- Database query performance

### 4.3 Prometheus Metrics (Custom)

```go
// Per service
http_requests_total{service, method, path, status}
http_request_duration_seconds{service, method, path}
http_request_size_bytes{service, method}
http_response_size_bytes{service, method}

// Business
music_streams_total{track_id, artist_id, source}
music_stream_duration_seconds{track_id}
music_unique_listeners{artist_id}
search_queries_total{type}
search_latency_seconds{type}
contribution_submissions_total{type, status}
moderation_queue_depth{priority}
gamification_xp_awarded_total{source}
```

### 4.4 Grafana Dashboards

| Dashboard | Audience | Refresh |
|-----------|----------|---------|
| Service Overview (RED) | All engineers | 1min |
| Database Performance | SRE / DBAs | 30s |
| Kafka Health | SRE | 30s |
| Business KPIs | Product / Management | 5min |
| User Growth & Retention | Product | 1h |
| Recommendation Quality | ML team | 1h |
| Moderation Pipeline | Trust & Safety | 5min |
| Cost Analysis | Engineering / Finance | 1h |

### 4.5 Alerting Rules

| Alert | Condition | Severity | Response |
|-------|-----------|----------|----------|
| High error rate | error_rate > 5% for 5min | Critical | PagerDuty |
| High latency | P99 > 1s for 5min | Warning | Slack |
| Service down | up == 0 for 1min | Critical | PagerDuty |
| DB connection pool exhaustion | connections > 80% for 2min | Critical | PagerDuty |
| Disk space low | disk_usage > 85% | Warning | Slack → auto PVC resize |
| Kafka consumer lag | lag > 10000 for 5min | Warning | Slack |
| Certificate expiry | < 30 days | Warning | Email |
| Daily active users drop | < 80% of 7-day avg | Warning | Slack |

### 4.6 Structured Logging (Zap)

```go
logger.Info("stream_completed",
    zap.String("user_id", userID),
    zap.String("track_id", trackID),
    zap.Int("duration_played", dur),
    zap.Float64("completion_rate", rate),
    zap.String("source", source),
    zap.Duration("latency", latency),
    zap.String("request_id", reqID),
)
```

Log levels: `debug` (dev only), `info`, `warn`, `error`, `panic`, `fatal`

No sensitive data (passwords, tokens, PII) in logs — use `zap.String("user_id", ...)` not email/phone.

### 4.7 Distributed Tracing (OpenTelemetry)

- **Instrumentation**: Auto-instrument Go services via OpenTelemetry SDK
- **Sampling**: Head-based with probability 0.1 (10%) for general traffic, 1.0 for slow requests (>500ms)
- **Context propagation**: W3C TraceContext across services and Kafka messages
- **Visualization**: Tempo traces linked to Grafana dashboards and Loki logs

---

## 5. Security & Compliance

### 5.1 Security Scanning

| Stage | Tool | What It Checks |
|-------|------|----------------|
| Commit | gitleaks | Secrets in code |
| Build | trivy | Container vulnerabilities |
| Build | snyk | Dependency vulnerabilities |
| Runtime | falco | Container behavior anomalies |
| Network | Cilium NetworkPolicy | Unauthorized connections |
| Periodic | kube-bench | CIS Kubernetes benchmark |

### 5.2 Network Security

- **Service mesh**: Istio with mTLS (STRICT mode) between all services
- **Network policies**: Default deny ingress, allow only from API Gateway and same-service
- **Ingress**: Only ports 443 (HTTPS) and 80 (redirect to 443) exposed
- **Egress**: Restricted to known endpoints (S3, CDN, email provider, payment gateway)
- **WAF**: ModSecurity with OWASP CRS + Persian-specific rules

### 5.3 Secrets Management

- **Vault** (HashiCorp) for storing secrets
- **External Secrets Operator** syncs Vault secrets to K8s secrets
- Secrets never written to logs, env vars, or config files
- Rotation: DB passwords every 90 days, API keys on compromise

---

## 6. Cost Optimization

| Area | Strategy | Estimated Savings |
|------|----------|------------------|
| Compute | Spot instances for worker pool | 60-70% |
| Compute | Right-sizing via VPA recommendations | 20-30% |
| Storage | Lifecycle policies on S3 (standard → IA → glacier) | 40-50% |
| Database | Read replicas only during peak hours (KEDA cron) | 30% |
| CDN | Cache optimization, pre-warming popular content | 20-30% |
| Monitoring | Metrics retention: 7d raw, 30d 5min, 1y 1h | 50% |

---

## 7. Disaster Recovery

### 7.1 Scenarios

| Scenario | RPO | RTO | Action |
|----------|-----|-----|--------|
| Single pod failure | 0 | < 30s | K8s auto-restart |
| Node failure | 0 | < 2min | Pod rescheduling |
| AZ outage | < 5min | < 15min | Failover to DR region |
| Entire region failure | < 15min | < 1h | Manual DNS switch to DR |
| Data corruption | < 1h (WAL) | < 2h | Point-in-time recovery |
| Catastrophic (both regions) | < 24h (backup) | < 48h | Restore from offsite backup |

### 7.2 Backup Schedule

| Component | Frequency | Retention | Storage |
|-----------|-----------|-----------|---------|
| PostgreSQL full | Daily | 30 days | S3 + cross-region |
| PostgreSQL WAL | Continuous | 7 days | S3 (same region) |
| Redis RDB | Every 5 min | 24 hours | S3 |
| OpenSearch snapshot | Daily | 14 days | S3 |
| Kafka log | Continuous | 7 days | Kafka retention |
| Application config | On change | Git history | GitHub |
| Kubernetes manifests | On change | Git history | GitHub |

---

## 8. Incident Response

### 8.1 Severity Levels

| Level | Description | Response Time | Communication |
|-------|-------------|---------------|---------------|
| SEV-1 | Complete outage, data loss | < 5min | PagerDuty + Slack + Email |
| SEV-2 | Partial outage, degraded performance | < 15min | Slack |
| SEV-3 | Minor issue, no user impact | < 1h | Slack |
| SEV-4 | Cosmetic, non-urgent | Next business day | Jira ticket |

### 8.2 Runbook Automation

- **Auto-remediation**: K8s pod restart, PVC resize, HPA adjustment
- **Manual runbooks**: Documented in `/docs/runbooks/` for complex incidents
- **Postmortem**: Required for all SEV-1 and SEV-2 incidents, blameless culture

---

## 9. Roadmap

### Phase 1: Foundation (Months 1-3)
- [x] Core backend: User, Music, Streams services
- [x] Core frontend: Player components, home, discover, search
- [x] Database schema and migrations
- [x] Design system and component library
- [ ] Authentication flow (register, login, OAuth, refresh)
- [ ] Basic playback (stream, pause, seek, queue)
- [ ] Search MVP (basic full-text, no Persian NLP yet)
- [ ] Social service MVP (follow, activity feed)
- [ ] CI/CD pipeline setup
- [ ] Kubernetes deployment (minimal)
- [ ] Monitoring stack (Prometheus + Grafana)
- [x] Architecture documentation

### Phase 2: Community & Discovery (Months 4-6)
- [ ] Persian search (transliteration, synonym, edge_ngram)
- [ ] Recommendation engine (collaborative + content-based)
- [ ] Contribution system (lyrics, metadata, translations)
- [ ] Moderation pipeline (AI + community)
- [ ] Gamification (XP, levels, badges, challenges)
- [ ] Listening parties
- [ ] Playlist collaboration
- [ ] Mobile responsive (PWA)
- [ ] Performance optimization (CDN, caching, lazy loading)

### Phase 3: Social & Creator Economy (Months 7-9)
- [ ] Live rooms (voice chat + shared queue)
- [ ] Music clubs (interest-based groups)
- [ ] Track discussions and ratings
- [ ] Creator dashboard (analytics, release scheduler)
- [ ] Creator verification and tipping
- [ ] Visualizer suite (WebGL spectrum, circular, particles, fluid)
- [ ] Fullscreen player modes (theatre, karaoke, ambient)
- [ ] Notification system (push, email, in-app)
- [ ] Admin panel (moderation, users, content, analytics)

### Phase 4: Monetization & Scale (Months 10-12)
- [ ] Subscription tiers (Free, Premium, Family, Student)
- [ ] Payment integration (Zarinpal, IDPay, Stripe)
- [ ] Ad system (audio ads, display ads, sponsored content)
- [ ] Creator payouts (pro-rata stream share)
- [ ] DRM for premium content (encrypted HLS)
- [ ] Regional deployment (MENA)
- [ ] Database sharding
- [ ] Performance tuning for 10M+ users
- [ ] Load testing and optimization

### Phase 5: Expansion (Year 2+)
- [ ] Multi-region active-active
- [ ] Podcast support
- [ ] Video clips and music videos
- [ ] AI-powered music production tools
- [ ] NFT / blockchain integration (for rare content)
- [ ] API marketplace for third-party developers
- [ ] Mobile native apps (iOS + Android)
- [ ] Desktop app (Electron/Tauri)
- [ ] Smart TV and car integration
- [ ] Lyrics translation marketplace
- [ ] Dastgah-based music theory education
- [ ] Persian poetry-to-music matching

---

## 10. Testing Strategy

| Test Type | Tool | Coverage Goal | Frequency |
|-----------|------|---------------|-----------|
| Unit (Go) | testing + testify | > 80% | Every commit |
| Unit (Vue) | vitest | > 70% | Every commit |
| Integration | testcontainers-go | Critical paths | Every PR |
| E2E | Playwright | User flows | Every deployment |
| API | Postman / Newman | All endpoints | Every deployment |
| Performance | k6 | Critical endpoints | Weekly |
| Load | Locust / k6 | Full system | Monthly |
| Chaos | Chaos Mesh | System resilience | Monthly |

### Testing Philosophy

- **TDD for business logic**: Write tests for all service handlers and repositories
- **Component testing for UI**: Test behavior, not implementation
- **Integration tests**: Use testcontainers for PostgreSQL, Redis, OpenSearch
- **Golden files**: For complex JSON responses, use golden file testing
- **Fuzz testing**: For input validation and anti-abuse logic
- **End-to-end**: Critical flows only (signup, search, play, contribute, moderate)

---

## 11. Development Workflow

### 11.1 Git Branching

```
main        ← Production-ready code
staging     ← Pre-production, passes all tests
develop     ← Integration branch
feat/*      ← Feature branches
fix/*       ← Bug fix branches
release/*   ← Release preparation branches
```

### 11.2 PR Review Checklist

- [ ] All tests pass (unit, integration, lint)
- [ ] No secrets or hardcoded credentials
- [ ] API changes documented in OpenAPI spec
- [ ] Database migrations backward-compatible
- [ ] Error handling covers all failure modes
- [ ] Logging at appropriate level (no sensitive data)
- [ ] Metrics added for new business events
- [ ] Frontend follows design system (glassmorphism, tokens)
- [ ] TypeScript strict — no `any` types
- [ ] Persian/RTL support verified (if UI change)
- [ ] Mobile responsive verified (if UI change)
