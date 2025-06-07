# Learn Event-Driven Microservices in Go

Here’s a suggested end-to-end learning path—spanning roughly 10 focused modules—to take you from “zero” to running production-grade event-driven microservices in Go, containerized on Kubernetes, documented with Markdown+Mermaid, visualized via Event Catalog, and structured via DDD + Clean Architecture.

---

## Module 1: Foundations of Go & Microservices (1 week)

**Objectives:**

* Install Go, set up workspace, understand modules & packages
* Write idiomatic Go: types, functions, error handling, interfaces
* Basic HTTP microservice with `net/http`

**Topics & Exercises:**

* Hello world & “CLI vs HTTP”
* Structs, methods, interfaces → define a `User` service interface
* Error patterns (`errors.Wrap`, custom types)
* Build a simple REST endpoint:

  ```go
  type User struct { ID, Name string }
  func handleGetUsers(w http.ResponseWriter, r *http.Request) { … }
  ```

* Hands-on: implement `GET /users` backed by in-memory slice

---

## Module 2: Advanced Go & Clean Architecture (1 week)

**Objectives:**

* Understand Clean Architecture layers: Entities, Use Cases, Interfaces, Infrastructure
* Organize Go folders: `domain/`, `usecase/`, `interface/`, `infra/`

**Topics & Exercises:**

* **Domain**: core `Order` and methods
* **Use Case**: `CreateOrder`, `GetOrder` in `usecase/order.go`
* **Interface**: HTTP handlers, messaging adapters
* **Infra**: Postgres repo, Kafka client
* **Example diagram:**

  ```mermaid
  graph LR
    UI -->|calls| UseCase
    UseCase -->|persists| Repo
    Repo -->|talks to| DB
  ```

* Hands-on: refactor Module 1 service into Clean Architecture

---

## Module 3: Domain-Driven Design Essentials (1 week)

**Objectives:**

* Bounded contexts, aggregates, entities vs value objects
* Ubiquitous language, strategic vs tactical DDD

**Topics & Exercises:**

* Identify bounded contexts: e.g. Orders vs Payments
* Model an `Order` aggregate: root entity, invariants
* Define value objects (`Money`, `Address`)
* Hands-on: DDD sketch of a “Booking” context in Markdown + Mermaid

---

## Module 4: Event-Driven Architecture Concepts (1 week)

**Objectives:**

* Patterns: pub/sub, event sourcing, CQRS
* Tradeoffs: consistency, durability, idempotency

**Topics & Exercises:**

* Event definition: name, schema, metadata
* Example event in Markdown:

  ```markdown
  ### OrderCreated
  - eventType: `order.created`
  - payload:
    - orderId: UUID
    - total: Money
  - metadata:
    - occurredAt: timestamp
  ```

* Mermaid sequence:

  ```mermaid
  sequenceDiagram
    Client->>OrderService: PlaceOrder
    OrderService->>EventBus: Emit OrderCreated
    PaymentService->>EventBus: Subscribe order.created
  ```

* Hands-on: design 3 key events for an e-commerce flow

---

## Module 5: Building Event-Driven Go Microservices (2 weeks)

**Objectives:**

* Use a message broker (e.g. Kafka, NATS JetStream) in Go
* Publish & subscribe, schema validation

**Topics & Exercises:**

* Choose broker client (e.g. sarama for Kafka)
* Define event structs, marshal/unmarshal to JSON or Protobuf
* Implement publisher:

  ```go
  func PublishOrderCreated(evt OrderCreated) error { … }
  ```

* Implement subscriber with at-least-once semantics & idempotency
* Hands-on: build two services—Order and Inventory—that communicate via events; write unit & integration tests

---

## Module 6: Containerization with Docker & Kubernetes Basics (1 week)

**Objectives:**

* Dockerfile for Go microservice, multi-stage builds
* Kubernetes primitives: Pod, Deployment, Service

**Topics & Exercises:**

* Write `Dockerfile`:

  ```dockerfile
  FROM golang:1.20 AS build
  WORKDIR /app
  COPY . .
  RUN go build -o order-service
  FROM gcr.io/distroless/base
  COPY --from=build /app/order-service /
  ENTRYPOINT ["/order-service"]
  ```

* Kubernetes YAML:

  ```yaml
  apiVersion: apps/v1
  kind: Deployment
  metadata: { name: order-service }
  spec:
    replicas: 2
    template:
      spec:
        containers:
        - name: order
          image: your-registry/order-service:latest
  ```

* Hands-on: deploy both services on a local K8s (e.g. kind/minikube) and test connectivity

---

## Module 7: Kubernetes Advanced Patterns (1 week)

**Objectives:**

* ConfigMaps, Secrets, health/liveness/readiness probes
* Helm Charts for templating

**Topics & Exercises:**

* Add probes to YAML
* Create ConfigMap for broker endpoints
* Write a basic Helm Chart to deploy Order + Inventory services together
* Hands-on: package Chart, deploy with overrides

---

## Module 8: Event Catalog for Visualization (1 week)

**Objectives:**

* Install and configure [Event Catalog](https://www.eventcatalog.dev/)
* Document and visualize your event schemas

**Topics & Exercises:**

* Define events in YAML for Event Catalog:

  ```yaml
  events:
    - name: order.created
      title: Order Created
      version: "1.0"
      schema:
        $ref: "./schemas/order-created.json"
  ```

* Run `ecctl` or Docker image, point at your YAML
* Browse UI to see lineage, schemas, owners
* Hands-on: document all events from Module 5, explore relationships

---

## Module 9: Documentation with Markdown & Mermaid (1 week)

**Objectives:**

* Write clear API & architecture docs in Markdown
* Use Mermaid for diagrams: class, sequence, component

**Topics & Exercises:**

* Example API doc snippet:

  ````markdown
  # Order Service API
  ## POST /orders
  **Request**  
  ```json
  { "items": [...] }
  ````

  **Response** 201

  ```json
  { "orderId": "123" }
  ```

* Mermaid component diagram:

  ```mermaid
  graph TD
    A[Client] --> B[API Gateway]
    B --> C[Order Service]
    B --> D[Inventory Service]
  ```

* Hands-on: produce a GitHub–style README covering setup, API, events, infra

---

## Module 10: Capstone Project & Production-Grade Practices (2 weeks)

**Objectives:**

* Bring it all together: design, code, deploy, observe
* CI/CD pipelines, monitoring, logging, tracing

**Topics & Exercises:**

* Set up GitHub Actions or GitLab CI: build, test, Docker push, Helm deploy
* Integrate Prometheus + Grafana for metrics, Jaeger for traces
* Define SLIs/SLOs, alerts
* **Capstone**:

  * Build a “Booking” system with 3 services (booking, payment, notification)
  * Event-driven flow, deployed on Kubernetes, documented in Markdown+Mermaid
  * Event Catalog populated with all event definitions
  * CI/CD, monitoring
* Code review & retrospective

---

### Additional Resources

* **Go**: “The Go Programming Language” (Alan A. A. Donovan)
* **Clean Architecture**: Robert C. Martin’s guide
* **DDD**: “Implementing Domain-Driven Design” (Vaughn Vernon)
* **K8s**: Official docs + “Kubernetes Patterns” (Bilgin Ibryam)
* **Event Catalog docs**: [https://www.eventcatalog.dev/docs/](https://www.eventcatalog.dev/docs/)
