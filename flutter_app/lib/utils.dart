// bool isSameDay(DateTime d1, DateTime d2) {
//   return d1.day == d2.day && d1.month == d2.month && d1.year == d2.year;
// }

DateTime getDateWithTime(DateTime date, double time) {
  return DateTime(date.year, date.month, date.day, time.truncate(),
      ((time - time.truncate()) * 60).truncate());
}
