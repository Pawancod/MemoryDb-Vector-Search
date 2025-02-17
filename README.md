# MemoryDb-Vector-Search

This project demonstrates how to perform vector search using AWS MemoryDB for Redis. The application stores questions and answers (Q&A) pairs along with their vector representations in Redis and retrieves the most relevant answer based on a user's query using vector similarity search.


**Note:**
- **Local Environment:** (can only be done via VPN on local)The code only serves the flow and Logic behind the vector search . However, AWS MemoryDB can only be accessed internally; external access is discouraged by AWS. Therefore, this setup should be deployed on machines like EC2 within the VPC with fine changes .
- **Note:** MemoryDB with vector earch is in preview state and is only present in some regions g us-east-1, make sure we have everything setup in the availability zone only
- **Purpose:** The code serves as a demonstration of the flow and logic for vector search.


## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Components](#components)


## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/Pawancod/MemoryDb-Vector-Search.git
   cd MemoryDb-Vector-Search

## Usage

1. **Configure AWS MemoryDB:**
   - Update the MemoryDB endpoint in `internal/database/redis.go`:
     ```go
     Addr: "memoryDB endpoiint :6379", // Replace with your MemoryDB endpoint
     ```

2. **Run the vectorization service:**
   - Ensure the vectorization service is running(using this since for generating vectors from an 3rd party api are paid, can be replaced with openAI/huggin face API for vector generation):
   - Install python dependencies- present in vectorization_service/requirements.txt
   - Run the following command in the dorectory to get vectorization service up in python
      ```sh
      python3 -m venv venv
      source venv/bin/activate
      python vectorization_service/main.py


3. **Run the Go application:**
   ```sh
   go run main.go

## Steps to make this work 
  - Make the Redis instance available by setting up the cluster endpoint.
  - Add questions and answers in the `model/qa.go` file to store them in Redis using `StoreQAPair` .
  - Lines 34 to 51 contain the logic for storing the question/answer pairs(main.go).
  - Lines 53 to 66 contain the logic for performing vector search(main.go).
  - You can update the query question to test.





## Components

### 1. `main.go`
- Initializes the Redis client.
- Stores predefined Q&A pairs in Redis.
- Creates an index for vector search.
- Performs vector search based on a user's query.

### 2. `internal/database/redis.go`
- Contains functions to interact with Redis:
  - `NewRedisClient`: Initializes a Redis client.
  - `CheckRedisConnection`: Checks the connection to Redis.
  - `StoreQAPair`: Stores Q&A pairs in Redis.
  - `CreateIndex`: Creates an index for vector search in Redis.
  - `floatsToString`: Helper function to convert a slice of floats to a space-separated string.

### 3. `internal/model/qa.go`
- Defines the `QAPair` struct.
- Retrieves predefined Q&A pairs and generates their vector representations.

### 4. `internal/service/vector_service.go`
- Contains functions for vector operations:
  - `PerformVectorSearch`: Performs vector search to find the most relevant answer.
  - `GenerateVector`: Calls an external service to generate a vector for a given question.
  - `floatsToString`: Helper function to convert a slice of floats to a space-separated string.

### 5. `vectorization_service/main.py`
- A Flask service that uses the Sentence Transformers library to generate vector representations of text.



