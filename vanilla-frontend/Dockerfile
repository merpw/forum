# Use a lightweight Node.js image as the base
FROM node:alpine

# Set the working directory in the container
WORKDIR /app

# Copy the package.json and package-lock.json files to the container
COPY package*.json ./

# Install the dependencies
RUN npm install

RUN npm build

# Copy the application code to the container
COPY . .

# Expose the default port (change if necessary)
EXPOSE 3000

# Set the command to run the application
CMD ["npm", "start"]
