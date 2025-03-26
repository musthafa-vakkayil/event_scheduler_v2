# Event Trigger System Using Golang and Postgres

### Credits
Below are the resources utilized to create this server:

1. **YouTube Tutorials**: 
   - [ Master Background Tasks with Asynq ](https://www.youtube.com/watch?v=g1gbjMuDDP0&t=761s) 
   - [How Caching your Go API Improves Performance](https://www.youtube.com/watch?v=7B5mXWYiQZE&t=641s)
2. **ChatGPT**: For suggestions on debugging and refactoring code.
3. **Stack Overflow**: Helpful threads on issues related to events, deployment, etc., with many useful answers.
4. [**Deployment Tutorial**](https://www.youtube.com/watch?v=sSAWMr_-Co4&t=60s): This video guided the deployment of the server on AWS.

## Developer Note
This is not a fully robust solution to the problem. Given the time constraints, I did my best to provide a working solution, but I am aware that it lacks some functionalities and does not handle certain edge cases. However, this has been a wonderful learning experience and an opportunity to expand my knowledge. Thank you for this opportunity; it has given me a solid foundation to improve and build further.


### About this Project

This project is a **Golang-based Trigger Scheduler App** that allows users to schedule and manage triggers. Triggers can be of two types:

1. **Scheduled Trigger**: Can be one-time or recurring.
2. **API Trigger**: Enables users to create a trigger to hit an API endpoint, which can be invoked later using another API.

### Logs
A key feature of the system is the detailed logs for executed triggers:
- Logs remain in the **ACTIVE** state for the first 2 hours after execution.
- After 2 hours, they are moved to the **ARCHIVE** state and are hidden by default.
- Archived logs can still be viewed and are deleted after 46 hours.

Timings for log states can be easily adjusted via the eniviornment variables.

### Additional Features
- Users can test **one-time scheduled triggers** and **API triggers** without storing them in the database by sending the same payload used for creating triggers. 
- User need to register their self an acquire token to use the api
- Recently called logs api has enabled caching for better performance
- Backgroud workers are intgrated work managing events and logs


## Interacting With the Server

- **Base URL**: [http://18.205.239.53:8000/](http://18.205.239.53:8080/)
- **Swagger Documentation**: [http://18.205.239.53:8000/swagger/index.html#/](http://18.205.239.53:8080/swagger/index.html#/)
- **UI for Monitoring tasks**: [http://18.205.239.53:8080/](http://18.205.239.53:8000/)


## Setting Up Locally

### Option 1: Using Docker
1. Clone this repository- (https://github.com/musthafa-vakkayil/event_scheduler_v2)
2. Open a terminal and execute the following command:

   ```bash
   docker-compose up --build
   ```

3. The server will be running at `localhost:8000`.

### Option 2: Without Docker
1. Clone this repository.
2. Update the `app.env` file in the root level with your Postgres redis database details.
3. Open a new terminal and run:

   ```bash
   go run main.go
   ```

4. The server will be running at `localhost:8000`.


## Technical Details

ERD Diagram is given below
![Blank diagram - Page 1](event_scheduler.png)

go to [database diagram](https://dbdocs.io/musthuvakkayil/event_scheduler) to view detailed table details (use password 'database')


### 1. Database Details
We use **PostgreSQL** as the database, consisting of five tables:   
- **Users**  - Stores User Data
- **Sessions** - Stores user session
- **Events** - Stores both API and Scheduled Event Data
- **Logs** - Stores Executed Events details
- **Events_Tasks** - Stores the backgroud task id for scheduled events

The events table and session table is connected with user table via username foreign key.

The logs table stay stand alone making it easier for the backgroud task to delete and archive logs. Also persist the log even the corresponding Event is deleted 

The Event_Task table stores the Background task id for scheduled events. if a scheduled event gets deteled or gets updated(the execution time) we will mark that task as canceled. So the background worker will check this table and skip the task accordingly

### 2. Server Details
The server is implemented using Golang and is responsible for:

- **Trigger scheduling**
- **Handling API requests**
- **Log management**

**Folder Structure:** The code is organized into separate folders for better readability and maintainability. Eg:
- Handlers: Handles incoming API requests.
- Repo: Manages database interactions.
- Cache: Responsible for API Caching
- Tasks: Create and load tasks in Redis Queue
- Worker: Backgroud worker which listens to Redis Queue and Process the tasks
- Migrations: Database Migrations
- Docs: Swagger Docs and Database docs

## API Docs

### CRUD API's for User

1. **Create User**
    - **Endpoint**: `/users` (POST)
    - Description: Use this API to create user.

    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/users' \
        --header 'Content-Type: application/json' \
        --data-raw '{
            "username": "berlin",
            "full_name": "berlin",
            "password": "berlin123",
            "email":"berlin@gmail.com"
        }'
        ```

    - **Response**
        ```json
        {
            "username": "berlin",
            "full_name": "berlin",
            "email": "berlin@gmail.com",
            "password_changed_at": "0001-01-01T00:00:00Z",
            "created_at": "2025-03-25T14:31:19.048657638+05:30"
        }
        ```

2. **Login**
    - **Endpoint**: `/login` (POST)
    - Description: Use this API to get Access and Refresh Tokens.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/login' \
        --header 'Content-Type: application/json' \
        --data '{
            "username": "berlin",
            "password": "berlin123"
        }'
        ```
    - **Response**
        ```json
        {
            "session_id": "efb14e71-7f7d-4d7f-9063-2c2c8dd2080d",
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjY0YjVkMjAwLTk4NjMtNDUyOS1iZWY0LThmYmZkMjkyYTJlNyIsIlVzZXJuYW1lIjoiYmVybGluIiwiZXhwIjoxNzQyODk4MjU2LCJuYmYiOjE3NDI4OTcwNTYsImlhdCI6MTc0Mjg5NzA1Nn0.iVpEJIeTRJ-SbAFZDHMue1ccFO94rXr1gBveyGbQeMA",
            "access_token_expiry": "2025-03-25T15:54:16+05:30",
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImVmYjE0ZTcxLTdmN2QtNGQ3Zi05MDYzLTJjMmM4ZGQyMDgwZCIsIlVzZXJuYW1lIjoiYmVybGluIiwiZXhwIjoxNzQyOTgzNDU2LCJuYmYiOjE3NDI4OTcwNTYsImlhdCI6MTc0Mjg5NzA1Nn0.qzXknDTOCRjVTiDqqNzH3hl9hPEYAzMyWzZNsyv0qlU",
            "refresh_token_expiry": "2025-03-26T15:34:16+05:30",
            "user": {
                "username": "berlin",
                "full_name": "berlin",
                "email": "berlin@gmail.com",
                "password_changed_at": "0001-01-01T05:53:28+05:53",
                "created_at": "2025-03-25T14:31:19.048657Z"
            }
        }
        ```
3. **Get User**
    - **Endpoint**: `/users/{username}` (GET) - Authorized Route 
    - Description: View the user details based on username.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/users/berlin' \
        --header 'Authorization: Bearer ••••••'
        ```
    - **Response**
        ```json
        {
            "username": "berlin",
            "full_name": "berlin",
            "email": "berlin@gmail.com",
            "password_changed_at": "0001-01-01T00:00:00Z",
            "created_at": "2025-03-24T15:43:52.37642Z"
        }
        ```
4. **Delete User**
    - **Endpoint**: `/users/{username}` (DELETE) - Authorized Route 
    - Description: Delete your User Account.
    - **Curl Command**
        ```bash
        curl --location --request DELETE 'http://localhost:8000/users/berlin' \
        --header 'Authorization: Bearer ••••••'
        ```
    - **Response**
        ```json
        {
        "OK"
        }
        ```
5. **Renew Access Token**
    - **Endpoint**: `/token/renew` (POST)
    - Description: Get a new access token using refresh token.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/token/renew' \
        --header 'Content-Type: application/json' \
        --data '{
            "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImFkYjM1ZTQ5LTIwMWMtNGY2MS04ZjViLWIwZDI2OGJhYmYzYyIsIlVzZXJuYW1lIjoiYmVybGluIiwiZXhwIjoxNzQyOTczMDY1LCJuYmYiOjE3NDI4ODY2NjUsImlhdCI6MTc0Mjg4NjY2NX0.F33XBMeDvhSY3-KF_-glXzuidCZ7nJ9jjsbKxvcCmUw"
        }'
        ```
    - **Response**
        ```json
        {
            "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImY5Yjg4OGExLWVkOGEtNDAxNS1iYWQwLWJlMTA4MDdiNmU5MyIsIlVzZXJuYW1lIjoiYmVybGluIiwiZXhwIjoxNzQyODg5MzI2LCJuYmYiOjE3NDI4ODgxMjYsImlhdCI6MTc0Mjg4ODEyNn0._dO2NZI5_ElbG5EcuMkXH5KJPyuwYWRBixPEYWWg1jk",
            "access_token_expiry": "2025-03-25T13:25:26+05:30"
        }
        ```
### CRUD for API Events

#### Common Event APIS
1. **View Event**
    - **Endpoint**: `/events/{id}` (GET) - Authorized Route 
    - Description: View Event Details.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events/4' \
        --header 'Authorization: Bearer ••••••'
        ```
    - **Response**
        ```json
        {
            "ID": 4,
            "name": "First Event",
            "type": "API",
            "api_endpoint": "https://httpbin.org/get",
            "api_method": "GET",
            "api_request_body": null,
            "run_at": null,
            "after_x_mins": 0,
            "interval": 0,
            "is_recurring": false,
            "created_by": "berlin",
            "CreatedAt": "2025-03-24T20:02:55.951311Z",
            "executed_at": "2025-03-24T20:14:19.656919Z"
        }
        ```
2. **List Events**
    - **Endpoint**: `/events` (GET) - Authorized Route 
    - Description: List All Events.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events?pageNumber=1&pageSize=5' \
        --header 'Authorization: Bearer ••••••'
        ```
    - **Response**
        ```json
        [
            {
                "ID": 1,
                "name": "First Event",
                "type": "API",
                "api_endpoint": "https://httpbin.org/get",
                "api_method": "GET",
                "api_request_body": null,
                "run_at": null,
                "after_x_mins": 0,
                "interval": 0,
                "is_recurring": false,
                "created_by": "berlin",
                "CreatedAt": "2025-03-24T15:45:13.860662Z",
                "executed_at": "2025-03-24T15:46:18.847641Z"
            },
            {
                "ID": 2,
                "name": "SCHEDULED Event",
                "type": "SCHEDULED",
                "api_endpoint": "",
                "api_method": "",
                "api_request_body": null,
                "run_at": null,
                "after_x_mins": 2,
                "interval": 0,
                "is_recurring": false,
                "created_by": "berlin",
                "CreatedAt": "2025-03-24T15:45:52.2128Z",
                "executed_at": "2025-03-24T15:47:56.525319Z"
            }
        ]
        ```
3. **Delete Event**
    - **Endpoint**: `/events/{id}` (DELETE) - Authorized Route 
    - Description: Delete an Event.
    - **Curl Command**
        ```bash
        curl --location --request DELETE 'http://localhost:8000/events/10' \
        --header 'Authorization: Bearer ••••••'
        ```
    - **Response**
        ```json
        {
        "OK"
        }
        ```

#### API Events

1. **Create API Event**
    - **Endpoint**: `/events/api` (POST) - Authorized Route 
    - Description: Create an API event.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events/api' \
        --header 'Content-Type: application/json' \
        --header 'Authorization: Bearer ••••••' \
        --data '{
            "name": "First Event",
            "type": "API",
            "api_endpoint": "https://httpbin.org/get",
            "api_method": "GET"
        }'
        ```
    - **Response**
        ```json
        {
            "ID": 1,
            "name": "First Event",
            "type": "API",
            "api_endpoint": "https://httpbin.org/get",
            "api_method": "GET",
            "api_request_body": null,
            "run_at": null,
            "after_x_mins": 0,
            "interval": 0,
            "is_recurring": false,
            "created_by": "berlin",
            "CreatedAt": "2025-03-24T15:45:13.860662292Z",
            "executed_at": null
        }
        ```
2. **Update API Event**
    - **Endpoint**: `/events/api/{id}` (PUT) - Authorized Route 
    - Description: Update an API event.
    - **Curl Command**
        ```bash
        ```
    - **Response**
        ```json
        ```
3. **Execute API Event**
    - **Endpoint**: `/events/{id}/execute` (GET) - Authorized Route 
    - Description: Execute an API event to make the API call
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events/6/execute' \
        --header 'Authorization: Bearer ••••••'
        ```
    - **Response**
        ```json
        {
            "args": {},
            "headers": {
            "Accept-Encoding": "gzip",
            "Host": "httpbin.org",
            "User-Agent": "Go-http-client/2.0",
            "X-Amzn-Trace-Id": "Root=1-67e247b5-030b4b1153d8e2403dfb2d35"
            },
            "origin": "14.195.101.90",
            "url": "https://httpbin.org/get"
        }
        ```
    
#### Scheduled Events
1. **Create Scheduled Event**
    - **Endpoint**: `/events/schedule` (POST) - Authorized Route 

    - Description: Schedule at a date.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events/schedule' \
        --header 'Content-Type: application/json' \
        --header 'Authorization: Bearer ••••••' \
        --data '{
            "name": "SCHEDULED Event 4",
            "type": "SCHEDULED",
            "run_at_this_date": "2025-03-25T15:36:00+05:30"
        }'
        ```
    - Description: Schedule after fixed mins.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events/schedule' \
        --header 'Content-Type: application/json' \
        --header 'Authorization: Bearer ••••••' \
        --data '{
            "name": "SCHEDULED Event 4",
            "type": "SCHEDULED",
            "run_after_x_mins": 10,
        }'
        ```
    - Description: Schedule Recurring event.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events/schedule' \
        --header 'Content-Type: application/json' \
        --header 'Authorization: Bearer ••••••' \
        --data '{
            "name": "SCHEDULED Event 4",
            "type": "SCHEDULED",
            "run_after_x_mins": 10,
            "is_recurring": true,
            "repeat_after_x_mins":30
        }'
        ```
    - **Response**
        ```json
        {
            "ID": 11,
            "name": "SCHEDULED Event 4",
            "type": "SCHEDULED",
            "api_endpoint": "",
            "api_method": "",
            "api_request_body": null,
            "run_at": "2025-03-25T15:36:00+05:30",
            "after_x_mins": 0,
            "interval": 0,
            "is_recurring": false,
            "created_by": "berlin",
            "CreatedAt": "2025-03-25T15:34:31.578401725+05:30",
            "executed_at": null
        }
        ```
2. **Update Scheduled Event**
    - **Endpoint**: `/events/schedule/{id}` (PUT) - Authorized Route 
    - Description: Update an Scheduled event.
    - **Curl Command**
        ```bash
        curl --location --request PUT 'http://localhost:8000/events/schedule/11' \
        --header 'Content-Type: application/json' \
        --header 'Authorization: Bearer ••••••' \
        --data '{
            "name": "SCHEDULED Event 4",
            "type": "SCHEDULED",
            "run_after_x_mins": 0,
            "run_at_this_date": "2025-03-25T15:38:00+05:30"
        }'
        ```
    - **Response**
        ```json
        {
            "ID": 11,
            "name": "SCHEDULED Event 4",
            "type": "SCHEDULED",
            "api_endpoint": "",
            "api_method": "",
            "api_request_body": null,
            "run_at": "2025-03-25T15:38:00+05:30",
            "after_x_mins": 0,
            "interval": 0,
            "is_recurring": false,
            "created_by": "",
            "CreatedAt": "0001-01-01T00:00:00Z",
            "executed_at": null
        }
        ```
### Log API
1. **List Logs**
    - **Endpoint**: `/logs` (GET) - Authorized Route 
    - Description: List All Executed Events.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/logs?pageNumber=1&pageSize=5' \
        --header 'Authorization: Bearer ••••••'
        ```
    - **Response**
        ```json
        [
            {
                "id": 1,
                "executed_by": "berlin",
                "triggered_on": "2025-03-24T15:46:18.84835Z",
                "status": "200 OK",
                "is_archived": false,
                "log_type": "API_EVENT",
                "api_payload": null
            }
        ]
        ```

### Test Event API's

1. **Create Test API Event**
    - **Endpoint**: `/test/events/api` (POST) - Authorized Route 
    - Description: Create an Test API event without storing in database.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events/api' \
        --header 'Content-Type: application/json' \
        --header 'Authorization: Bearer ••••••' \
        --data '{
            "name": "First Event",
            "type": "API",
            "api_endpoint": "https://httpbin.org/get",
            "api_method": "GET"
        }'
        ```
    - **Response**
        ```json
        {
            "args": {},
            "headers": {
            "Accept-Encoding": "gzip",
            "Host": "httpbin.org",
            "User-Agent": "Go-http-client/2.0",
            "X-Amzn-Trace-Id": "Root=1-67e247b5-030b4b1153d8e2403dfb2d35"
            },
            "origin": "14.195.101.90",
            "url": "https://httpbin.org/get"
        }
        ```
2. **Create Test Scheduled Event**
    - **Endpoint**: `/test/events/schedule` (POST) - Authorized Route 

    - Description: Schedule at a date.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events/schedule' \
        --header 'Content-Type: application/json' \
        --header 'Authorization: Bearer ••••••' \
        --data '{
            "name": "SCHEDULED Event 4",
            "type": "SCHEDULED",
            "run_at_this_date": "2025-03-25T15:36:00+05:30"
        }'
        ```
    - Description: Schedule after fixed mins.
    - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events/schedule' \
        --header 'Content-Type: application/json' \
        --header 'Authorization: Bearer ••••••' \
        --data '{
            "name": "SCHEDULED Event 4",
            "type": "SCHEDULED",
            "run_after_x_mins": 10,
        }'
        ```
    - **Response**
        ```json
        {
        "OK"
        }
        ```