class CalendarEvent {
  final String title;
  final DateTime startDate;
  final DateTime endDate;

  double get startTime => startDate.hour + startDate.minute / 60;
  double get endTime => endDate.hour + endDate.minute / 60;

  const CalendarEvent({
    required this.title,
    required this.startDate,
    required this.endDate,
  });

  @override
  String toString() {
    return '{CalendarEvent title:"$title" startDate:$startDate endDate:$endDate}';
  }
}

class CalendarEvents {
  // A map mapping a string representing the day ("7-15-2022") to the events
  // occurring on that day
  final Map<String, List<CalendarEvent>> _eventsByDay =
      <String, List<CalendarEvent>>{};

  CalendarEvents({
    required List<CalendarEvent> events,
  }) {
    for (CalendarEvent event in events) {
      String dayString = _getDayString(event.startDate);
      if (_eventsByDay[dayString] == null) {
        _eventsByDay[dayString] = [event];
      } else {
        _eventsByDay[dayString]!.add(event);
      }
    }
  }

  List<CalendarEvent>? getEventsForDay(DateTime date) {
    return _eventsByDay[_getDayString(date)];
  }

  String _getDayString(DateTime date) {
    return date.toLocal().toIso8601String().substring(0, 10);
  }
}
