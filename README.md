## Prerequisites
  * Install Node.js Angular CLI (version 15)
    ```
    brew install node
    npm install -g @angular/cli@15
    ```
  * Download docker desktop
  * Setup Go
    ```
    brew install go
    ```

## Run mongoDB
```
cd backend
docker-compose up -d
```
## Run Golang backend server
```
cd backend
go run main.go
```

## Test Server
```
curl -X POST http://localhost:8080/graphql \
-H "Content-Type: application/json" \
-d '{"query": "query ($filter: String) { states(filter: $filter) { name } }", "variables": {"filter": "Calif"}}'
```

## Run Angular frontend server
```
cd frontend
ng serve
```

## Accesss
  - Frontend: 'http://localhost:4200'
  - Backend: 'http://localhost:8080/graphql'