# POC Temporal

Simulating ride hailing booking workflow using Temporal.io 

## Requirements
- Go 1.14+
- JDK 8+
- Docker and docker-compose
- [Buf](http://buf.build) 

## Running
1. Generate stubs `buf generate`
2. Run Temporal server `docker-compose up`
3. Run all the servers (server/main.go and Application.java)
4. Run all the workers (worker/main.go)
5. Make new booking
    ```json
    POST http://localhost:8090/book
    {
        "userId": 1, 
        "trip": {
            "start": {
              "latitude": 324,
              "longitude": 567
            },
            "end": {
              "latitude": 324,
              "longitude": 567
            }
        }
    }
    ```
6. Finish the trip
    ```json
    POST http://localhost:8091/arrive?uid=1
    ```