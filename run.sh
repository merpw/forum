#!/bin/bash

# Function to check if a command is available
command_exists() {
    command -v "$1" &> /dev/null
}

install_nginx_linux() {
    # Linux installation
    if ! command_exists nginx; then
        echo "Nginx is not installed. Installing Nginx..."
        sudo apt-get install nginx -y
        if [ $? -ne 0 ]; then
            echo "Nginx installation failed. Please check the installation command."
            exit 1
        fi
    fi
    current_dir=$(pwd)
    echo "You are about to set the repositorys' nginx.conf as your configuration file: \
    events {}
    http {
    include mime.types;
    server {
        listen "http";
        server_name localhost;
        root $current_dir/vanilla-frontend; <--- Changed
        location /api/ {
        proxy_pass http://localhost:8080;
        }
    }
    }"

    read -n 1 -p "Do you want to continue? (y/n): " choice
    echo

    if [[ "$choice" == "y" || "$choice" == "Y" ]]; then
        current_dir=$(pwd)
        sudo bash -c "cat > $current_dir/nginx.conf <<EOL
    events {}
    http {
    include /etc/nginx/mime.types;
    server {
        listen "http";
        server_name localhost;
        root $current_dir/vanilla-frontend;

        location /api/ {
        proxy_pass http://localhost:8080;
        }
        location /api/internal/ {
            return 404;
        }
        }
    }"
    fi
    nginx -c $(pwd)/nginx.conf
    sudo services restart nginx
}

install_nginx_macos() {
    if ! command_exists brew; then
        echo "Homebrew is not installed. Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        if [ $? -ne 0 ]; then
            echo "Homebrew installation failed. Please check the installation command."
            exit 1
        fi
    fi
    if ! brew ls --versions nginx > /dev/null; then
        echo "Nginx is not installed. Installing Nginx..."
        brew install nginx
        if [ $? -ne 0 ]; then
            echo "nginx installation failed. Please check the installation command."
            exit 1
        fi
    fi
    current_dir=$(pwd)
    echo "You are about to set the repositorys' nginx.conf as your configuration file: \
    events {}
    http {
    include mime.types;
    server {
        listen "http";
        server_name localhost;
        root $current_dir/vanilla-frontend;
        location /api/ {
        proxy_pass http://localhost:8080;
        }
        location /api/internal/ {
        return 404;
            }
        }
    }"

    read -n 1 -p "Do you want to continue? (y/n): " choice
    echo

    if [[ "$choice" == "y" || "$choice" == "Y" ]]; then
        current_dir=$(pwd)
        sudo bash -c "cat > $current_dir/nginx.conf <<EOL
    events {}
    http {
    include /opt/homebrew/etc/nginx/mime.types;
    server {
        listen "http";
        server_name localhost;
        root $current_dir/vanilla-frontend;

        location /api/ {
        proxy_pass http://localhost:8080;
        }
        location /api/internal/ {
            return 404;
        }
    }
    }"
    fi
    brew services restart nginx
}

echo "Checking for Nginx installation..."
# Check the operating system
if [[ $(uname) == "Linux" ]]; then
    install_nginx_linux
   
elif [[ $(uname) == "Darwin" ]]; then
    install_nginx_macos
else
    # Unsupported operating system (Windows mostly)
    echo "Unsupported operating system"
    exit 1
fi

echo "Restarting nginx..."

echo "Compiling TypeScript files..."
(cd vanilla-frontend && npx tsc -b --force)

# Starting Go server
echo "Open page here: http://localhost:80"
(cd backend && go run main.go)
