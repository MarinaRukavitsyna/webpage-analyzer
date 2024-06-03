# Web Application for Analyzing Webpages

A web application that performs an analysis of a webpage/URL.

### Features
- **Form with Text Field**: Users can type in the URL of the webpage to be analyzed.
- **Submit Button**: A button to send a request to the server for analysis.

### Usage
- Enter the URL of the webpage in the text field.
- Click the submit button to initiate the analysis.

## Run and Test the Application

### Run Unit Tests

- Go to the "webpage-analyzer-service/cmd/api/analyzer/" folder
- To run unit tests, use the "go test" command
```sh
cd webpage-analyzer-service/cmd/api/analyzer/
go test
```

### Srart w/a Docker

Start frontend:

-  Go to the "frontend" folder
-  Type "go run ./cmd/web"
-  You should see the following text: "Starting front end service on port 80"
```sh
cd frontend
go run ./cmd/web
```

Start backend:

- Go to the "webpage-analyzer-service/cmd/api" folder
- Run the application using "go run main.go"
- Open a web browser and go to http://localhost:8080 to test the functionality
```sh
cd webpage-analyzer-service/cmd/api
go run main.go
```

### Srart with Docker

Start frontend:

- Go to the "project" folder
- Build the frontend binary:
```sh
cd project
make build_frontend 
```
- To start the frontend, use the following command:
```sh
cd project
make start 
```
- To stop the frontend, use "make stop"
```sh
make stop 
```

Start backend:

- Go to the "project" folder
- Build all projects and start docker compose:
```sh
cd project
make up_build 
```
- To start all containers in the background without forcing build, use the following command:
```sh
cd project
make up 
```
- To stop the server, use "make down"
```sh
make down 
```

Open a web browser and go to http://localhost:8080 to test the functionality
  

## Assumptions and Decisions

### Unclear Requirements:
- Assumed basic error handling only.
- Assumed the analysis includes basic HTML properties such as version, title, number of headings, internal/external links, inaccessible links, and the presence of a login form.
- Only most popular doctypes (HTML version) are analyzed. Doctypes are taken from https://www.w3.org/QA/2002/04/valid-dtd-list.html.

### Frontend Integration:
- Chose to implement a loading panel to improve user experience during the processing of the request.
- Chose to implement a simple UI.

### Backend Integration:
- Implemented basic error handling to return error messages to the client in case of invalid input or server-side errors.
- Use goroutine to parallelise page scraping.
- Use Docker and Make automation platform for deployment. Assumed that they are installed on the test machine.

## Suggestions for Improvement

### Enhanced Error Handling:
- Provide more detailed error messages and user feedback in case of various failure scenarios.
- Test several more complex pages to detect and handle different types of errors.
- Implement retries for network-related errors when fetching the URL content.

### Security Improvements:
- Validate and sanitize user input more rigorously to prevent security vulnerabilities such as XSS or server-side injection attacks.
- Use HTTPS for secure communication between client and server.

### Scalability:
- Optimize the backend to handle concurrent requests efficiently.
- Use caching mechanisms to store analysis results for frequently analyzed URLs to reduce redundant processing.

### User Interface Enhancements:
- Improve the UI/UX with better styling and responsive design.
- Provide more detailed insights and visualizations (e.g., graphs, charts) for the analysis results.

### Logging and Monitoring:
- Implement logging for better traceability and debugging.
- Set up monitoring tools to track the performance and health of the application.

### Extend Functionality:
- Add support for analyzing more web page elements (e.g., scripts, images).
- Provide options for users to customize the analysis (e.g., choose specific elements to analyze).
- Add more doctypes for analysys.
- Research different libraries for web scraping.
- Use Swagger to implement API-first approach.
