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
  ADD_CALENDAR_ACCOUNT: "add-calendar-account",
  ADD_CALENDAR_ACCOUNT_FROM_EDIT: "add-calendar-account-from-edit",
})

export const eventTypes = Object.freeze({
  SPECIFIC_DATES: "specific_dates",
  DOW: "dow",
})

export const dayIndexToDayString = Object.freeze([
  "2018-06-17", // Sunday
  "2018-06-18", // Monday
  "2018-06-19", // Tuesday
  "2018-06-20", // Wednesday
  "2018-06-21", // Thursday
  "2018-06-22", // Friday
  "2018-06-23", // Saturday
])
