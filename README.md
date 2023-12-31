
# Maani google search project
## Introduction
This guide provides detailed instructions to run the google search project locally. The project is structured with three main components: Database, Store, and Retrieval. The Store component comprises servers for the backend and admin, while the Retrieval component acts as an API gateway.this project is mainly providing uploading and downloading features for various file types in addition to an API for searching with google and retrieving images based on the provided query and saving the results . 

## project structure
This Go project follows a microservices architecture, featuring two main components:

- **Store Microservice:** Responsible for managing admin-related endpoints and a backend connected to the database.


- **Retrieval Microservice:** Handles user authentication and serves as a gateway for the Store Microservice.

The communication between microservices is facilitated through HTTP requests, promoting efficient interaction. The project relies on PostgreSQL as its chosen database.

The microservices communicate with each other using HTTP requests, and the project utilizes PostgreSQL as the database.

## Run Locally

### Getting Started

1. **Clone the Project:**

    ```bash
    git clone https://github.com/lebleuciel/maani
    ```

2. **Navigate to the Project Directory:**

    ```bash
    cd maani
    ```

### Build Images and Start Containers

To build images from the Docker files and start containers, use the following command:

```bash
make run
```

This command initializes three images:

- **Database:** A container for our PostgreSQL Database.
- **Store:** Contains two servers—one for the backend on port 9000 and another for admin on port 9001.
- **Retrieval:** Acts as an API gateway on port 8000.

### Seed Data

To generate users in the database, run the following command:

```bash
make seed-data
```
This command creates two users:
- **Admin User:**
  - **Email:** admin@example.com
  - **Password:** password
  - **Access Type:** Admin

- **Regular User:**
  - **Email:** customer@example.com
  - **Password:** password
  - **Access Type:** User




## helper functions

In the Maani project, the `helper` part within the packages directory is dedicated to providing additional functionalities and various utilities. These utilities are designed to enhance the overall capabilities of the project.

### Overview

The `helper` part encompasses several utilities that cater to different aspects of the Maani project forexample working with filesystem inorder to read and store input files.

## Usage

To interact with the program, you can choose between two options: Postman and Swagger.

### Swagger

To generate Swagger documentation, use the following command:

```bash
make gateway-api
```
Executing this command serves the Swagger documentation, providing an interactive interface to explore and interact with the APIs.

#### Authentication
Upon login, you will receive an authentication token. Use this token in the Bearer section for authorization in subsequent requests.

### Postman
For Postman users, a convenient script has been provided to automate token handling. You can download the Postman workspace from the following link and add it to your workspaces:

Postman doc: https://documenter.getpostman.com/view/13169243/2s9YeG5rNc
 

This workspace includes predefined settings for handling tokens, eliminating the need for manual entry during login.

### documentation
If you prefer exploring the project through Go documentation, use the following command to serve Go documentation:
```bash
make godoc
```
Feel free to choose the method that best suits your workflow, whether it's interacting with APIs through Swagger, utilizing Postman for streamlined testing, or exploring detailed Go documentation.

