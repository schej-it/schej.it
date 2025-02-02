# Database Restructuring: Event Responses Optimization

## Current Structure

Currently, event responses are stored as a map in the Event model:

```go
type Event struct {
    // ... other fields
    Responses map[string]*Response // Map of userId -> response
}
```

## Goal

Restructure the database to allow for efficient indexing of responses by userId while maintaining backward compatibility with the frontend.

## Steps

1. **Create New Response Structure**

   ```go
   type EventResponse struct {
       UserId   primitive.ObjectID `bson:"userId"`
       Response *Response          `bson:"response"`
   }

   type Event struct {
       // ... other fields
       ResponsesList []EventResponse     `bson:"responses"`      // New field for indexed queries
       Responses    map[string]*Response `bson:"responsesMap"`  // Keep for backward compatibility
   }
   ```

2. **Create Database Migration Script**

   - Create a new migration script in `server/scripts/YYYYMMDD_event_responses_restructure/`
   - Script should:
     1. Query all events
     2. For each event:
        - Convert existing responses map to array format
        - Keep the original map format in `responsesMap`
        - Update the document with both fields

3. **Create Index on ResponsesList.UserId**

   ```javascript
   db.events.createIndex({ "responses.userId": 1 });
   ```

4. **Update Event Creation/Update Logic**

   - When creating/updating events, populate both `responses` array and `responsesMap`
   - Example:
     ```go
     // When adding a response:
     event.ResponsesList = append(event.ResponsesList, EventResponse{
         UserId: userId,
         Response: response,
     })
     event.ResponsesMap[userId.Hex()] = response
     ```

5. **Update Query in getEvents**

   - Replace current query:
     ```go
     bson.M{"responses." + userId.Hex(): bson.M{"$exists": true}}
     ```
   - With new indexed query:
     ```go
     bson.M{"responses.userId": userId}
     ```

6. **Update Response Processing**

   - In getEvents and getEvent endpoints:
     - Return only the `responsesMap` field to maintain frontend compatibility
     - Can be done through MongoDB projection or in Go code

7. **Testing Steps**

   - Write tests to verify:
     - Migration script works correctly
     - New events are created with both response formats
     - Queries using the new index are faster
     - Frontend still works with the returned response format
     - Response updates work correctly

8. **Deployment Considerations**

   - Run migration during low-traffic period
   - Consider running migration in batches if dataset is large
   - Have rollback plan ready
   - Monitor query performance before/after

9. **Future Optimization**
   - Once frontend is updated to use array format:
     - Remove `responsesMap` field
     - Update all queries to use only `responses` array
     - Create new migration to remove old format

## Performance Impact

- Queries filtering by userId in responses will be significantly faster due to index
- Small increase in storage space due to duplicate data during transition
- Slight increase in write operations (maintaining both formats)

## Monitoring

Monitor the following metrics after deployment:

- Query execution time for getEvents
- Storage size of events collection
- Write latency for event updates
