# Event Trigger System Using Golang and PostgreSQL

## 🚀 **Credits**
Resources used for building this server:

1. **YouTube Tutorials:**  
   - [Master Background Tasks with Asynq](https://www.youtube.com/watch?v=g1gbjMuDDP0&t=761s)  
   - [How Caching Your Go API Improves Performance](https://www.youtube.com/watch?v=7B5mXWYiQZE&t=641s)  
2. **ChatGPT:** For code debugging and refactoring suggestions.  
3. **Stack Overflow:** For helpful discussions on events, deployment, and related issues.  
4. **Deployment Tutorial:**  
   - [AWS Deployment Guide](https://www.youtube.com/watch?v=sSAWMr_-Co4&t=60s)

---

## 🛠️ **Developer Note**
This is version 2 of the [event-trigger-using-go](https://github.com/musthafa-vakkayil/event_trigger_using_go) project. The goal is to enhance performance, optimize the solution, and add new features.  

The project uses the [Asynq package](https://github.com/hibiken/asynq), a Redis-based background job queue, to manage scheduled tasks efficiently. I am continuously learning and open to optimizing it further with better solutions.  

Thank you for your support! 

---

## 📚 **About This Project**
This is a **Golang-based Trigger Scheduler App** that allows users to schedule and manage triggers.  

There are two types of triggers:  
- **Scheduled Trigger:** Supports one-time and recurring execution.  
- **API Trigger:** Allows creating triggers to call API endpoints, which can be invoked later using another API.  

---

## 🔍 **Logs**
The system maintains detailed logs of executed triggers:  
- Logs stay in the **ACTIVE** state for 2 hours after execution.  
- After 2 hours, logs move to the **ARCHIVE** state (hidden by default).  
- Archived logs remain accessible and are deleted after 46 hours.  

Log timing settings can be adjusted using environment variables.  

---

## ⚙️ **Additional Features**
- You can test **one-time scheduled triggers** and **API triggers** without storing them in the database by using the same payload.  
- Users must register and acquire a token to access the API.  
- **Caching** is enabled for the recently called logs API to boost performance.  
- **Background workers** manage events and logs efficiently.  

---

## 🌐 **Interacting With the Server**
- **Base URL:** [http://18.205.239.53:8080/](http://18.205.239.53:8080/)  
- **Swagger Documentation:** [http://18.205.239.53:8080/swagger/index.html#/](http://18.205.239.53:8080/swagger/index.html#/)  
- **UI for Monitoring Tasks:** [http://18.205.239.53:8000/](http://18.205.239.53:8000/)  

---


## 🔥 **Setting Up Locally**

### ✅ **Option 1: Using Docker**
1. Clone the repository:  
   ```bash
   git clone https://github.com/musthafa-vakkayil/event_scheduler_v2
   ```
2. Start the server:
   ```bash
   docker-compose up --build
   ```

3. The server will be running at `localhost:8000`.

### 🔧 Option 2: Without Docker
1. Clone this repository.
2. Update the app.env file with your PostgreSQL and Redis details.
3. Start the server:

   ```bash
   go run main.go
   ```

4. The server will run at localhost:8000.


## 🛠️ Technical Details
📊 ERD Diagram
![Blank diagram - Page 1](event_scheduler.png)

For detailed table information, visit:
[Database Diagram](https://dbdocs.io/musthuvakkayil/event_scheduler) (password: database)


### 1.🗄️ Database Schema 
- **Users**  - Stores user data.
- **Sessions** - Stores user sessions.
- **Events** - Stores both API and scheduled event data.
- **Logs** - Stores Executed Events details
- **Events_Tasks** - Stores background task IDs for scheduled events.

The Events and Sessions tables link to Users via username as a foreign key.
The Logs table is standalone, simplifying background task operations like log deletion and archiving, even if the corresponding event is deleted.

### 2.🌟 Server Details
The server, built with Golang, handles:

- **Trigger scheduling**
- **Handling API requests**
- **Log management**

**📁 Folder Structure:** 
- Handlers: Manages incoming API requests.
- Repo: Handles database interactions.
- Cache: Manages API caching.
- Tasks: Creates and loads tasks in the Redis queue.
- Worker: Background worker that processes tasks from the Redis queue.
- Migrations: Contains database migrations.
- Docs: Contains Swagger and database documentation.

## 🚀 Deployment
For deployment, I used an **AWS EC2 instance** with my free-tier subscription. The deployment process was automated using **GitHub Actions**.


### 💰 Pricing
The approximate cost of running a **t3.micro EC2 instance** with a **16 GB EBS volume** for 30 days (24x7):

#### 1. EC2 Instance Costs:
- **Pricing**: $0.0104/hour in the US East (N. Virginia) region.
- **Monthly Usage**: 30 days × 24 hours = 720 hours.
- **Instance Cost**: 720 × $0.0104 = **$7.49**.

#### 2. EBS Volume Costs:
- **Volume Type**: General Purpose SSD (gp3).
- **Pricing**: $0.08/GB/month.
- **For a 16 GB volume**: 16 × $0.08 = **$1.28**.

### Total Cost:
**$7.49 + $1.28 = $8.77/month**.

## 🔥 API Documentation

### 🛡️ User APIs

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
### 🚀 Event APIs

#### Common APIS
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
        - **Curl Command**
        ```bash
        curl --location 'http://localhost:8000/events/api/1' \
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
### 🛠️ Log APIs
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