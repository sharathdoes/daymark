# Daily Brief - Setup Guide

## Project Structure

This is a monorepo with two main components:

- **Backend**: Go API server (in root directory)
- **Frontend**: Next.js 15 application (in `/frontend` directory)

## Frontend Setup

The frontend will be shown in the preview automatically. Here's what's happening:

### Configuration

- **vercel.json**: Configured to build and run the Next.js app from the `/frontend` directory
- **package.json** (frontend): Contains all necessary dependencies (Next.js, React, Tailwind, Zustand)
- **.env.local**: Pre-configured to connect to the backend API at `http://localhost:8080`

### What You Should See

When the preview loads, you'll see the Daily Brief home page with:

1. **Header**: "Daily Brief" title with navigation
2. **Hero Section**: "Test Your News Knowledge" heading
3. **Category Selection**: Available news categories (loaded from backend)
4. **Difficulty Selector**: Easy, Medium, Hard options
5. **Start Button**: Begins the quiz

### If the Preview is Blank

The most common reason is that dependencies need to be installed. The preview should handle this automatically, but if not:

1. Open the file browser in v0
2. Navigate to `/frontend`
3. The system will auto-install dependencies when you first run the dev server

### Features

#### Home Page (`/`)
- Fetches categories from backend API
- Select category and difficulty
- Start quiz button triggers quiz generation

#### Quiz Page (`/quiz`)
- Single-question interface
- Progress bar showing current question
- Answer selection with highlight on selected answer
- Next/Previous navigation buttons
- Submit quiz when all questions answered

#### Results Page (`/results/:sessionId`)
- Displays final score and percentage
- Reviews all questions with:
  - Your answer vs correct answer
  - Explanation for each question
  - Visual feedback (correct/incorrect)
- "Take Another Quiz" button to return home

### Components

- **Header**: Navigation with "Daily Brief" logo and back button
- **LoadingOverlay**: Animated loading indicator with rotating messages (5-30s)
- **ErrorBanner**: Error messages with retry and dismiss options

### Styling

- **Editorial Design**: Clean, journalistic aesthetic
- **Typography**: Crimson Text (headings) + Inter (body)
- **Colors**: Monochrome base + blue accent (#0066cc)
- **Responsive**: Mobile-first, works great on all screen sizes

### API Integration

The frontend connects to backend endpoints:

- `GET /api/categories` - List all categories
- `POST /api/quizzes/generate` - Create new quiz session
- `POST /api/quizzes/:id/submit` - Submit answers
- `GET /api/quizzes/:id/results` - Get results

### Environment Variables

- `NEXT_PUBLIC_API_URL`: Backend API URL (default: `http://localhost:8080`)

### Troubleshooting

**Blank screen**: 
- Check that dependencies are installed in `/frontend`
- Verify the API endpoint is accessible
- Check browser console for errors

**Can't connect to backend**:
- Ensure the backend is running on port 8080
- Verify `NEXT_PUBLIC_API_URL` in `.env.local`

**Categories not loading**:
- Check that the backend API is running
- Verify the `/api/categories` endpoint is accessible
- Check browser Network tab for failed requests

### Next Steps

1. The frontend will auto-load in the preview
2. Click on a category and difficulty to start a quiz
3. Answer the questions and see your results
4. Review detailed explanations for each question

Enjoy your Daily Brief quiz experience!
