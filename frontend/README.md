# Daily Brief - Frontend

A modern Next.js frontend for the Daily Brief news quiz application.

## Features

- **Category Selection**: Choose from multiple news categories
- **Difficulty Levels**: Easy, Medium, and Hard difficulty options
- **Interactive Quiz**: Single-question interface with progress tracking
- **Result Review**: Detailed score breakdown with question explanations
- **Responsive Design**: Mobile-first, accessible interface
- **Real-time Progress**: Visual progress bar and question counter

## Tech Stack

- **Next.js 15** - React framework with App Router
- **TypeScript** - Type safety
- **Tailwind CSS** - Utility-first CSS
- **Zustand** - Lightweight state management
- **Crimsons Text & Inter** - Editorial typography

## Getting Started

### Prerequisites

- Node.js 18+ (or pnpm, yarn, bun)
- Backend API running at `http://localhost:8080/api` (default)

### Installation

1. Install dependencies:
```bash
npm install
# or
pnpm install
```

2. Create `.env.local` file:
```bash
cp .env.example .env.local
```

3. Update `NEXT_PUBLIC_API_URL` if your backend runs on a different port

### Running the Development Server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

```
frontend/
├── app/
│   ├── layout.tsx          # Root layout with fonts
│   ├── globals.css         # Global styles and design tokens
│   ├── page.tsx            # Home page (category/difficulty selection)
│   ├── quiz/
│   │   └── page.tsx        # Quiz page (single-question interface)
│   └── results/
│       └── [sessionId]/
│           └── page.tsx    # Results page (score and review)
├── components/
│   ├── Header.tsx          # Navigation header
│   ├── LoadingOverlay.tsx  # Loading state with rotating messages
│   └── ErrorBanner.tsx     # Error message display
├── lib/
│   ├── api.ts              # API client functions
│   └── quizStore.ts        # Zustand state store
├── types.ts                # TypeScript interfaces
├── tailwind.config.ts      # Tailwind configuration
├── next.config.js          # Next.js configuration
└── package.json
```

## API Integration

The frontend communicates with the backend API:

- `GET /api/categories` - Fetch available categories
- `POST /api/quizzes/generate` - Generate a new quiz
- `POST /api/quizzes/:id/submit` - Submit quiz answers
- `GET /api/quizzes/:id/results` - Fetch quiz results

## State Management

Uses Zustand for global quiz state:

- `currentSession` - Active quiz session
- `loading` - Loading state flag
- `error` - Error information
- `loadingMessage` - Current loading message

## Design System

### Colors
- **Background**: White (#ffffff)
- **Foreground**: Dark gray (#1a1a1a)
- **Accent**: Blue (#0066cc)
- **Muted**: Light gray (#f5f5f5)
- **Success**: Green (#16a34a)
- **Destructive**: Red (#dc2626)

### Typography
- **Serif**: Crimson Text (headings)
- **Sans**: Inter (body)

## Build

```bash
npm run build
npm start
```

## Deployment

To deploy on Vercel:

```bash
vercel
```

Set the environment variable `NEXT_PUBLIC_API_URL` in Vercel project settings to point to your production backend.

## Development

- Hot Module Replacement (HMR) enabled
- TypeScript strict mode
- Tailwind CSS with custom design tokens
- Accessible components (ARIA labels, keyboard navigation)

## License

MIT
