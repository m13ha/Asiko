#!/bin/bash

cd backend

# Start the backend application

docker compose up

#!/bin/bash


# Change to the frontend directory
cd ../frontend/

# Install dependencies if needed
npm install

# Set environment variables for development
export REACT_APP_API_URL=http://127.0.0.1:8888

# Start the frontend application
npm start