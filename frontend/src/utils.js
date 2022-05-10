export const getDateString = (date) => {
  return `${date.getMonth()+1}/${date.getDate()}`
}

export const getDateRangeString = (date1, date2) => {
  console.log(date1.getMonth())
  return  getDateString(date1) + ' - ' + getDateString(date2)
}