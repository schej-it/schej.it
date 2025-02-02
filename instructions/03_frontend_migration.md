# 3. Frontend Migration

## tRPC Client Setup

1. Create tRPC client utilities (`src/utils/trpc.ts`):

```typescript
import { createTRPCNext } from "@trpc/next";
import { httpBatchLink } from "@trpc/client";
import type { AppRouter } from "@/server/routers/_app";

export const trpc = createTRPCNext<AppRouter>({
  config() {
    return {
      links: [
        httpBatchLink({
          url: "/api/trpc",
        }),
      ],
    };
  },
});
```

2. Set up tRPC provider in `src/app/providers.tsx`:

```typescript
"use client";

import { useState } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { trpc } from "@/utils/trpc";

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient());
  const [trpcClient] = useState(() => trpc.createClient());

  return (
    <trpc.Provider client={trpcClient} queryClient={queryClient}>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </trpc.Provider>
  );
}
```

3. Create tRPC route handlers (`src/server/routers/event.ts`):

```typescript
import { z } from "zod";
import { router, protectedProcedure, publicProcedure } from "../trpc";

export const eventRouter = router({
  list: publicProcedure.query(async ({ ctx }) => {
    return ctx.prisma.event.findMany({
      include: {
        creator: true,
        invitees: true,
      },
    });
  }),

  create: protectedProcedure
    .input(
      z.object({
        title: z.string(),
        description: z.string().optional(),
        startTime: z.date(),
        endTime: z.date(),
        inviteeIds: z.array(z.string()),
      })
    )
    .mutation(async ({ ctx, input }) => {
      return ctx.prisma.event.create({
        data: {
          title: input.title,
          description: input.description,
          startTime: input.startTime,
          endTime: input.endTime,
          creatorId: ctx.user.id,
          invitees: {
            connect: input.inviteeIds.map((id) => ({ id })),
          },
        },
      });
    }),
});
```

## Component Migration

1. Analyze Vue components structure:

- Map Vue components to React components
- Create equivalent React components in `src/components/`
- Migrate component logic from Options API to React hooks

## Component Migration Guidelines

1. Convert Vue Single File Components (.vue) to React TypeScript (.tsx) with tRPC:

```typescript
// Example: Converting a Vue component to React with tRPC
// Old (Vue):
<template>
  <div class="events-list">
    <div v-for="event in events" :key="event.id" class="event-card">
      <h2>{{ event.title }}</h2>
      <p>{{ event.description }}</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      events: []
    }
  },
  async mounted() {
    const response = await axios.get('/api/events');
    this.events = response.data;
  }
}
</script>

// New (React):
'use client';

import { trpc } from '@/utils/trpc';

export function EventsList() {
  const { data: events, isLoading } = trpc.event.list.useQuery();

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="events-list">
      {events?.map((event) => (
        <div key={event.id} className="event-card">
          <h2>{event.title}</h2>
          <p>{event.description}</p>
        </div>
      ))}
    </div>
  );
}
```

2. State Management Migration:

- Replace Vuex with React Context or Redux Toolkit
- Create equivalent stores in `src/store/`
- Use tRPC for server state management
- Implement React hooks for local state management

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
└── api/
    └── trpc/
        └── [trpc]/
            └── route.ts    # tRPC handler
```

## Authentication

1. Implement NextAuth.js with tRPC:

- Set up authentication providers
- Create protected routes using `protectedProcedure`
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
