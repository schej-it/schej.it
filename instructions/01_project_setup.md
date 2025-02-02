# 1. Project Setup

## Initial Setup

1. Create a new Next.js project with TypeScript:

```bash
npx create-next-app@latest schej-next --typescript --tailwind --eslint
cd schej-next
```

2. Install necessary dependencies:

```bash
# Base dependencies
npm install @prisma/client mongodb mongoose next-auth axios date-fns

# tRPC dependencies
npm install @trpc/server @trpc/client @trpc/react-query @trpc/next @tanstack/react-query zod

# Dev dependencies
npm install -D prisma @types/node
```

3. Set up environment variables:

- Create `.env.local` file with the following variables:
  ```
  MONGODB_URI=mongodb://localhost/schej-it
  NEXTAUTH_SECRET=your-secret-here
  NEXTAUTH_URL=http://localhost:3000
  ```

## Project Structure

Create the following directory structure:

```
schej-next/
├── src/
│   ├── app/             # Next.js 13+ app directory
│   ├── components/      # Reusable components
│   ├── lib/            # Utility functions and shared logic
│   ├── models/         # Database models
│   ├── server/         # tRPC server code
│   │   ├── routers/    # tRPC route handlers
│   │   ├── context.ts  # tRPC context
│   │   └── trpc.ts     # tRPC initialization
│   ├── types/          # TypeScript type definitions
│   └── utils/          # Helper functions
├── prisma/             # Database schema and migrations
└── public/            # Static assets
```

## tRPC Setup

1. Create the tRPC initialization file (`src/server/trpc.ts`):

```typescript
import { initTRPC, TRPCError } from "@trpc/server";
import { Context } from "./context";

const t = initTRPC.context<Context>().create();

export const router = t.router;
export const publicProcedure = t.procedure;
export const middleware = t.middleware;

// Protected procedure middleware
const isAuthed = middleware(({ ctx, next }) => {
  if (!ctx.session?.user) {
    throw new TRPCError({ code: "UNAUTHORIZED" });
  }
  return next({
    ctx: {
      ...ctx,
      user: ctx.session.user,
    },
  });
});

export const protectedProcedure = t.procedure.use(isAuthed);
```

2. Create the tRPC context (`src/server/context.ts`):

```typescript
import { inferAsyncReturnType } from "@trpc/server";
import { getServerSession } from "next-auth";
import { prisma } from "@/lib/prisma";

export async function createContext({ req, res }: any) {
  const session = await getServerSession(req, res);

  return {
    req,
    res,
    prisma,
    session,
  };
}

export type Context = inferAsyncReturnType<typeof createContext>;
```

3. Set up the root router (`src/server/routers/_app.ts`):

```typescript
import { router } from "../trpc";
import { eventRouter } from "./event";
import { userRouter } from "./user";

export const appRouter = router({
  event: eventRouter,
  user: userRouter,
});

export type AppRouter = typeof appRouter;
```

## Database Setup

1. Initialize Prisma:

```bash
npx prisma init
```

2. Configure Prisma for MongoDB:

- Update `prisma/schema.prisma` to use MongoDB connector
- Migrate existing MongoDB schemas to Prisma schema format
