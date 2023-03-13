import { get } from './fetch_utils'
import { isBetween } from './general_utils'
/* 
  Date utils 
*/

/** Returns a string representation of the given date, i.e. May 14th is "5/14" */
export const getDateString = (date) => {
  date = new Date(date)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

/** Returns an ISO formatted date string */
export const getISODateString = (date, utc=false) => {
  date = new Date(date)
  if (utc) {
    return date.toISOString().substring(0, 10)
  }

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

/** Returns a string representing date range from date1 to date2, i.e. "5/14 - 5/27" */
export const getDateRangeString = (date1, date2) => {
  date1 = new Date(date1)
  date2 = new Date(date2)

  // Correct date2 if time is 12am (because ending at 12am doesn't begin the next day)
  if (date2.getHours() == 0) {
    date2 = getDateDayOffset(date2, -1)
  }

  return getDateString(date1) + ' - ' + getDateString(date2)
}

/** Returns a string representing the date range for the provided event */
export const getDateRangeStringForEvent = (event) => {
  const startDate = new Date(event.dates[0])
  const endDate = new Date(event.dates[event.dates.length-1])
  return getDateRangeString(startDate, endDate);
}

/** Returns a new date object with the given date (e.g. 5/2/2022) and the specified time (e.g. "11:30") */
export const getDateWithTime = (date, timeString) => {
  date = new Date(date)

  const { hours, minutes } = splitTime(timeString)
  return new Date(
    date.getFullYear(),
    date.getMonth(),
    date.getDate(),
    hours,
    minutes
  )
}

/** Returns a new date object with the given date (e.g. 5/2/2022) and the specified timeNum (e.g. 11.5) */
export const getDateWithTimeNum = (date, timeNum, utc=false) => {
  date = new Date(date)

  const hours = parseInt(timeNum)
  const minutes = (timeNum - hours) * 60
  if (!utc) {
    return new Date(
      date.getFullYear(),
      date.getMonth(),
      date.getDate(),
      hours,
      minutes
    )
  } else {
    return new Date(Date.UTC(
      date.getUTCFullYear(),
      date.getUTCMonth(),
      date.getUTCDate(),
      hours,
      minutes
    ))
  }
}

/** Takes a time string (e.g. 13:30) and splits it into hours and minutes, returning an object of the form { hours, minutes } */
export const splitTime = (timeString) => {
  const [hours, minutes] = timeString.split(':')
  return { hours: parseInt(hours), minutes: parseInt(minutes) }
}

/** Takes a timeNum (e.g. 9.5) and splits it into hours and minutes, returning an object of the form { hours, minutes } */
export const splitTimeNum = (timeNum) => {
  const hours = Math.floor(timeNum)
  const minutes = Math.floor((timeNum - hours) * 60)
  return { hours, minutes }
}

/** Returns the specified date offset by the given number of days (can be positive or negative) */
export const getDateDayOffset = (date, offset) => {
  date = new Date(date)
  return new Date(date.getTime() + offset * 24 * 60 * 60 * 1000)
}

/** Returns the specified date offset by the given number of hours */
export const getDateHoursOffset = (date, hoursOffset) => {
  const { hours, minutes } = splitTimeNum(hoursOffset)
  const newDate = new Date(date)
  newDate.setHours(newDate.getHours() + hours)
  newDate.setMinutes(newDate.getMinutes() + minutes)
  return newDate
}

/** Converts a timeNum (e.g. 13) to a timeText (e.g. "1 pm") */
export const timeNumToTimeText = (timeNum) => {
  const hours = Math.floor(timeNum)
  const minutesDecimal = timeNum - hours
  const minutesString = minutesDecimal > 0 ? `:${String(Math.floor(minutesDecimal*60)).padStart(2, '0')}` : ''
  

  if (timeNum >= 0 && timeNum < 1) return `12${minutesString} am`
  else if (timeNum < 12) return `${hours}${minutesString} am`
  else if (timeNum >= 12 && timeNum < 13) return `12${minutesString} pm`
  return `${hours - 12}${minutesString} pm`
}

/** Converts a timeNum (e.g. 9.5) to a timeString (e.g. 09:30:00) */
export const timeNumToTimeString = (timeNum) => {
  const hours = Math.floor(timeNum)
  const minutesDecimal = timeNum - hours
  const paddedHours = String(hours).padStart(2, '0');
  const paddedMinutes = String(Math.floor(minutesDecimal*60)).padStart(2, '0');

  return `${paddedHours}:${paddedMinutes}:00`
}

/** Converts a date to a timeNum (e.g. 9.5) */
export const dateToTimeNum = (date, utc=false) => {
  date = new Date(date)
  if (utc) {
    return date.getUTCHours() + date.getUTCMinutes() / 60
  }
  return date.getHours() + date.getMinutes() / 60
}

/** Clamps the date to the given time, type can either be "upper" or "lower" */
export const clampDateToTimeNum = (date, timeNum, type) => {
  const diff = dateToTimeNum(date) - timeNum
  if (type === 'upper' && diff < 0) {
    return getDateWithTimeNum(date, timeNum)
  } else if (type === 'lower' && diff > 0) {
    return getDateWithTimeNum(date, timeNum)
  }

  // Return original date
  return date
}

/** Returns negative if date1 < date2, positive if date2 > date1, and 0 if date1 == date2 */
export const dateCompare = (date1, date2) => {
  date1 = new Date(date1)
  date2 = new Date(date2)
  return date1.getTime() - date2.getTime()
}

/** returns -1 if a is before b, 1 if a is after b, 0 otherwise */
export const compareDateDay = (a, b) => {
  a = new Date(a)
  b = new Date(b)
  if (a.getFullYear() !== b.getFullYear()) {
    return a.getFullYear() - b.getFullYear()
  } else if (a.getMonth() !== b.getMonth()) {
    return a.getMonth() - b.getMonth()
  } else {
    return a.getDate() - b.getDate()
  }
}

/**
Returns whether the given timeNum is between date1 and date2 
such that date1.getHour() <= timeNum <= date2.getHour(), accounting 
for the possibility that date1 and date2 might be on separate days
*/
export const isTimeNumBetweenDates = (timeNum, date1, date2) => {
  const hour1 = date1.getHours()
  const hour2 = date2.getHours()

  if (hour1 <= hour2) {
    return hour1 <= timeNum && timeNum <= hour2
  } else {
    return (
      (hour1 <= timeNum && timeNum < 24) || (0 <= timeNum && timeNum <= hour2)
    )
  }
}

/** Returns whether date is in between startDate and startDate + duration (in hours) */
export const isDateInRange = (date, startDate, duration) => {
  const endDate = new Date(startDate)
  endDate.setHours(endDate.getHours() + duration)
  return startDate <= date && date <= endDate
}

/** Converts a utc time int to a local time int based on the timezoneOffset */
export const utcTimeToLocalTime = (timeNum, timezoneOffset = new Date().getTimezoneOffset()) => {
  let localTimeNum = timeNum - timezoneOffset / 60
  localTimeNum %= 24
  if (localTimeNum < 0) localTimeNum += 24

  return localTimeNum
}

/** Returns a string representing the current timezone */
export const getCurrentTimezone = () => {
  return new Date()
    .toLocaleTimeString('en-us', { timeZoneName: 'short' })
    .split(' ')[2]
}

/** 
  Returns an array of the user's calendar events for the given event, filtering for events
  only between the time ranges of the event and clamping calendar events that extend beyond the time
  ranges
*/
export const getCalendarEventsByDay = async (event) => {
  let timeMin = new Date(event.dates[0]).toISOString()
  let timeMax = getDateDayOffset(new Date(event.dates[event.dates.length - 1]), 2).toISOString()

  // Fetch calendar events from Google Calendar
  const calendarEvents = await get(
    `/user/calendar?timeMin=${timeMin}&timeMax=${timeMax}`
  ).then(events => events.map(e => {
    e.startDate = new Date(e.startDate)
    e.endDate = new Date(e.endDate)
    return e
  }))
  calendarEvents.sort((a, b) => dateCompare(a.startDate, b.startDate))

  // Iterate through all dates and add calendar events to array
  const calendarEventsByDay = []
  for (const i in event.dates) {
    calendarEventsByDay[i] = []

    if (calendarEvents.length == 0) break

    const start = new Date(event.dates[i])
    const end = new Date(start)
    end.setHours(start.getHours() + event.duration)

    // Keep iterating through calendar events until it's empty or there are no more events for the current date
    while (calendarEvents.length > 0 && end > calendarEvents[0].startDate) {
      const [calendarEvent] = calendarEvents.splice(0, 1)

      // Check if calendar event overlaps with event time ranges
      const startDateWithinRange = isBetween(calendarEvent.startDate, start, end)
      const endDateWithinRange = isBetween(calendarEvent.endDate, start, end)
      const rangeWithinCalendarEvent = isBetween(start, calendarEvent.startDate, calendarEvent.endDate) && isBetween(end, calendarEvent.startDate, calendarEvent.endDate)
      if (startDateWithinRange || endDateWithinRange || rangeWithinCalendarEvent) {
        const rangeStartWithinCalendarEvent = isBetween(start, calendarEvent.startDate, calendarEvent.endDate)
        const rangeEndWithinCalendarEvent = isBetween(end, calendarEvent.startDate, calendarEvent.endDate)
        if (rangeStartWithinCalendarEvent) {
          // Clamp calendarEvent start
          calendarEvent = { ...calendarEvent, startDate: start }
        }
        if (rangeEndWithinCalendarEvent) {
          // Clamp calendarEvent end
          calendarEvent = { ...calendarEvent, endDate: end }
        }

        // The number of hours since start time
        const hoursOffset =
          (calendarEvent.startDate.getTime() - start.getTime()) /
          (1000 * 60 * 60)

        // The length of the event in hours
        const hoursLength =
          (calendarEvent.endDate.getTime() -
            calendarEvent.startDate.getTime()) /
          (1000 * 60 * 60)

        // Don't display event if the event is 0 hours long
        if (hoursLength == 0) continue

        calendarEventsByDay[i].push({
          ...calendarEvent,
          hoursOffset,
          hoursLength,
        })
      }
    }
  }

  return calendarEventsByDay
}