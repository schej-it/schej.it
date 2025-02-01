# 3. Frontend Migration

## Component Migration

1. Analyze Vue components structure:

- Map Vue components to React components
- Create equivalent React components in `src/components/`
- Migrate component logic from Options API to React hooks

## Component Migration Guidelines

1. Convert Vue Single File Components (.vue) to React TypeScript (.tsx):

```typescript
// Example: Converting a Vue component to React
// Old (Vue):
<template>
  <div class="event-card">
    <h2>{{ event.title }}</h2>
    <p>{{ event.description }}</p>
  </div>
</template>

<script>
export default {
  props: ['event']
}
</script>

// New (React):
interface EventProps {
  event: {
    title: string;
    description?: string;
  }
}

export function EventCard({ event }: EventProps) {
  return (
    <div className="event-card">
      <h2>{event.title}</h2>
      <p>{event.description}</p>
    </div>
  );
}
```

2. State Management Migration:

- Replace Vuex with React Context or Redux Toolkit
- Create equivalent stores in `src/store/`
- Implement React hooks for state management

## Routing Migration

1. Convert Vue Router routes to Next.js pages:

- Move views from `frontend/src/views/` to `src/app/`
- Implement Next.js file-based routing
- Convert dynamic routes to Next.js dynamic route format

Example structure:

```
src/app/
├── page.tsx                 # Home page
├── events/
│   ├── page.tsx            # Events list
│   └── [id]/
│       └── page.tsx        # Single event page
├── profile/
│   └── page.tsx            # User profile
└── auth/
    └── [...nextauth]/
        └── route.ts        # Auth API routes
```

## API Integration

1. Create API routes in `src/app/api/`:

```typescript
// src/app/api/events/route.ts
import { NextResponse } from "next/server";
import { prisma } from "@/lib/prisma";

export async function GET() {
  const events = await prisma.event.findMany({
    include: {
      creator: true,
      invitees: true,
    },
  });
  return NextResponse.json(events);
}
```

2. Update frontend API calls:

- Replace Axios calls with Next.js API routes
- Implement proper error handling
- Add loading states

## Styling Migration

1. Convert Vue styles to Tailwind CSS:

- Migrate scoped styles to Tailwind classes
- Create custom Tailwind components for reusable styles
- Implement responsive design patterns

## Authentication

1. Implement NextAuth.js:

- Set up authentication providers
- Create protected routes
- Migrate existing auth logic to NextAuth.js

## Progressive Migration Strategy

1. Implement features in phases:

- Start with core features (events, users)
- Add authentication and authorization
- Migrate advanced features
- Test thoroughly between phases

2. Run both systems in parallel during migration:

- Deploy Next.js app to a staging environment
- Test feature parity
- Gradually transition users to new system
