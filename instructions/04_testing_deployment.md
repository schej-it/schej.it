# 4. Testing and Deployment

## Testing Setup

1. Set up testing environment:

```bash
npm install -D jest @testing-library/react @testing-library/jest-dom
```

2. Configure Jest for Next.js:

```javascript
// jest.config.js
module.exports = {
  testEnvironment: "jsdom",
  setupFilesAfterEnv: ["<rootDir>/jest.setup.js"],
  moduleNameMapper: {
    "^@/(.*)$": "<rootDir>/src/$1",
  },
};
```

## Test Migration

1. Migrate existing tests:

- Convert Vue Test Utils tests to React Testing Library
- Update API tests for Next.js API routes
- Add new integration tests for Next.js specific features

Example test:

```typescript
// __tests__/components/EventCard.test.tsx
import { render, screen } from "@testing-library/react";
import { EventCard } from "@/components/EventCard";

describe("EventCard", () => {
  it("renders event details correctly", () => {
    const event = {
      title: "Test Event",
      description: "Test Description",
    };

    render(<EventCard event={event} />);

    expect(screen.getByText("Test Event")).toBeInTheDocument();
    expect(screen.getByText("Test Description")).toBeInTheDocument();
  });
});
```

## E2E Testing

1. Set up Playwright or Cypress:

```bash
npm install -D @playwright/test
# or
npm install -D cypress
```

2. Write E2E tests for critical user flows:

- Authentication
- Event creation and management
- User interactions
- Calendar integration

## Performance Testing

1. Implement performance monitoring:

- Set up Lighthouse CI
- Monitor Core Web Vitals
- Test loading performance
- Measure Time to Interactive (TTI)

## Deployment Setup

1. Configure deployment environment:

```bash
# Install deployment dependencies
npm install -D @vercel/node
```

2. Create deployment configuration:

```javascript
// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
  images: {
    domains: ["your-image-domain.com"],
  },
};

module.exports = nextConfig;
```

## Deployment Steps

1. Production build:

```bash
npm run build
```

2. Environment setup:

- Configure production environment variables
- Set up MongoDB production connection
- Configure NextAuth.js for production

3. Deployment options:
   a. Vercel deployment:

```bash
vercel deploy
```

b. Docker deployment:

```dockerfile
# Dockerfile
FROM node:18-alpine AS base

FROM base AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci

FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

FROM base AS runner
WORKDIR /app
ENV NODE_ENV production
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000
ENV PORT 3000
CMD ["node", "server.js"]
```

## Monitoring and Maintenance

1. Set up monitoring:

- Implement error tracking (e.g., Sentry)
- Set up performance monitoring
- Configure logging

2. Create maintenance procedures:

- Database backups
- Regular updates
- Security patches
- Performance optimization

## Rollback Plan

1. Prepare rollback strategy:

- Keep old system running in parallel
- Maintain database backups
- Document rollback procedures
- Test rollback process

2. Define rollback triggers:

- Critical bugs
- Performance issues
- Data inconsistencies
- User feedback
