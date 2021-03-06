// Returns a new date from the given day and time
DateTime getDateWithTime(DateTime date, double time, {bool local = false}) {
  if (local) {
    return DateTime(date.year, date.month, date.day, time.truncate(),
        ((time - time.truncate()) * 60).truncate());
  }
  return DateTime.utc(date.year, date.month, date.day, time.truncate(),
      ((time - time.truncate()) * 60).truncate());
}

// Returns a local datetime with the same day as the given datetime
DateTime getLocalDayFromUtcDay(DateTime day) {
  return DateTime(day.year, day.month, day.day);
}
