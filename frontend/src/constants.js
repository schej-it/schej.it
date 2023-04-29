// Urls
export const serverURL =
  process.env.NODE_ENV === "development" ? "http://localhost:3002" : "/api"
export const socketURL =
  process.env.NODE_ENV === "development" ? "http://localhost:3002" : "/"

// Errors enum
export const errors = Object.freeze({
  JsonError: "json-error",
  NotSignedIn: "not-signed-in",
  UserDoesNotExist: "user-does-not-exist",
  EventNotFound: "event-not-found",
})

// Auth types
export const authTypes = Object.freeze({
  EVENT_ADD_AVAILABILITY: "event-add-availability",
  EVENT_SIGN_IN: "event-sign-in",
  EVENT_SCHEDULE: "event-schedule",
})
