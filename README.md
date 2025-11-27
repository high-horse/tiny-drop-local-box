# Simple File Sharing App (Tiny Drop OR Local Box)

A simple file sharing app built with Go and HTML. Users can upload files, which are made available within their local network for 1 hour before being automatically cleaned up. Other users on the same network can download the files. If the user is in a different network, they can use a special query parameter to access the files.

## Features
- **Upload Files**: Users can upload files to the server.
- **Network-Only Access**: Files are available only within the user's local network.
- **Automatic Cleanup**: Files are automatically removed 1 hour after upload.
- **Access via IP**: Users can access files by appending `?ip=<client-ip>` to the URL if they are in a different network.

## How It Works

1. **Upload a File**: When a file is uploaded, it is saved on the server and made available for download by users within the same local network.
2. **File Access**: 
    - **Same Network**: Other users on the same network can access the file directly through the app's URL.
    - **Different Network**: If a user is in a different network, they can access the file by using `?ip=<client-ip>` in the URL (e.g., `http://your-server-ip/?ip=client-ip`).
3. **File Expiry**: Files will automatically be removed from the server 1 hour after being uploaded.

## Setup

1. **Install Go**: Make sure you have Go installed on your machine. You can download it from [here](https://golang.org/dl/).

2. **Clone the Repository**:
```bash
git clone https://github.com/high-horse/tiny-drop-local-box.git
cd tiny-drop-local-box
```

3. **Run the App**:
```bash
go run main.go
```

4. **Access the App**: Open your web browser and navigate to http://localhost:8080 (or the IP address of your server) to access the file upload page.

## Configuration
- Check config.go 
