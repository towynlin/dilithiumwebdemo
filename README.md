# Dilithium Web Demo

I wanted to do a quick spike to show a breadth of skills.

This repo is a small web service with a client and server written in Go.

It uses a post-quantum safe signature algorithm,
[CRYSTALS-Dilithium](https://pq-crystals.org/dilithium/index.shtml),
as a custom signing method for JWT authentication.

The app is Dockerized and pushed to Docker Hub.

I used Terraform to deploy the app on a Kubernetes cluster on GCP
with nginx ingress and a postegres database.

I used Let's Encrypt to generate a TLS certificate and bought a domain.

It's deployed at https://postquantumcryptography.rocks/. I will take it down shortly — don't want to be paying for the infrastructure. 😉

Initial creation of db tables can be done with `go run ./cmd/server/admin/createtables`.
You'll need to set the `DATABASE_URL` env var.

Users can be added to the database with `go run ./cmd/server/admin/adduser`.
The public key and ID get added to the database.
You must save the secret key to use for authentication.

To run the client, you'll need to set the `SERVER_URL`
without a trailing slash, as well as the signing ID and key, like this.

```
export SERVER_URL=https://postquantumcryptography.rocks
export SIGNING_ID=c98cbf05-3642-49c3-8d14-6d8df355d82e
export SIGNING_KEY=aUdneh59sh...
```

Then you can run the client with `go run ./cmd/client`.
Subcommands are `get`, `post`, and `delete <id>`.

I wanted to write a React Native front end and tests for the Go code,
but this was what I got in a 3-day spike. 🤷‍♂️😅
