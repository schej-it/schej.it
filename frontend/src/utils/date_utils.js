import { get } from './fetch_utils'
/* 
  Date utils 
*/

/** Returns a string representation of the given date, i.e. May 14th is "5/14" */
export const getDateString = (date) => {
  date = new Date(date)
  return `${date.getMonth() + 1}/${date.getDate()}`
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
  let startDate
  let endDate
  if (event.startDate) {
    // Legacy date representation
    startDate = new Date(event.startDate)
    endDate = new Date(event.endDate)
  } else {
    // New date representation
    startDate = getDateWithTimeInt(new Date(event.dates[0]), event.startTime, true)
    endDate = getDateWithTimeInt(new Date(event.dates[event.dates.length - 1]), event.startTime, true)
  }
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

/** Returns a new date object with the given date (e.g. 5/2/2022) and the specified timeInt (e.g. 11.5) */
export const getDateWithTimeInt = (date, timeInt, utc=false) => {
  date = new Date(date)

  const hours = parseInt(timeInt)
  const minutes = (timeInt - hours) * 60
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

/** Returns the specified date offset by the given number of days (can be positive or negative) */
export const getDateDayOffset = (date, offset) => {
  date = new Date(date)
  return new Date(date.getTime() + offset * 24 * 60 * 60 * 1000)
}

/** Converts a timeInt (e.g. 13) to a timeText (e.g. "1 pm") */
export const timeIntToTimeText = (timeInt) => {
  const hours = Math.floor(timeInt)
  const minutesDecimal = timeInt - hours
  const minutesString = minutesDecimal > 0 ? `:${Math.floor(minutesDecimal*60)}` : ''
  

  if (timeInt >= 0 && timeInt < 1) return `12${minutesString} am`
  else if (timeInt < 12) return `${hours}${minutesString} am`
  else if (timeInt >= 12 && timeInt < 13) return `12${minutesString} pm`
  return `${hours - 12}${minutesString} pm`
}

/** Converts a date to a timeInt (e.g. 9.5) */
export const dateToTimeInt = (date) => {
  date = new Date(date)
  return date.getHours() + date.getMinutes() / 60
}

/** Clamps the date to the given time, type can either be "upper" or "lower" */
export const clampDateToTimeInt = (date, timeInt, type) => {
  const diff = dateToTimeInt(date) - timeInt
  if (type === 'upper' && diff < 0) {
    return getDateWithTimeInt(date, timeInt)
  } else if (type === 'lower' && diff > 0) {
    return getDateWithTimeInt(date, timeInt)
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
Returns whether the given timeInt is between date1 and date2 
such that date1.getHour() <= timeInt <= date2.getHour(), accounting 
for the possibility that date1 and date2 might be on separate days
*/
export const isTimeIntBetweenDates = (timeInt, date1, date2) => {
  const hour1 = date1.getHours()
  const hour2 = date2.getHours()

  if (hour1 <= hour2) {
    return hour1 <= timeInt && timeInt <= hour2
  } else {
    return (
      (hour1 <= timeInt && timeInt < 24) || (0 <= timeInt && timeInt <= hour2)
    )
  }
}

/** Converts a utc time int to a local time int based on the timezoneOffset */
export const utcTimeToLocalTime = (timeInt, timezoneOffset = new Date().getTimezoneOffset()) => {
  let localTimeInt = timeInt - timezoneOffset / 60
  localTimeInt %= 24
  if (localTimeInt < 0) localTimeInt += 24

  return localTimeInt
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
export const getCalendarEvents = (event) => {
  let timeMin
  let timeMax
  let startTime
  let endTime
  if (event.startDate) {
    // Legacy date representation
    timeMin = event.startDate.toISOString()
    timeMax = getDateDayOffset(event.endDate, 1).toISOString()

    startTime = event.startTime
    endTime = event.endTime
  } else {
    // New date representation
    timeMin = new Date(event.dates[0]).toISOString()
    timeMax = getDateDayOffset(new Date(event.dates[event.dates.length - 1]), 1).toISOString()

    startTime = utcTimeToLocalTime(event.startTime)
    endTime = utcTimeToLocalTime(event.endTime)
  }

  return get(
    `/user/calendar?timeMin=${timeMin}&timeMax=${timeMax}`
  ).then((data) => {
    return data
      .map((calendarEvent) => {
        // If calendarEvent has a time int between the start and end dates, clamp it based on whether it's the starttime or endtime
        calendarEvent.startDate = new Date(calendarEvent.startDate)
        calendarEvent.endDate = new Date(calendarEvent.endDate)
        const { startDate, endDate } = calendarEvent
        if (isTimeIntBetweenDates(startTime, startDate, endDate)) {
          return {
            ...calendarEvent,
            startDate:
              startDate.getHours() <= startTime
                ? getDateWithTimeInt(startDate, startTime)
                : getDateWithTimeInt(endDate, startTime),
          }
        } else if (isTimeIntBetweenDates(endTime, startDate, endDate)) {
          return {
            ...calendarEvent,
            endDate:
              endDate.getHours() >= endTime
                ? getDateWithTimeInt(endDate, endTime)
                : getDateWithTimeInt(startDate, endTime),
          }
        } else {
          return calendarEvent
        }
      })
      .filter((calendarEvent) => {
        // Filter calendarEvent based on whether it's completely in between start time and end time

        // calendarEventDayStart is a date representation of the event start time for the day the calendar event takes place
        const calendarEventDayStart = getDateWithTimeInt(calendarEvent.startDate, startTime)
        if (calendarEventDayStart.getTime() > calendarEvent.startDate.getTime()) {
          // Go back a day if calendarEventDayStart is past the calendarEvent start time
          calendarEventDayStart.setDate(calendarEventDayStart.getDate() - 1);
        }

        // calendarEventDayEnd is a date representation of the event end time for the day the calendar event takes place
        const calendarEventDayEnd = new Date(calendarEventDayStart)
        if (endTime > startTime) {
          calendarEventDayEnd.setHours(calendarEventDayEnd.getHours() + (endTime - startTime));
        } else {
          calendarEventDayEnd.setHours(calendarEventDayEnd.getHours() + (endTime + 24-startTime));
        }

        const isBetween = calendarEvent.startDate >= calendarEventDayStart && calendarEvent.endDate <= calendarEventDayEnd
        return isBetween
      })
  })
}