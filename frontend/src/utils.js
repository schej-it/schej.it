export const serverURL = process.env.NODE_ENV === 'development' ? 'http://localhost:3000' : '/api'
export const socketURL = process.env.NODE_ENV === 'development' ? 'http://localhost:3000' : '/'

export const getDateString = (date) => {
  return `${date.getMonth()+1}/${date.getDate()}`
}

export const getDateRangeString = (date1, date2) => {
  return  getDateString(date1) + ' - ' + getDateString(date2)
}

export const getDateWithTime = (date, timeString) => {
  /* Returns a new date object with the given date (e.g. 5/2/2022) and the specified time (e.g. 11:30) */
  
  const { hours, minutes } = splitTime(timeString)
  return new Date(date.getFullYear(), date.getMonth(), date.getDate(), hours, minutes)
}

export const splitTime = (timeString) => {
  /* Takes a time string (e.g. 13:30) and splits it into hours and minutes, returning an object of the form { hours, minutes } */
  const [hours, minutes] = timeString.split(':')
  return { hours: parseInt(hours), minutes: parseInt(minutes) }
}

export const getDateDayOffset = (date, offset) => {
  /* Returns the specified date offset by the given number of days (can be positive or negative) */
  return new Date(date.getTime() + offset * 24*60*60*1000)
}

export const timeIntToTimeText = (timeInt) => {
  /* Converts a timeInt (e.g. 13) to a timeText (e.g. "1 pm") */
  if (timeInt == 0) return '12 am'
  else if (timeInt <= 11) return `${timeInt} am`
  else if (timeInt == 12) return '12 pm'
  return `${timeInt - 12} pm`
}

export const dateCompare = (date1, date2) => {
  /* Returns negative if date1 < date2, positive if date2 > date1, and 0 if date1 == date2 */
  return date1.getTime() - date2.getTime()
}

export const get = (route) => {
  return fetchMethod('GET', route)
}

export const post = (route, body={}) => {
  return fetchMethod('POST', route, body)
}

export const patch = (route, body={}) => {
  return fetchMethod('PATCH', route, body)
}

export const _delete = (route, body={}) => {
  return fetchMethod('DELETE', route, body)
}

export const fetchMethod = (method, route, body={}) => {
  /* Calls the given route with the give method and body */
  const params = {
    method,
    credentials: 'include',
  }

  if (method !== 'GET') {
    // Add params specific to POST/PATCH/DELETE
    params.headers = {
      'Content-Type': 'application/json'
    }
    params.body = JSON.stringify(body)
  }

  return fetch(serverURL + route, params).then(async res => {
    // Parse data if it is json, otherwise just return an empty object
    const text = await res.text()
    try {
      return JSON.parse(text)
    } catch (err) {
      return {}
    }
  }).then(data => {
    // Throw an error if one occurred
    if (data.error)
      throw data.error
    
    return data
  })
}