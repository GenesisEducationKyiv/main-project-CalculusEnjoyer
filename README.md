# Ultimate Crypto Currency Rate Provider
### By Volodymyr Kravchuk

This project is a powerful and scalable solution for fetching real-time cryptocurrency 
rates and dispatching to subscribed emails. Implemented using a microservice architecture of 4 services 
and the gRPC for connection between them. All services are independent and run in separate Docker containers.

<br />
<img src="https://github.com/GenesisEducationKyiv/main-project-CalculusEnjoyer/blob/main/docs/arch.png">
<br />

## Running the application
Open `services` directory in the terminal and run:

```docker compose up```

Application will run on localhost:8080 by default.
## Setting up sender email
Test credentials for quick testing are already set, so you do not have to do anything 
in order to set up it. But if you want to change sender email, you can do this in `email` service
in `.env` file (path: `services/email/.env`)
## Requests

```
GET  -> http://localhost:8080/api/rate
POST -> http://localhost:8080/api/subscribe             
POST -> http://localhost:8080/api/sendEmails
```
[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-718a45dd9cf7e7f842a935f5ebbe5719a5e09af4491e668f4dbf3b35d5cca122.svg)](https://classroom.github.com/online_ide?assignment_repo_id=11353472&assignment_repo_type=AssignmentRepo)
