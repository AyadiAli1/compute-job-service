# compute-job-service

An asynchronous compute-job service: a Go gRPC server accepts calculation jobs
(DC power flow), executes them in a worker pool and streams job status updates
to clients. Built as a hands-on DevOps/cloud portfolio project.

> 🚧 Work in progress

## Roadmap

- [ ] gRPC API and Go service
- [ ] Python client with pytest integration tests
- [ ] Docker and Docker Compose
- [ ] CI/CD with GitHub Actions and GHCR
- [ ] Kubernetes (kind, Helm)
- [ ] Infrastructure as Code with Terraform on AWS
- [ ] Monitoring with Prometheus and Grafana

## License

MIT