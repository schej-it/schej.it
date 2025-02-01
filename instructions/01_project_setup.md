# 1. Project Setup

## Initial Setup

1. Create a new Next.js project with TypeScript:

```bash
npx create-next-app@latest schej-next --typescript --tailwind --eslint
cd schej-next
```

2. Install necessary dependencies:

```bash
npm install @prisma/client mongodb mongoose next-auth axios date-fns
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
│   ├── types/          # TypeScript type definitions
│   └── utils/          # Helper functions
├── prisma/             # Database schema and migrations
└── public/            # Static assets
```

## Database Setup

1. Initialize Prisma:

```bash
npx prisma init
```

2. Configure Prisma for MongoDB:

- Update `prisma/schema.prisma` to use MongoDB connector
- Migrate existing MongoDB schemas to Prisma schema format
