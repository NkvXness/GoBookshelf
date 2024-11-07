# GoBookshelf

GoBookshelf is a full-stack web application for managing a personal book collection. Built with Go and React, it provides a simple and efficient way to manage your books with proper ISBN validation and data management.

## Features

### Backend

- RESTful API built with Go
- SQLite database for data persistence
- ISBN-13 validation and formatting
- Request validation and error handling
- Pagination support
- CORS support
- Unit and integration tests

### Frontend

- Modern React application with Vite
- Clean and responsive UI with Tailwind CSS
- Real-time data updates with React Query
- Form validation with proper error handling
- ISBN formatting and validation
- Loading states and error feedback

## Tech Stack

### Backend

- Go 1.21+
- SQLite3
- go-playground/validator

### Frontend

- React 18
- Vite
- Tailwind CSS
- React Query
- Axios

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- npm or yarn

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/GoBookshelf.git
cd GoBookshelf
```

2. Install backend dependencies:

```bash
go mod tidy
```

3. Install frontend dependencies:

```bash
cd frontend
npm install
```

## Running the Application

### Backend

1. Navigate to the project root directory
2. Run the server:

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Frontend

1. Navigate to the frontend directory:

```bash
cd frontend
```

2. Start the development server:

```bash
npm run dev
```

The application will be available at `http://localhost:5173`

## Using the Application

### Managing Books

1. **View Books**

   - Open the application in your browser
   - Books are displayed in a paginated table
   - Use the pagination controls at the bottom to navigate

2. **Add a Book**

   - Click the "Add New Book" button
   - Fill in the required fields:
     - Title (required, max 200 characters)
     - Author (required, max 100 characters)
     - ISBN-13 (required, must be valid)
     - Published Date (required)
   - Click "Add Book" to save

3. **Edit a Book**

   - Click the pencil icon next to a book
   - Modify the desired fields
   - Click "Save" to update or "Cancel" to discard changes

4. **Delete a Book**
   - Click the trash icon next to a book
   - Confirm the deletion when prompted

### ISBN Format

The application expects ISBN-13 format:

- 13 digits (e.g., 978-3-16-148410-0)
- Automatically validates checksum
- Automatically formats with hyphens

## API Endpoints

- `GET /books?page=1&page_size=10` - List books with pagination
- `GET /books/{id}` - Get a specific book
- `POST /books` - Create a new book
- `PUT /books/{id}` - Update an existing book
- `DELETE /books/{id}` - Delete a book

## Development

### Running Tests

Backend tests:

```bash
go test ./...
```

### Project Structure

```
GoBookshelf/
├── cmd/
│   └── server/           # Application entrypoint
├── configs/              # Configuration files
├── frontend/            # React application
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── utils/      # Utility functions
│   │   └── ...
├── internal/
│   ├── api/            # API handlers
│   ├── errors/         # Error handling
│   ├── models/         # Data models
│   └── storage/        # Database operations
└── migrations/         # Database migrations
```

## License

This project is licensed under the MIT License
