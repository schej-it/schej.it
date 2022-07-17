// Returns a new date from the given day and time
DateTime getDateWithTime(DateTime date, double time) {
  return DateTime.utc(date.year, date.month, date.day, time.truncate(),
      ((time - time.truncate()) * 60).truncate());
}
