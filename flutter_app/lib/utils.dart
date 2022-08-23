// Returns a new date from the given day and time
DateTime getDateWithTime(DateTime date, double time, {bool local = false}) {
  if (local) {
    return DateTime(date.year, date.month, date.day, time.truncate(),
        ((time - time.truncate()) * 60).truncate());
  }
  return DateTime.utc(date.year, date.month, date.day, time.truncate(),
      ((time - time.truncate()) * 60).truncate());
}

DateTime getLocalDateWithTime(DateTime date, double time) {
  return getDateWithTime(date, time, local: true);
}

// Returns a local datetime with the same day as the given datetime
DateTime getLocalDayFromUtcDay(DateTime day) {
  return DateTime(day.year, day.month, day.day);
}

// Returns if [val] is in between the given range
bool inRange(Comparable val, Comparable min, Comparable max) {
  if (min.compareTo(max) > 0) {
    // Swap min and max if min > max
    final tmp = min;
    min = max;
    max = tmp;
  }
  return val.compareTo(min) >= 0 && val.compareTo(max) <= 0;
}
