# 2. Data Migration

## Database Schema Migration

1. Analyze existing MongoDB collections from Go backend:

- Events
- Users

2. Create Prisma schema in `prisma/schema.prisma`:

```prisma
generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "mongodb"
  url      = env("MONGODB_URI")
}

model User {
  id            String    @id @default(auto()) @map("_id") @db.ObjectId
  email         String    @unique
  name          String?
  createdEvents Event[]   @relation("CreatedEvents")
  invitedEvents Event[]   @relation("InvitedEvents")
}

model Event {
  id          String    @id @default(auto()) @map("_id") @db.ObjectId
  title       String
  description String?
  startTime   DateTime
  endTime     DateTime
  creator     User      @relation("CreatedEvents", fields: [creatorId], references: [id])
  creatorId   String    @db.ObjectId
  invitees    User[]    @relation("InvitedEvents")
}
```

## Data Migration Script

1. Create a migration script in `scripts/migrate-data.ts`:

```typescript
import { PrismaClient } from "@prisma/client";
import { MongoClient } from "mongodb";

const prisma = new PrismaClient();
const oldDb = new MongoClient("mongodb://localhost");

async function migrateData() {
  // Connect to old database
  await oldDb.connect();
  const database = oldDb.db("schej-it");

  // Migrate users
  const users = await database.collection("users").find({}).toArray();
  for (const user of users) {
    await prisma.user.create({
      data: {
        id: user._id.toString(),
        email: user.email,
        name: user.name,
      },
    });
  }

  // Migrate events
  const events = await database.collection("events").find({}).toArray();
  for (const event of events) {
    await prisma.event.create({
      data: {
        id: event._id.toString(),
        title: event.title,
        description: event.description,
        startTime: event.startTime,
        endTime: event.endTime,
        creatorId: event.creatorId.toString(),
        invitees: {
          connect: event.inviteeIds.map((id: any) => ({ id: id.toString() })),
        },
      },
    });
  }
}

migrateData()
  .catch(console.error)
  .finally(async () => {
    await prisma.$disconnect();
    await oldDb.close();
  });
```

2. Run the migration:

```bash
npx ts-node scripts/migrate-data.ts
```

## Verification Steps

1. Compare record counts between old and new databases
2. Verify data integrity for a sample of records
3. Test relationships between collections
4. Backup old database before switching to new system
