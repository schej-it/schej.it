import Vue from 'vue'

import store from '@/store'
import { serverURL, socketURL, errors } from '@/constants'

/* 
  Date utils 
*/
export const getDateString = (date) => {
  date = new Date(date)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

export const getDateRangeString = (date1, date2) => {
  date1 = new Date(date1)
  date2 = new Date(date2)

  // Correct date2 if time is 12am (because ending at 12am doesn't begin the next day)
  if (date2.getHours() == 0) {
    date2 = getDateDayOffset(date2, -1)
  }

  return getDateString(date1) + ' - ' + getDateString(date2)
}

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

export const getDateWithTime = (date, timeString) => {
  /* Returns a new date object with the given date (e.g. 5/2/2022) and the specified time (e.g. 11:30) */
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

export const getDateWithTimeInt = (date, timeInt, utc=false) => {
  /* Returns a new date object with the given date (e.g. 5/2/2022) and the specified timeInt (e.g. 11.5) */
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

export const splitTime = (timeString) => {
  /* Takes a time string (e.g. 13:30) and splits it into hours and minutes, returning an object of the form { hours, minutes } */
  const [hours, minutes] = timeString.split(':')
  return { hours: parseInt(hours), minutes: parseInt(minutes) }
}

export const getDateDayOffset = (date, offset) => {
  /* Returns the specified date offset by the given number of days (can be positive or negative) */
  date = new Date(date)
  return new Date(date.getTime() + offset * 24 * 60 * 60 * 1000)
}

export const timeIntToTimeText = (timeInt) => {
  /* Converts a timeInt (e.g. 13) to a timeText (e.g. "1 pm") */
  if (timeInt == 0) return '12 am'
  else if (timeInt <= 11) return `${timeInt} am`
  else if (timeInt == 12) return '12 pm'
  return `${timeInt - 12} pm`
}

export const dateToTimeInt = (date) => {
  /* Converts a date to a timeInt (e.g. 9.5) */
  date = new Date(date)
  return date.getHours() + date.getMinutes() / 60
}

export const clampDateToTimeInt = (date, timeInt, type) => {
  /* Clamps the date to the given time, type can either be "upper" or "lower" */
  const diff = dateToTimeInt(date) - timeInt
  if (type === 'upper' && diff < 0) {
    return getDateWithTimeInt(date, timeInt)
  } else if (type === 'lower' && diff > 0) {
    return getDateWithTimeInt(date, timeInt)
  }

  // Return original date
  return date
}

export const dateCompare = (date1, date2) => {
  /* Returns negative if date1 < date2, positive if date2 > date1, and 0 if date1 == date2 */
  date1 = new Date(date1)
  date2 = new Date(date2)
  return date1.getTime() - date2.getTime()
}

export const compareDateDay = (a, b) => {
  // returns -1 if a is before b, 1 if a is after b, 0 otherwise
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

export const isTimeIntBetweenDates = (timeInt, date1, date2) => {
  /* 
  Returns whether the given timeInt is between date1 and date2 
  such that date1.getHour() <= timeInt <= date2.getHour(), accounting 
  for the possibility that date1 and date2 might be on separate days
  */

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

export const areDatesInTimeRanges = (date1, date2, timeRanges) => {
  /* Returns whether both date1 and date2 are fully contained within time1 and time2 */

  const time1 = date1.getTime()
  const time2 = date2.getTime()
  for (const range of timeRanges) {
    if (range.start <= time1 && time2 <= range.end) {
      return true
    }
  }
  return false
}

export const utcTimeToLocalTime = (timeInt, timezoneOffset = new Date().getTimezoneOffset()) => {
  let localTimeInt = timeInt - timezoneOffset / 60
  localTimeInt %= 24
  if (localTimeInt < 0) localTimeInt += 24

  return localTimeInt
}

export const getCalendarEvents = (event) => {
  /* 
    Returns an array of the user's calendar events for the given event, filtering for events
    only between the time ranges of the event and clamping calendar events that extend beyond the time
    ranges
  */

  const timeRanges = []
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

    let curDate = event.startDate
    while (curDate.getTime() < event.endDate.getTime()) {
      const nextDate = getDateDayOffset(curDate, 1)

      let end
      if (event.startTime <= event.endTime) {
        end = getDateWithTimeInt(curDate, event.endTime).getTime()
      } else {
        end = getDateWithTimeInt(nextDate, event.endTime).getTime()
      }
      timeRanges.push({
        start: curDate.getTime(),
        end,
      })

      curDate = nextDate
    }
  } else {
    // New date representation
    timeMin = new Date(event.dates[0]).toISOString()
    timeMax = getDateDayOffset(new Date(event.dates[event.dates.length - 1]), 1).toISOString()

    startTime = utcTimeToLocalTime(event.startTime)
    endTime = utcTimeToLocalTime(event.endTime)

    for (const date of event.dates) {
      const paddedStartTime = String(event.startTime).padStart(2, '0');
      const curDate = new Date(`${date}T${paddedStartTime}:00:00Z`);
      const nextDate = getDateDayOffset(curDate, 1)

      let end
      if (startTime <= endTime) {
        end = getDateWithTimeInt(curDate, endTime).getTime()
      } else {
        end = getDateWithTimeInt(nextDate, endTime).getTime()
      }
      timeRanges.push({
        start: curDate.getTime(),
        end,
      })

      curDate = nextDate
    }
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
      .filter(({ startDate, endDate }) => {
        // Filter calendarEvent based on whether it's completely in between start time and end time
        return areDatesInTimeRanges(startDate, endDate, timeRanges)
      })
  })
}

/* 
  Fetch utils
*/
export const get = (route) => {
  return fetchMethod('GET', route)
}

export const post = (route, body = {}) => {
  return fetchMethod('POST', route, body)
}

export const patch = (route, body = {}) => {
  return fetchMethod('PATCH', route, body)
}

export const _delete = (route, body = {}) => {
  return fetchMethod('DELETE', route, body)
}

export const fetchMethod = (method, route, body = {}) => {
  /* Calls the given route with the give method and body */
  const params = {
    method,
    credentials: 'include',
  }

  if (method !== 'GET') {
    // Add params specific to POST/PATCH/DELETE
    params.headers = {
      'Content-Type': 'application/json',
    }
    params.body = JSON.stringify(body)
  }

  return fetch(serverURL + route, params)
    .then(async (res) => {
      const text = await res.text()

      // Check if response was ok
      if (!res.ok) {
        throw JSON.parse(text)
      }

      // Parse data if it is json, otherwise throw an error
      try {
        return JSON.parse(text)
      } catch (err) {
        throw { error: errors.JsonError }
      }
    })
    .then((data) => {
      return data
    })
}

/*
  Other
*/

export const signInGoogle = (state = null, consent = false) => {
  /* Redirects user to the correct google sign in page */
  const clientId =
    '523323684219-jfakov2bgsleeb6den4ktpohq4lcnae2.apps.googleusercontent.com'
  const redirectUri = `${window.location.origin}/auth`
  const scope = encodeURIComponent(
    'openid email profile https://www.googleapis.com/auth/calendar.calendarlist.readonly https://www.googleapis.com/auth/calendar.events.readonly'
  )

  let stateString = ''
  if (state !== null) {
    state = encodeURIComponent(JSON.stringify(state))
    stateString = `&state=${state}`
  }

  let promptString = ''
  if (consent) {
    promptString = '&prompt=consent'
  } else {
    promptString = '&prompt=select_account'
  }

  window.location.href = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${clientId}&redirect_uri=${redirectUri}&response_type=code&scope=${scope}&access_type=offline${promptString}${stateString}`
}

var timeoutId
export const onLongPress = (element, callback, capture = false) => {
  /* Calls callback() on long press */

  element.addEventListener(
    'touchstart',
    function (e) {
      timeoutId = setTimeout(function () {
        timeoutId = null
        e.stopPropagation()
        callback(e.target)
      }, 500)
    },
    capture
  )

  element.addEventListener(
    'contextmenu',
    function (e) {
      e.preventDefault()
    },
    capture
  )

  element.addEventListener(
    'touchend',
    function () {
      if (timeoutId) clearTimeout(timeoutId)
    },
    capture
  )

  element.addEventListener(
    'touchmove',
    function () {
      if (timeoutId) clearTimeout(timeoutId)
    },
    capture
  )
}

export const isBetween = (value, lower, upper, inclusive = true) => {
  /* Returns whether the given value is between lower and upper */
  if (inclusive) {
    return value >= lower && value <= upper
  } else {
    return value > lower && value < upper
  }
}

export const clamp = (value, lower, upper) => {
  /* Clamps the given value between the given ranges */
  if (value < lower) return lower
  if (value > upper) return upper
  return value
}

export const isPhone = (vuetify) => {
  return vuetify.breakpoint.name === 'xs'
}

export const br = (vuetify, breakpoint) => {
  return vuetify.breakpoint.name === breakpoint
}

export const dataURItoBlob = (dataURI) => {
  // convert base64 to raw binary data held in a string
  // doesn't handle URLEncoded DataURIs - see SO answer #6850276 for code that does this
  var byteString = atob(dataURI.split(',')[1])

  // separate out the mime component
  var mimeString = dataURI.split(',')[0].split(':')[1].split(';')[0]

  // write the bytes of the string to an ArrayBuffer
  var ab = new ArrayBuffer(byteString.length)
  var ia = new Uint8Array(ab)
  for (var i = 0; i < byteString.length; i++) {
    ia[i] = byteString.charCodeAt(i)
  }

  return new Blob([ab], { type: mimeString })
}

export const processEvent = (event) => {
  /* Reformats the given event object to the format we want */
  if (event.startDate) {
    event.startDate = new Date(event.startDate)
    event.endDate = new Date(event.endDate)
    event.startTime = event.startDate.getHours()
    event.endTime = event.endDate.getHours()
  } 
}

export const getCurrentTimezone = () => {
  return new Date()
    .toLocaleTimeString('en-us', { timeZoneName: 'short' })
    .split(' ')[2]
}
