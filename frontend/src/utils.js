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

export const timeIntToTimeString = (timeInt) => {
  /* Converts a timeInt (e.g. 13) to a timeString (e.g. "1 pm") */
  if (timeInt == 0) return '12 am'
  else if (timeInt <= 11) return `${timeInt} am`
  else if (timeInt == 12) return '12 pm'
  return `${timeInt - 12} pm`
}